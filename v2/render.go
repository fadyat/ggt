package main

import (
	"fmt"
	"go/ast"
	"go/printer"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/fadyat/ggt/v2/lo"
)

type stdoutWriter struct{}

func (s *stdoutWriter) Write(p []byte) (n int, err error) { return os.Stdout.Write(p) }
func (s *stdoutWriter) Close() error                      { return nil }

type writerCreator func(string) (io.WriteCloser, error)

func renderPackage(pkgInfo *packagePair, wcreator writerCreator) {
	if pkgInfo.normal == nil {
		return
	}

	filePairs := make(map[string]*filePair)
	for filePath, fi := range pkgInfo.normal.files {
		testFilePath := strings.TrimSuffix(filePath, ".go") + "_test.go"
		filePairs[testFilePath] = &filePair{normal: fi}
	}

	if pkgInfo.test != nil {
		for filePath, fi := range pkgInfo.test.files {
			if _, ok := filePairs[filePath]; ok {
				filePairs[filePath].test = fi
			}
		}
	}

	for filePath, fp := range filePairs {
		if err := renderFile(filePath, fp, wcreator); err != nil {
			fmt.Println(err)
		}
	}
}

type filePair struct {
	normal *fileInfo
	test   *fileInfo
}

func (f *filePair) pkgName() string {
	if f.test != nil {
		return f.test.file.Name.Name
	}

	return f.normal.file.Name.Name
}

func renderFile(filePath string, filePair *filePair, wcreator writerCreator) error {
	out, err := wcreator(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer out.Close()

	mergedImports := make(map[string]*ast.ImportSpec)
	for _, fi := range []*fileInfo{filePair.normal, filePair.test} {
		if fi == nil {
			continue
		}

		for _, imp := range fi.file.Imports {
			mergedImports[imp.Path.Value] = imp
		}
	}

	return renderTemplate(out, map[string]any{
		// "PkgName": filePair.pkgName(), // todo: need to make pluggable if want to use with _test files
		"Imports": lo.MapToSlice(mergedImports, func(_ string, v *ast.ImportSpec) *importInfo { // todo: need to to make pluggable to support additional imports for file, when generating mocks
			return &importInfo{
				Path: v.Path.Value,
				Alias: lo.
					If(v.Name == nil, "").
					ElseF(func() string { return v.Name.Name }),
			}
		}),
		//"RemainingTests": lo.
		//	If[[]*funcInfo](filePair.test == nil, nil).
		//	ElseF(func() []*funcInfo {
		//		return lo.MapToSlice(filePair.test.funcs, func(_ string, v *funcInfo) *funcInfo {
		//			return v
		//		})
		//	}),

		// remaining content is a test file data excluding the imports
		// and the package name
		"RemainingContent": func() string {
			if filePair.test == nil {
				// todo: need to render new test file
				return ""
			}

			var b strings.Builder

			newTestFile := *filePair.test.file
			//newTestFile.Imports = lo.MapToSlice(mergedImports, func(_ string, v *ast.ImportSpec) *ast.ImportSpec {
			//	return v
			//})

			_ = printer.Fprint(&b, fset, &newTestFile)
			return b.String()
		}(),
	})
}

func renderTemplate(out io.Writer, data any) error {
	t, err := template.
		New("tmpl").
		Parse(tmpl)

	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	if err = t.ExecuteTemplate(out, "tmpl", data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

const tmpl = `
{{ .RemainingContent }}
`
