package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"io"
	"log"
	"log/slog"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/fadyat/ggt/v2/lo"
)

const (
	mode = packages.NeedName |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo |
		packages.NeedImports

	pattern = "./..."
)

var (
	lg   = slog.Default()
	fset = token.NewFileSet()
)

//	todo: if file path is fully provided, need also include the test file automatically
//	 command-line-arguments if file path is fully provided
//
// todo: add checks, that file exists, now disabling this feature
func withPatterns() []string {
	if strings.HasSuffix(pattern, "_test.go") {
		return []string{strings.TrimSuffix(pattern, "_test.go") + ".go", pattern}
	}

	if strings.HasSuffix(pattern, ".go") {
		return []string{pattern, strings.TrimSuffix(pattern, ".go") + "_test.go"}
	}

	return []string{pattern}
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelError)
	flag.Parse()

	cfg := &packages.Config{Fset: fset, Mode: mode, Dir: flag.Arg(0), Tests: true}

	pkgs, err := packages.Load(cfg, pattern)
	packages.PrintErrors(pkgs)
	if err != nil {
		log.Fatalf("loading packages: %v", err)
	}

	packagePairs := make(map[string]*packagePair)
	for _, pkg := range pkgs {
		// skipping compiled test files, required for the test files
		if strings.HasSuffix(pkg.PkgPath, ".test") {
			continue
		}

		pkgPath := strings.TrimSuffix(pkg.PkgPath, "_test")
		if _, ok := packagePairs[pkgPath]; !ok {
			packagePairs[pkgPath] = &packagePair{}
		}

		pkgInfo := processPackage(pkg)
		packagePairs[pkgPath].add(pkgInfo)
	}

	for _, pkgPair := range packagePairs {
		renderPackage(pkgPair, func(_ string) (io.WriteCloser, error) {
			return &stdoutWriter{}, nil
		})
	}
}

type packagePair struct {
	normal *packageInfo
	test   *packageInfo
}

func (p *packagePair) add(pi *packageInfo) {
	if strings.HasSuffix(pi.pkg.PkgPath, "_test") {
		p.test = pi
	} else {
		p.normal = pi
	}
}

type packageInfo struct {
	pkg *packages.Package

	files map[string]*fileInfo
}

type fileInfo struct {
	file *ast.File

	funcs map[string]*funcInfo
}

type funcInfo struct {
	funcDecl *ast.FuncDecl

	params     []*paramInfo
	typeParams []*paramInfo
	results    []*paramInfo
	recv       *recvInfo
}

func (f *funcInfo) TestName() string {
	var sb strings.Builder
	sb.WriteString("Test_")
	if f.recv != nil {
		sb.WriteString(fmt.Sprintf("%s_", f.recv.recvParam.name))
	}

	sb.WriteString(f.funcDecl.Name.Name)
	return sb.String()
}

func (f *funcInfo) RenderUnchanged() string {
	var b bytes.Buffer
	_ = printer.Fprint(&b, fset, f.funcDecl)
	return b.String()
}

func processPackage(pkg *packages.Package) *packageInfo {
	lg.Debug("package", slog.String("path", pkg.PkgPath))

	pi := packageInfo{pkg: pkg, files: make(map[string]*fileInfo)}
	for _, fileAst := range pkg.Syntax {
		fi := processFile(fileAst, pkg.TypesInfo)
		pi.files[fset.Position(fileAst.Package).Filename] = fi
	}

	return &pi
}

func processFile(fileAst *ast.File, tinfo *types.Info) *fileInfo {
	lg.Debug("file", slog.String("path", fset.Position(fileAst.Package).Filename))

	fi := fileInfo{file: fileAst, funcs: make(map[string]*funcInfo)}
	for _, decl := range fileAst.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			fni := processFuncDecl(funcDecl, tinfo)
			fi.funcs[funcDecl.Name.Name] = fni
		}
	}

	return &fi
}

func processFuncDecl(fd *ast.FuncDecl, tinfo *types.Info) *funcInfo {
	return &funcInfo{
		funcDecl: fd,
		params: lo.FlatMap(fd.Type.Params.List, func(item *ast.Field, _ int) []*paramInfo {
			return fieldParams(item, tinfo)
		}),
		typeParams: lo.
			If[[]*paramInfo](fd.Type.TypeParams == nil, nil).
			ElseF(func() []*paramInfo {
				return lo.FlatMap(fd.Type.TypeParams.List, func(item *ast.Field, _ int) []*paramInfo {
					return fieldParams(item, tinfo)
				})
			}),
		results: lo.
			If[[]*paramInfo](fd.Type.Results == nil, nil).
			ElseF(func() []*paramInfo {
				return lo.FlatMap(fd.Type.Results.List, func(item *ast.Field, _ int) []*paramInfo {
					return fieldParams(item, tinfo)
				})
			}),
		recv: lo.
			If[*recvInfo](fd.Recv == nil, nil).
			ElseF(func() *recvInfo {
				return processRecv(fd.Recv.List[0], tinfo)
			}),
	}
}

type paramInfo struct {
	field *ast.Field

	name  string
	ftype *typeInfo
}

type typeInfo struct {
	syntaxType   ast.Expr
	semanticType types.Type
}

func (t *typeInfo) String() string {
	if t.syntaxType != nil {
		var b bytes.Buffer
		_ = printer.Fprint(&b, fset, t.syntaxType)
		return b.String()
	}

	return t.semanticTypeString(t.semanticType)
}

func (t *typeInfo) semanticTypeString(stype types.Type) string {
	switch st := stype.(type) {
	case *types.Pointer:
		return "*" + t.semanticTypeString(st.Elem())
	case *types.Named:
		packagePath := strings.Split(st.String(), "/")
		packageWithStruct := strings.Split(packagePath[len(packagePath)-1], ".")
		structName := packageWithStruct[len(packageWithStruct)-1]

		pkg := st.Obj().Pkg()
		if pkg == nil {
			return structName
		}

		return fmt.Sprintf("%s.%s", pkg.Name(), structName)
	}

	return stype.String()
}

func fieldParams(field *ast.Field, tinfo *types.Info) []*paramInfo {
	ftype := typeInfo{
		syntaxType:   field.Type,
		semanticType: tinfo.Types[field.Type].Type,
	}

	return lo.
		If(field.Names == nil, []*paramInfo{{field: field, name: "", ftype: &ftype}}).
		ElseF(func() []*paramInfo {
			return lo.Map(field.Names, func(name *ast.Ident, _ int) *paramInfo {
				return &paramInfo{field: field, name: name.Name, ftype: &ftype}
			})
		})
}

type recvInfo struct {
	recv *ast.Field

	recvParam  *paramInfo
	typeParams []*paramInfo
	fields     []*paramInfo
}

func processRecv(recv *ast.Field, tinfo *types.Info) *recvInfo {
	recvParam := fieldParams(recv, tinfo)[0]
	namedType := getNamedType(recvParam.ftype.semanticType)

	return &recvInfo{
		recv:      recv,
		recvParam: recvParam,
		typeParams: lo.
			If[[]*paramInfo](namedType == nil, nil).
			ElseF(func() []*paramInfo {
				var typeParams = make([]*paramInfo, 0, namedType.TypeParams().Len())
				for i := 0; i < namedType.TypeParams().Len(); i++ {
					typeParams = append(typeParams, &paramInfo{
						field: recv,
						name:  namedType.TypeParams().At(i).Obj().Name(),
						ftype: &typeInfo{
							syntaxType:   nil, // skipping due to complexity and unimportance
							semanticType: namedType.TypeParams().At(i).Constraint(),
						},
					})
				}

				return typeParams
			}),
		fields: lo.
			If[[]*paramInfo](namedType == nil, nil).
			ElseF(func() []*paramInfo {
				structType, ok := namedType.Underlying().(*types.Struct)
				if !ok {
					return nil
				}

				var fields = make([]*paramInfo, 0, structType.NumFields())
				for i := 0; i < structType.NumFields(); i++ {
					fields = append(fields, &paramInfo{
						field: recv,
						name:  structType.Field(i).Name(),
						ftype: &typeInfo{
							syntaxType:   nil, // skipping due to complexity and unimportance
							semanticType: structType.Field(i).Type(),
						},
					})
				}

				return fields
			}),
	}
}

func getNamedType(tt types.Type) *types.Named {
	switch t := tt.(type) {
	case *types.Pointer:
		return getNamedType(t.Elem())
	case *types.Named:
		return t
	}

	return nil
}

type importInfo struct {
	Alias string
	Path  string
}

func (i *importInfo) Render() string {
	var sb strings.Builder
	if i.Alias != "" {
		sb.WriteString(i.Alias)
		sb.WriteString(" ")
	}

	sb.WriteString(i.Path)
	return sb.String()
}
