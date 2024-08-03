package plugins

import "github.com/fadyat/ggt/internal"

type PluggableFile struct {
	PackageName string
	Imports     []string
	Functions   []*PluggableFn
}

type PluggableFn struct {
	*internal.Fn

	Verification string
}

func NewPluggableFile(f *internal.File) *PluggableFile {
	return &PluggableFile{
		PackageName: f.PackageName,
		Imports:     f.Imports,
		Functions:   newPluggableFns(f.Functions),
	}
}

func newPluggableFns(fns []*internal.Fn) []*PluggableFn {
	var (
		pluggableFns = make([]*PluggableFn, 0, len(fns))
		rplugs       = newResultPlugins()
	)

	for _, fn := range fns {
		pluggableFns = append(pluggableFns, &PluggableFn{
			Fn:           fn,
			Verification: WithResultPlugins(fn, rplugs),
		})
	}

	return pluggableFns
}
