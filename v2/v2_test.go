package main_test

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

const (
	mode = packages.NeedName |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo

	directory = "./play/"
	pattern   = "./..."
)

// Test_ggt it's integration like testing framework.
// It's parse the files from .play directory and generates the test cases.
// Each function at .play directory have comment with expected output.
// The expected output will be used to compare with the generated test function.
func Test_ggt(t *testing.T) {
	cfg := &packages.Config{Fset: token.NewFileSet(), Mode: mode, Dir: directory}
	pkgs, err := packages.Load(cfg, pattern)
	require.NoError(t, err)

	for _, pkg := range pkgs {
		for _, fileAst := range pkg.Syntax {
			ast.Inspect(fileAst, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					t.Run(fmt.Sprintf("%s/%s", pkg.Name, funcDecl.Name.Name), func(t *testing.T) {
						args, expected := parseFuncDoc(funcDecl)

						// todo: need to call the implementation with generation of test cases.
						//  and pass args as os.Args to create flags.
						got := ""
						_ = args

						assert.Equal(t, expected, got, args)
					})
				}

				return true
			})
		}
	}
}

// parseFuncDoc parse the function documentation and returns the arguments and the expected output.
// arguments are written in the first line of the comment, \n, and the expected output.
// also here we removes the first \t from the comment, because it's using for valid comment.
func parseFuncDoc(fd *ast.FuncDecl) ([]string, string) {
	split := strings.Split(fd.Doc.Text(), "\n")
	if len(split) == 0 {
		return []string{}, ""
	}

	args := strings.Trim(split[0], "[]")
	lines := make([]string, 0, len(split)-1)

	for _, line := range split[1:] {
		line = strings.TrimPrefix(line, "\t")
		if line != "" {
			lines = append(lines, line)
		}
	}

	return strings.Split(args, ", "), strings.Join(lines, "\n")
}
