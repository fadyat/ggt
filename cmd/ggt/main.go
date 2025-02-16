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
	err = r.Render(plugins.NewPluggableFile(file))
	exit(err, "render tests")

	out, err := exec.Command("gofmt", "-w", f.OutputFile).CombinedOutput()
	if err != nil {
		exit(fmt.Errorf("%s: %s", err, out), "format generated file")
	}
}
