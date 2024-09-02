package observability

import (
	"fmt"

	"github.com/fadyat/ggt/internal/plugins"
)

func ShowTree(f *plugins.PluggableFile) {
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
