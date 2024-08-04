package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/fadyat/ggt/internal"
	"github.com/fadyat/ggt/internal/plugins"
	"github.com/fadyat/ggt/internal/renderer"
)

func exit(err error, msg string) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", msg, err)
		os.Exit(1)
	}
}

func main() {
	f, err := internal.ParseFlags()
	exit(err, "parse flags")

	parser := internal.NewParser(f)
	file, err := parser.GenerateMissingTests()
	if err != nil {
		if errors.Is(err, internal.ErrNoMissingTests) {
			fmt.Println("no missing tests")
			return
		}

		exit(err, "generate tests")
	}

	r := renderer.NewRenderer(f)
	pf := plugins.NewPluggableFile(file)

	// showTree(pf) // todo: remove me

	err = r.Render(pf)
	exit(err, "render tests")

	out, err := exec.Command("gofmt", "-w", f.OutputFile).CombinedOutput()
	if err != nil {
		exit(fmt.Errorf("%s: %s", err, out), "format generated file")
	}
}

func showTree(f *plugins.PluggableFile) {
	fmt.Println(f.PackageName)
	for _, fn := range f.Functions {
		fmt.Printf("\t%s\n", fn.TestName())

		if fn.Struct != nil {
			fmt.Printf("\t\t%s\n", fn.Struct.Name)
			for _, field := range fn.Struct.Fields {
				fmt.Printf("\t\t\t%s %s\n", field.Name, field.Type)
			}
		}

		fmt.Printf("\t\t%s\n", fn.Name)
		for _, arg := range fn.Args {
			fmt.Printf("\t\t\t%s %s\n", arg.Name, arg.Type)
		}

		for _, res := range fn.Results {
			fmt.Printf("\t\t\t%s %s\n", res.Name, res.Type)
		}

		fmt.Println()
	}
}
