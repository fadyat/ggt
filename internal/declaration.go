package internal

import (
	"fmt"

	"github.com/fadyat/gutify/internal/lo"
)

type File struct {
	PackageName string
	Imports     []string
	Functions   []*Fn
}

type Fn struct {
	Name     string
	Args     []*argument
	Generics []*argument
	Results  []*result
}

func newFn(name string) *Fn {
	return &Fn{
		Name: name,
	}
}

func (f *Fn) generateFriendlyResultNames() {
	var (
		countTypes = lo.CountValuesBy(f.Results, func(res *result) string {
			if res.Type == "error" {
				return res.Type
			}

			return "not_error"
		})

		getName = func(initialCount, currentCount *int, prefix string) string {
			name := prefix
			if *initialCount > 1 {
				name = fmt.Sprintf("%s%d", prefix, *initialCount-*currentCount+1)
			}

			*currentCount--
			return name
		}

		initialErrorCount                       = countTypes["error"]
		initialNotErrorCount                    = countTypes["not_error"]
		currentErrorCount, currentNotErrorCount = initialErrorCount, initialNotErrorCount
	)

	for _, res := range f.Results {

		// if any of the results already has a name, it means,
		// that the names are named by the user, we should not change them
		if res.Name != "" {
			return
		}

		var name string
		if res.Type == "error" {
			name = getName(&initialErrorCount, &currentErrorCount, "wantErr")
		} else {
			name = getName(&initialNotErrorCount, &currentNotErrorCount, "want")
		}

		res.Name = name
	}
}

type argument struct {
	Name string
	Type string
}

func newArgument(name, typ string) *argument {
	return &argument{
		Name: name,
		Type: typ,
	}
}

type result struct {
	Name string
	Type string
}

func newResult(name, typ string) *result {
	return &result{
		Name: name,
		Type: typ,
	}
}
