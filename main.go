package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

func exitOnError(err error, msg string) {
	if err == nil {
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", msg, err)
	os.Exit(1)
}

func getTestFileName(fileName string) string {
	return fmt.Sprintf(
		"%s_test.go",
		strings.TrimSuffix(fileName, ".go"),
	)
}

const (
	tmplFile                  = "no-mocks.tmpl"
	generatedResultPrefixName = "want"
)

// idea: generate tests for all functions in the file
//
// 1. find all functions in the file (can be methods and functions)
// 2. create *_test.go file with tests for each function (need to skip tests for functions that already have tests)
// 3. write tests in specific format (can be selected from a list of formats) to the *_test.go file
//
// functional requirements:
// - generate tests for all functions in the file
// - skip functions that already have tests
// - write tests in specific format (can be selected from a list of formats)
//
// non-functional requirements:
// - mocks support (can be inside package or in a separate package (e.g. mocks))
// - table-driven tests support
// - parallel tests support
// - outer package tests support (generate tests not only for the current package, but also for the outer packages with *_test package)
func main() {
	var (
		inputFile  = "main.go"
		outputFile = getTestFileName(inputFile)

		// todo: for PoC doing all in-memory, later will need to write to the file
		_ = outputFile
	)

	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, inputFile, nil, parser.AllErrors)
	exitOnError(err, "parse file")

	// searching all functions in the file
	var functions = make([]*fn, 0)
	for _, decl := range file.Decls {
		f, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		functions = append(functions, parseFn(f))
	}

	err = renderTemplate(os.Stdout, tmplFile, map[string]any{
		// todo: parse package name from the file
		"PackageName": "main",
		// todo: parse imports from the file
		"Imports": []string{
			"fmt",
			"go/ast",
			"go/parser",
			"go/token",
			"io",
			"os",
			"strconv",
			"strings",
			"text/template",
		},
		"Functions": functions,
	})
	exitOnError(err, "render template")
}

func kekis(aboba1, aboba2 string) (string, string) {
	return aboba1, aboba2
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"collect": func(field string, slice any) []string {
			v := reflect.ValueOf(slice)
			if v.Kind() != reflect.Slice {
				panic("collect: first argument must be a slice")
			}

			var values []string
			for i := 0; i < v.Len(); i++ {
				values = append(values, v.Index(i).FieldByName(field).String())
			}

			return values
		},
		"prefix": func(prefix string, value []string) []string {
			for i, v := range value {
				value[i] = prefix + v
			}

			return value
		},
		// to_got is a function, which converts the variables names
		// from the want-like names to the got-like names
		"to_got": func(value []string) []string {
			for i, v := range value {
				if strings.HasPrefix(v, generatedResultPrefixName) {
					value[i] = strings.Replace(v, generatedResultPrefixName, "got", 1)
				}
			}

			return value
		},
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
	}
}

func renderTemplate(out io.Writer, templatePath string, data map[string]any) error {
	tmpl, err := template.
		New(templatePath).
		Funcs(funcMap()).
		ParseFiles(templatePath)

	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	if err := tmpl.ExecuteTemplate(out, templatePath, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

func parseFn(f *ast.FuncDecl) *fn {
	var function = &fn{
		Name:    f.Name.Name,
		Args:    make([]argument, 0, len(f.Type.Params.List)),
		Results: nil,
	}

	for _, arg := range f.Type.Params.List {
		argType := getTypeName(arg.Type)

		for _, name := range arg.Names {
			function.Args = append(function.Args, argument{
				Name: name.Name,
				Type: argType,
			})
		}
	}

	if f.Type.Results == nil {
		return function
	}

	function.Results = make([]result, 0, len(f.Type.Results.List))
	for i, res := range f.Type.Results.List {
		resType := getTypeName(res.Type)

		if len(res.Names) == 0 {
			// todo: for better user-experience, need to generated better names
			//  like want instead of want1, when only single arg is returned;
			//  like wantErr for errors
			//  like return type related names? (probably AI can help with that)
			function.Results = append(function.Results, result{
				Name: generatedResultPrefixName + strconv.Itoa(i+1),
				Type: resType,
			})
			continue
		}

		for _, name := range res.Names {
			function.Results = append(function.Results, result{
				Name: name.Name,
				Type: resType,
			})
		}
	}

	// todo: add type params (generics) support
	//  need to move some logic in following function
	//  to not to skip the type params, if the results are missing
	return function
}

func getTypeName(t ast.Expr) string {
	switch v := t.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.StarExpr:
		return "*" + getTypeName(v.X)
	case *ast.SelectorExpr:
		return getTypeName(v.X) + "." + v.Sel.Name
	case *ast.ArrayType:
		return "[]" + getTypeName(v.Elt)
	case *ast.MapType:
		return "map[" + getTypeName(v.Key) + "]" + getTypeName(v.Value)

		// todo: support other types, e.g. chan, func, interface, struct, etc.
	}

	panic(fmt.Sprintf("unknown type %T", t))
}

type fn struct {
	Name    string
	Args    []argument
	Results []result
}

type argument struct {
	Name string
	Type string
}

type result struct {
	Name string
	Type string
}
