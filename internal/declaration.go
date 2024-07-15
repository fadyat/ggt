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
	Receiver *identifier
	Args     []*identifier
	Generics []*identifier
	Results  []*identifier
}

func newFn(name string) *Fn {
	return &Fn{
		Name: name,
	}
}

func (f *Fn) generateFriendlyNames(iterable []*identifier) {
	var (
		countTypes = lo.CountValuesBy(iterable, func(res *identifier) string {
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

	for _, res := range iterable {
		// setting names only for none user-defined names or empty names
		if res.Name != "" && res.Name != "_" {
			continue
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

type identifier struct {
	Name string
	Type string
}

func newIdentifier(name, typ string) *identifier {
	return &identifier{
		Name: name,
		Type: typ,
	}
}
