package play

// []
//
//	func Test_emptyFunc(t *testing.T) {
//		testcases := []struct {
//			name string
//		}{
//			{},
//		}
//
//		for _, tt := range testcases {
//			t.Run(tt.name, func(t *testing.T) {
//				emptyFunc()
//			})
//		}
//	}
func emptyFunc() {}

// []
//
//	func Test_withStringArgument(t *testing.T) {
//		type args struct {
//			s string
//		}
//
//		testcases := []struct {
//			name string
//			args args
//		}{
//			{},
//		}
//
//		for _, tt := range testcases {
//			t.Run(tt.name, func(t *testing.T) {
//				withStringArgument(tt.args.s)
//			})
//		}
//	}
func withStringArgument(s string) {}
