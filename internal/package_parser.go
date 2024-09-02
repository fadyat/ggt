package internal

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/fadyat/ggt/internal/lo"
)

var (
	ErrNoMissingTests = errors.New("no missing tests")
)

// PackageParser is required in cases, when we need to generate the
// testcase for some method from one file, but the struct definition
// is stored in another file. In this case, we need to perform the lazy
// parsing of all files inside the package, until we don't fine the required
// struct definition.
//
// This struct is responsible for parsing the package and storing the
// parsed files in the memory, so that we can access them later.
type PackageParser struct {
	flags *Flags

	inputFileSet  *token.FileSet
	outputFileSet *token.FileSet
	inputAst      *ast.File
	outputAst     *ast.File

	// currentPackageFileFileSet will contain the parsed files of the package
	// which is being processed, caching only the last one, because we are
	// doing the lazy parsing of the package.
	currentPackageFileFileSet *token.FileSet
	currentPackageFileAst     *ast.File
}

func NewParser(flags *Flags) *PackageParser {
	return &PackageParser{
		flags: flags,
	}
}

func (p *PackageParser) GenerateMissingTests(inputFnFilters ...func(*Fn) bool) (f *File, err error) {
	p.inputFileSet, p.inputAst, err = p.parseFile(p.flags.InputFile)
	if err != nil {
		return nil, fmt.Errorf("parse input file: %w", err)
	}

	p.outputFileSet, p.outputAst, err = p.parseFile(p.flags.OutputFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("parse output file: %w", err)
	}

	missingTests := p.getMissingTests(p.createInputFnFilter(inputFnFilters))
	if len(missingTests) == 0 {
		return nil, ErrNoMissingTests
	}

	if err = p.getStructsForMethods(missingTests); err != nil {
		return nil, fmt.Errorf("get structs for methods: %w", err)
	}

	file := &File{
		Functions: missingTests,
	}

	if p.outputAst == nil {
		file.PackageName = p.inputAst.Name.Name
		file.Imports = p.getImports(p.inputAst)
	}

	return file, nil
}

func (p *PackageParser) createInputFnFilter(inputFnFilters []func(*Fn) bool) func(*Fn) bool {
	if len(inputFnFilters) == 0 {
		return func(*Fn) bool { return true }
	}

	return func(fn *Fn) bool {
		for _, filter := range inputFnFilters {
			if !filter(fn) {
				return false
			}
		}

		return true
	}
}

func (p *PackageParser) getImports(f *ast.File) []*Import {
	fromFile := lo.Map(f.Imports, func(imp *ast.ImportSpec, _ int) *Import {
		var alias string
		if imp.Name != nil {
			alias = imp.Name.Name
		}

		return newImport(alias, imp.Path.Value)
	})

	// appending empty import in cases when no imports exist, but
	// need to generate imports from the template
	return append(fromFile, newImport("", ""))
}

func (p *PackageParser) parseFile(path string) (*token.FileSet, *ast.File, error) {
	tokenFileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(tokenFileSet, path, nil, parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}

	return tokenFileSet, astFile, nil
}

func (p *PackageParser) getMissingTests(inputFilter func(*Fn) bool) []*Fn {
	inputFuncs := getFuncs(p.inputFileSet, p.inputAst, func(fs *token.FileSet, decl *ast.FuncDecl) (*Fn, bool) {
		ff := parseFn(fs, decl)
		ff.generateFriendlyNames(ff.Args)
		ff.generateFriendlyNames(ff.Results)
		return ff, inputFilter(ff)
	})

	outputFuncs := getFuncs(p.outputFileSet, p.outputAst, func(fs *token.FileSet, decl *ast.FuncDecl) (*Fn, bool) {
		ff := parseFn(fs, decl)
		return ff, true
	})

	return lo.FilterMap(inputFuncs, func(item *Fn, _ int) (*Fn, bool) {
		return item, !lo.ContainsBy(outputFuncs, func(out *Fn) bool {
			return item.TestName() == out.Name
		})
	})
}

func (p *PackageParser) parseAndMatchStructs(missingStructsFn map[string]*Fn) {
	fileStructs := lo.SliceToMap(
		getStructs(p.currentPackageFileFileSet, p.currentPackageFileAst, parseStructs),
		func(s *Struct) (string, *Struct) { return s.Name, s },
	)

	for _, method := range missingStructsFn {
		structType := method.structTypeBasedOnReceiver()
		if s, ok := fileStructs[structType]; ok {
			method.Struct = s
			delete(missingStructsFn, method.TestName())
		}
	}
}

func (p *PackageParser) getStructsForMethods(methods []*Fn) error {
	missingStructsFn := lo.SliceToMap(
		lo.FilterMap(methods, func(method *Fn, _ int) (*Fn, bool) {
			return method, method.Receiver != nil
		}),
		func(f *Fn) (string, *Fn) { return f.TestName(), f },
	)

	if len(missingStructsFn) == 0 {
		return nil
	}

	inputFileDir := filepath.Dir(p.flags.InputFile)
	packageFiles, err := listPackageFiles(inputFileDir, defaultExcludeFunc(p.flags.InputFile))
	if err != nil {
		return fmt.Errorf("list package files: %w", err)
	}

	p.currentPackageFileFileSet, p.currentPackageFileAst = p.inputFileSet, p.inputAst

	p.parseAndMatchStructs(missingStructsFn)
	if len(missingStructsFn) == 0 {
		return nil
	}

	// doing the same logic, but for the rest of the files in the package
	for _, file := range packageFiles {
		p.currentPackageFileFileSet, p.currentPackageFileAst, err = p.parseFile(filepath.Join(inputFileDir, file))
		if err != nil {
			return fmt.Errorf("parse file: %w", err)
		}

		p.parseAndMatchStructs(missingStructsFn)
		if len(missingStructsFn) == 0 {
			return nil
		}
	}

	return fmt.Errorf(
		"missing structs for the following methods: %s",
		lo.MapToSlice(missingStructsFn, func(k string, _ *Fn) string { return k }),
	)
}

func getStructs(fs *token.FileSet, f *ast.File, parser func(*token.FileSet, *ast.GenDecl) []*Struct) []*Struct {
	if f == nil {
		return nil
	}

	return lo.FlatMap(f.Decls, func(item ast.Decl, _ int) []*Struct {
		gen, ok := item.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			return nil
		}

		return parser(fs, gen)
	})
}

func parseStructs(fs *token.FileSet, decl *ast.GenDecl) []*Struct {
	var structs = make([]*Struct, 0)

	for _, spec := range decl.Specs {
		typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
		if !isTypeSpec {
			continue
		}

		structType, isStructType := typeSpec.Type.(*ast.StructType)
		if !isStructType {
			continue
		}

		var s = newStruct(typeSpec.Name.Name)
		if typeSpec.TypeParams != nil {
			s.Generics = lo.FlatMap(typeSpec.TypeParams.List, func(typeParam *ast.Field, _ int) []*Identifier {
				typeParamType := getTypeName(fs, typeParam.Type)
				return lo.Map(typeParam.Names, func(name *ast.Ident, _ int) *Identifier {
					return newIdentifier(name.Name, typeParamType)
				})
			})
		}

		s.Fields = lo.FlatMap(structType.Fields.List, func(field *ast.Field, _ int) []*Identifier {
			fieldType := getTypeName(fs, field.Type)
			if len(field.Names) == 0 {
				var split = strings.Split(fieldType, ".")
				var fieldName = split[len(split)-1]
				return []*Identifier{newIdentifier(fieldName, fieldType)}
			}

			return lo.Map(field.Names, func(name *ast.Ident, _ int) *Identifier {
				return newIdentifier(name.Name, fieldType)
			})
		})

		structs = append(structs, s)
	}

	return structs
}

func listPackageFiles(path string, exclude func(s string) bool) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("read directory: %w", err)
	}

	var result = make([]string, 0)
	for _, file := range files {
		var fileName = file.Name()
		if !exclude(fileName) {
			result = append(result, fileName)
		}
	}

	return result, nil
}

func defaultExcludeFunc(inputFile string) func(string) bool {
	return func(s string) bool {
		if s == inputFile {
			return true
		}

		if !strings.HasSuffix(s, ".go") {
			return true
		}

		if strings.HasSuffix(s, "_test.go") {
			return true
		}

		return false
	}
}

func getFuncs(
	fs *token.FileSet,
	f *ast.File,
	parser func(*token.FileSet, *ast.FuncDecl) (*Fn, bool),
) []*Fn {
	if f == nil {
		return nil
	}

	return lo.FilterMap(f.Decls, func(item ast.Decl, _ int) (*Fn, bool) {
		fn, ok := item.(*ast.FuncDecl)
		if !ok {
			return nil, false
		}

		return parser(fs, fn)
	})
}

func parseFn(fs *token.FileSet, f *ast.FuncDecl) *Fn {
	var function = newFn(f.Name.Name)

	if f.Recv != nil {
		var (
			receiverType = getTypeName(fs, f.Recv.List[0].Type)
			receiverName = string(receiverType[0]) // todo: mb can be better
		)

		if len(f.Recv.List[0].Names) > 0 {
			receiverName = f.Recv.List[0].Names[0].Name
		}

		function.Receiver = newIdentifier(receiverName, receiverType)
	}

	if f.Type.TypeParams != nil {
		function.Generics = lo.FlatMap(f.Type.TypeParams.List, func(typeParam *ast.Field, _ int) []*Identifier {
			typeParamType := getTypeName(fs, typeParam.Type)
			return lo.Map(typeParam.Names, func(name *ast.Ident, _ int) *Identifier {
				return newIdentifier(name.Name, typeParamType)
			})
		})
	}

	function.Args = lo.FlatMap(f.Type.Params.List, func(arg *ast.Field, _ int) []*Identifier {
		argType := getTypeName(fs, arg.Type)
		if len(arg.Names) == 0 {
			return []*Identifier{newIdentifier("", argType)}
		}

		return lo.Map(arg.Names, func(name *ast.Ident, _ int) *Identifier {
			return newIdentifier(name.Name, argType)
		})
	})

	if f.Type.Results != nil {
		function.Results = lo.FlatMap(f.Type.Results.List, func(res *ast.Field, _ int) []*Identifier {
			resType := getTypeName(fs, res.Type)
			if len(res.Names) == 0 {
				return []*Identifier{newIdentifier("", resType)}
			}

			return lo.Map(res.Names, func(name *ast.Ident, _ int) *Identifier {
				return newIdentifier(name.Name, resType)
			})
		})
	}

	return function
}

func getTypeName(fs *token.FileSet, expr ast.Expr) string {
	var b bytes.Buffer
	_ = printer.Fprint(&b, fs, expr)
	return b.String()
}
