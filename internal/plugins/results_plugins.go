package plugins

import (
	"fmt"
	"strings"

	"github.com/fadyat/ggt/internal"
	"github.com/fadyat/ggt/internal/lo"
)

// todo: also need to create the preparation plugins, which accepts the list of arguments
//  and made some preparation before the test function will be called.
//
//  - following function need to be added to the testcase structure.
//  - it accepts the list of arguments (which are determined by the plugin itself)
//  - it generates the prepare function call with required arguments.

// func test_name {
//    struct_fields ( interfaces or structs which doesn't require mocks/stubs
//   				  	+ pluggable to support concrete types for interfaces )		*
//    func_args		( args are immutable, defined inside the testcase )
//    func_results	( immutable results + pluggable results ) 						v
//
//    testcases := []struct {
//        name string ( immutable )
//        fields struct_fields ( immutable )
//        args func_args ( immutable )
//        want func_results ( immutable )
//
//        + pluggable preparation function											*
//    }{
//        {},
//    }
//
//    for _, tt := range testcases {
//        t.Run(tt.name, func(t *testing.T) {
// 		      struct_creation 	( pluggable to support preparation /
//		   						  different creation techs )	    				*
// 		      function_call     ( immutable )
// 		      check_results		( pluggable to support different checks )			v
//        })
//    }
// }

// ResultsPlugin is a subset of plugins responsible for modifying function return values
// and checking subsequent results.
type ResultsPlugin interface {

	// PatchResults changes the format of result values
	// for further custom validation logic.
	PatchResults([]*internal.Identifier) []*internal.Identifier

	// VerifyResults changes the validation logic for the results.
	// Map will contain the list of templates, which can be used for particular
	// result validation.
	VerifyResults([]*internal.Identifier, map[string][]string)
}

func WithResultsPlugins(fn *internal.Fn, plugins []ResultsPlugin) string {
	var (
		results       = fn.Results
		verifications = make(map[string][]string)
	)

	for _, plugin := range plugins {
		plugin.VerifyResults(results, verifications)
		results = plugin.PatchResults(results)
	}

	fn.Results = results
	return strings.Join(
		lo.MapToSlice(verifications, func(_ string, v []string) string {
			return strings.Join(v, "\n")
		}),
		"\n",
	)
}

// coreDefaultResultsPlugin is a default implementation of the ResultsPlugin interface, which
// will be used as a base for all other plugins.
// It doesn't change the results and doesn't provide any additional verification logic.
type coreDefaultResultsPlugin struct{}

func (c *coreDefaultResultsPlugin) PatchResults(identifiers []*internal.Identifier) []*internal.Identifier {
	return identifiers
}

func toGotSingle(v string) string { // todo: remove me
	if strings.HasPrefix(v, "want") {
		return strings.Replace(v, "want", "got", 1)
	}

	return v
}

func (c *coreDefaultResultsPlugin) VerifyResults(identifiers []*internal.Identifier, m map[string][]string) {
	for _, identifier := range identifiers {
		m[identifier.Name] = []string{fmt.Sprintf(
			"require.Equal(t, tt.want.%s, %s)",
			identifier.Name,
			toGotSingle(identifier.Name),
		)}
	}
}

// errorAssertionPlugin is a plugin, which replaces all the error type results with the
// special assertion function, which called after the function execution.
type errorAssertionPlugin struct{}

func (e *errorAssertionPlugin) PatchResults(identifiers []*internal.Identifier) []*internal.Identifier {
	for _, identifier := range identifiers {
		if identifier.Type == "error" {
			identifier.Type = "require.ErrorAssertionFunc"
		}
	}

	return identifiers
}

func (e *errorAssertionPlugin) VerifyResults(identifiers []*internal.Identifier, m map[string][]string) {
	for _, identifier := range identifiers {
		if identifier.Type == "error" {
			m[identifier.Name] = []string{fmt.Sprintf(
				"tt.want.%s(t, %s)",
				identifier.Name,
				toGotSingle(identifier.Name),
			)}
		}
	}
}

func newResultsPlugins() []ResultsPlugin {
	return []ResultsPlugin{
		&coreDefaultResultsPlugin{},
		&errorAssertionPlugin{},
	}
}
