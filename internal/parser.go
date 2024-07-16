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
	"strings"

	"github.com/fadyat/gutify/internal/lo"
)

type Parser struct {
	flags *Flags

	inputFileSet  *token.FileSet
	outputFileSet *token.FileSet
	inputAst      *ast.File
	outputAst     *ast.File
}

func NewParser(flags *Flags) *Parser {
	return &Parser{
		flags: flags,
	}
}

func (p *Parser) GenerateMissingTests() (f *File, err error) {
	p.inputFileSet = token.NewFileSet()
	p.inputAst, err = parser.ParseFile(p.inputFileSet, p.flags.InputFile, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("parse input file: %w", err)
	}

	p.outputFileSet = token.NewFileSet()
	p.outputAst, err = parser.ParseFile(p.outputFileSet, p.flags.OutputFile, nil, parser.AllErrors)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("parse output file: %w", err)
	}

	inputFuncs := getFuncs(p.inputFileSet, p.inputAst, func(fs *token.FileSet, decl *ast.FuncDecl) *Fn {
		ff := parseFn(fs, decl)
		ff.generateFriendlyNames(ff.Args)
		ff.generateFriendlyNames(ff.Results)
		return ff
	})

	outputFuncs := getFuncs(p.outputFileSet, p.outputAst, func(fs *token.FileSet, decl *ast.FuncDecl) *Fn {
		ff := parseFn(fs, decl)
		ff.Name = strings.TrimPrefix(ff.Name, "Test_")
		return ff
	})

	testsMissingFn := lo.FilterMap(inputFuncs, func(item *Fn, _ int) (*Fn, bool) {
		return item, !lo.ContainsBy(outputFuncs, func(out *Fn) bool {
			return item.Name == out.Name
		})
	})

	if len(testsMissingFn) == 0 {
		return nil, ErrNoMissingTests
	}

	file := &File{
		Functions: testsMissingFn,
	}

	if p.outputAst == nil {
		file.PackageName = p.inputAst.Name.Name
		file.Imports = lo.Map(p.inputAst.Imports, func(imp *ast.ImportSpec, _ int) string {
			return imp.Path.Value
		})

		// appending empty string in cases when no imports exist, but
		// need to generate imports from the template
		file.Imports = append(file.Imports, "")
	}

	return file, nil
}

func getFuncs(fs *token.FileSet, f *ast.File, parser func(*token.FileSet, *ast.FuncDecl) *Fn) []*Fn {
	if f == nil {
		return nil
	}

	return lo.FilterMap(f.Decls, func(item ast.Decl, _ int) (*Fn, bool) {
		fn, ok := item.(*ast.FuncDecl)
		if !ok {
			return nil, false
		}

		return parser(fs, fn), true
	})
}

func parseFn(fs *token.FileSet, f *ast.FuncDecl) *Fn {
	var function = newFn(f.Name.Name)

	if f.Recv != nil {
		var (
			receiverType = getTypeName(fs, f.Recv.List[0].Type)
			receiverName = ""
		)

		if len(f.Recv.List[0].Names) > 0 {
			receiverName = f.Recv.List[0].Names[0].Name
		}

		function.Receiver = newIdentifier(receiverName, receiverType)
	}

	if f.Type.TypeParams != nil {
		function.Generics = lo.FlatMap(f.Type.TypeParams.List, func(typeParam *ast.Field, _ int) []*identifier {
			typeParamType := getTypeName(fs, typeParam.Type)
			return lo.Map(typeParam.Names, func(name *ast.Ident, _ int) *identifier {
				return newIdentifier(name.Name, typeParamType)
			})
		})
	}

	function.Args = lo.FlatMap(f.Type.Params.List, func(arg *ast.Field, _ int) []*identifier {
		argType := getTypeName(fs, arg.Type)
		if len(arg.Names) == 0 {
			return []*identifier{newIdentifier("", argType)}
		}

		return lo.Map(arg.Names, func(name *ast.Ident, _ int) *identifier {
			return newIdentifier(name.Name, argType)
		})
	})

	if f.Type.Results != nil {
		function.Results = lo.FlatMap(f.Type.Results.List, func(res *ast.Field, _ int) []*identifier {
			resType := getTypeName(fs, res.Type)
			if len(res.Names) == 0 {
				return []*identifier{newIdentifier("", resType)}
			}

			return lo.Map(res.Names, func(name *ast.Ident, _ int) *identifier {
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
