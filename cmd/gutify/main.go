package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/fadyat/gutify/internal"
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

	for _, fn := range file.Functions {
		fmt.Println(fn.Name)
		for _, arg := range fn.Args {
			fmt.Println("  arg:", arg.Name, arg.Type)
		}

		for _, res := range fn.Results {
			fmt.Println("  res:", res.Name, res.Type)
		}

		if fn.Receiver != nil {
			fmt.Println("  receiver:", fn.Receiver.Name, fn.Receiver.Type)
			for _, field := range fn.Struct.Fields {
				fmt.Println("    field:", field.Name, field.Type)
			}
		}

		fmt.Println()
	}

	renderer := internal.NewRenderer(f)
	err = renderer.Render(file)
	exit(err, "render tests")

	out, err := exec.Command("gofmt", "-w", f.OutputFile).CombinedOutput()
	if err != nil {
		exit(fmt.Errorf("%s: %s", err, out), "format generated file")
	}
}
