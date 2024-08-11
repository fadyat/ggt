package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
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
		packages.NeedTypesInfo

	pattern = "./.play"
)

var (
	lg   = slog.Default()
	fset = token.NewFileSet()
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	flag.Parse()

	cfg := &packages.Config{Fset: fset, Mode: mode, Dir: flag.Arg(0), Tests: false}

	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		log.Fatalf("loading packages: %v", err)
	}

	for _, pkg := range pkgs {
		processPackage(pkg)
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

func processPackage(pkg *packages.Package) *packageInfo {
	lg.Debug("processing package", slog.String("name", pkg.Name))

	pi := packageInfo{pkg: pkg, files: make(map[string]*fileInfo)}
	for _, fileAst := range pkg.Syntax {
		fi := processFile(fileAst, pkg.TypesInfo)
		pi.files[fileAst.Name.Name] = fi
	}

	return &pi
}

func processFile(fileAst *ast.File, tinfo *types.Info) *fileInfo {
	lg.Debug("processing file", slog.String("name", fileAst.Name.Name))

	fi := fileInfo{file: fileAst, funcs: make(map[string]*funcInfo)}
	for _, decl := range fileAst.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			fni := processFuncDecl(funcDecl, tinfo)
			fi.funcs[funcDecl.Name.Name] = fni

			if fni.recv != nil {
				fmt.Println(
					"recv:",
					fni.recv.recvParam.name,
					fni.recv.recvParam.ftype.Syntax(),
					"generics",
					strings.Join(lo.Map(fni.recv.typeParams, func(p *paramInfo, _ int) string {
						return p.name + " " + p.ftype.semanticType.String()
					}), ", "),
					"fields",
					strings.Join(lo.Map(fni.recv.fields, func(p *paramInfo, _ int) string {
						return p.name + " " + p.ftype.semanticType.String()
					}), ", "),
				)
			}
		}
	}

	return &fi
}

func processFuncDecl(fd *ast.FuncDecl, tinfo *types.Info) *funcInfo {
	lg.Debug("processing func", slog.String("name", fd.Name.Name))

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

func (t *typeInfo) Syntax() string {
	var b bytes.Buffer
	_ = printer.Fprint(&b, fset, t.syntaxType)
	return b.String()
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
						name:  namedType.TypeParams().At(i).Obj().Name(), // todo: want from syntax ??
						ftype: &typeInfo{
							syntaxType:   nil, // todo: skipping syntax type for now
							semanticType: namedType.TypeParams().At(i),
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
						name:  structType.Field(i).Name(), // todo: want from syntax ??
						ftype: &typeInfo{
							syntaxType:   nil, // todo: skipping syntax type for now
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
