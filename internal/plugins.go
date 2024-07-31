package internal

type Plugin interface {
	Name() string
}

type _ struct {
	_ *Fn
	_ []Plugin
}

func _() []*identifier {
	// plugin1 -> plugin2 -> plugin3 -> plugin4
	//
	// each plugin accepts the result of the previous plugin
	// and returns the result of the next plugin
	//
	// accepts the list of arguments, which will be returned by the function
	// and returns new list of arguments.
	//
	// also for following list of arguments, the plugin will recreate the validation
	// function at the end of the test function.
	//
	// Example:
	//
	// func add(a, b int) (int, error) { ... }
	//
	// wantArgs: (int, error) ->
	//			 (int, require.ErrorAssertionFunc)
	//
	// wantCheck: (require.Equal(t, tt.want.a, got.a), require.Equal(t, tt.want.err, got.err)) ->
	//			  (require.Equal(t, tt.want.a, got.a), tt.want.err(t, got.err))

	return nil
}

// todo: also need to create the preparation plugins, which accepts the list of arguments
//  and made some preparation before the test function will be called.
//
//  - following function need to be added to the testcase structure.
//  - it accepts the list of arguments (which are determined by the plugin itself)
//  - it generates the prepare function call with required arguments.

// func test_name { ( accepts only name based on some pluggable rules )			    *
//    struct_fields ( interfaces or structs which doesn't require mocks/stubs
//   				  	+ pluggable to support concrete types for interfaces )		*
//    func_args		( args are immutable, defined inside the testcase )
//    func_results	( immutable results + pluggable results ) 						*
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
// 		      check_results		( pluggable to support different checks )			*
//        })
//    }
// }
