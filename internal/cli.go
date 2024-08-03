package internal

import (
	"flag"
	"fmt"
	"strings"
)

type Flags struct {
	InputFile  string
	OutputFile string
}

func ParseFlags() (*Flags, error) {
	var f = &Flags{
		InputFile:  "<from-user>.go",
		OutputFile: "<from-user>_test.go",
	}

	flag.StringVar(&f.InputFile, "file", "", "input file")
	flag.StringVar(&f.OutputFile, "output", "", "output file")
	flag.Parse()

	if !strings.HasSuffix(f.InputFile, ".go") {
		return nil, fmt.Errorf("input file must have .go extension")
	}

	if f.OutputFile == "" {
		f.OutputFile = fmt.Sprintf("%s_test.go", strings.TrimSuffix(f.InputFile, ".go"))
	} else if !strings.HasSuffix(f.OutputFile, "_test.go") {
		return nil, fmt.Errorf("output file must have _test.go extension")
	}

	return f, nil
}
