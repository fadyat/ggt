package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jessevdk/go-flags"
)

var ErrParseFlags = errors.New("error parsing flags")

type Flags struct {
	InputFile  string `short:"i" long:"input" description:"input file"`
	OutputFile string `short:"o" long:"output" description:"output file"`
	Debug      bool   `short:"d" long:"debug" description:"debug mode"`
}

// mergeWithRemainingArgs writes the argument value, if not already set, from the
// flag to the corresponding field in the Flags struct.
func (f *Flags) mergeWithRemainingArgs(args []string) {
	if len(args) > 0 && f.InputFile == "" {
		f.InputFile = args[0]
	}

	if f.OutputFile == "" {
		if len(args) > 1 {
			f.OutputFile = args[1]
		} else {
			f.OutputFile = fmt.Sprintf("%s_test.go", strings.TrimSuffix(f.InputFile, ".go"))
		}
	}
}

func (f *Flags) validate() error {
	if !strings.HasSuffix(f.InputFile, ".go") {
		return fmt.Errorf("input file %q expected format: *.go", f.InputFile)
	}

	if !strings.HasSuffix(f.OutputFile, "_test.go") {
		return fmt.Errorf("output file %q expected format: *_test.go", f.OutputFile)
	}

	return nil
}

func ParseFlags(args []string) (*Flags, error) {
	var opts Flags
	args, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return nil, ErrParseFlags
	}

	opts.mergeWithRemainingArgs(args)
	if err = opts.validate(); err != nil {
		return nil, err
	}

	return &opts, nil
}
