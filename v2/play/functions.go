package play

import "time"

// []
//
//	func Test_withNoArgument(t *testing.T) {
//		testcases := []struct {
//			name string
//		}{
//			{},
//		}
//
//		for _, tt := range testcases {
//			t.Run(tt.name, func(t *testing.T) {
//				withNoArgument()
//			})
//		}
//	}
func withNoArgument() {}

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

// []
//
//	func Test_withPointerArgument(t *testing.T) {
//		type args struct {
//			s *string
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
//				withPointerArgument(tt.args.s)
//			})
//		}
//	}
func withPointerArgument(s *string) {}

// []
//
//	func Test_withSliceArgument(t *testing.T) {
//		type args struct {
//			s []string
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
//				withSliceArgument(tt.args.s)
//			})
//		}
//	}
func withSliceArgument(s []string) {}

// []
//
//	func Test_withVariadicArgument(t *testing.T) {
//		type args struct {
//			s []string
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
//				withVariadicArgument(tt.args.s...)
//			})
//		}
//	}
func withVariadicArgument(s ...string) {}

// []
//
//	func Test_withStructArgument(t *testing.T) {
//		type args struct {
//			s struct{ Name string }
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
//				withStructArgument(tt.args.s)
//			})
//		}
//	}
func withStructArgument(s struct{ Name string }) {}

// []
//
//	func Test_withInterfaceArgument(t *testing.T) {
//		type args struct {
//			s interface{ Name() string }
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
//				withInterfaceArgument(tt.args.s)
//			})
//		}
//	}
func withInterfaceArgument(s interface{ Name() string }) {}

// []
//
//	func Test_withFunctionArgument(t *testing.T) {
//		type args struct {
//			s func() string
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
//				withFunctionArgument(tt.args.s)
//			})
//		}
//	}
func withFunctionArgument(s func() string) {}

// []
//
//	func Test_withChannelArgument(t *testing.T) {
//		type args struct {
//			s chan<- string
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
//				withChannelArgument(tt.args.s)
//			})
//		}
//	}
func withChannelArgument(s chan<- string) {}

// []
//
//	func Test_withMapArgument(t *testing.T) {
//		type args struct {
//			s map[string]string
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
//				withMapArgument(tt.args.s)
//			})
//		}
//	}
func withMapArgument(s map[string]string) {}

// []
//
//	func Test_withPointerToSliceArgument(t *testing.T) {
//		type args struct {
//			s *[]string
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
//				withPointerToSliceArgument(tt.args.s)
//			})
//		}
//	}
func withPointerToSliceArgument(s *[]string) {}

// []
//
//	func Test_withStdLibStructArgument(t *testing.T) {
//		type args struct {
//			s time.Time
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
//				withStdLibStructArgument(tt.args.s)
//			})
//		}
//	}
func withStdLibStructArgument(s time.Time) {}

// []
//
//	func Test_withStdLibInterfaceArgument(t *testing.T) {
//		type args struct {
//			s error
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
//				withStdLibInterfaceArgument(tt.args.s)
//			})
//		}
//	}
func withStdLibInterfaceArgument(s error) {}

// []
//
//	func Test_withStdLibInterfaceArgument(t *testing.T) {
//		type args struct {
//			s string
//			t string
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
//				withMultipleSameTypeArguments(tt.args.s, tt.args.t)
//			})
//		}
//	}
func withMultipleSameTypeArguments(s, t string) {}

// []
//
//	func Test_withMultipleDifferentTypeArguments(t *testing.T) {
//		type args struct {
//			s string
//			t int
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
//				withMultipleDifferentTypeArguments(tt.args.s, tt.args.t)
//			})
//		}
//	}
func withMultipleDifferentTypeArguments(s string, t int) {}

type stringTypeAlias = string

// []
//
//	func Test_withAliasArgument(t *testing.T) {
//		type args struct {
//			s stringTypeAlias
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
//				withAliasArgument(tt.args.s)
//			})
//		}
//	}
func withAliasArgument(s stringTypeAlias) {}

// []
//
//	func Test_withNoArgumentName(t *testing.T) {
//		type args struct {
//			in1 int
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
//				withNoArgumentName(tt.args.in1)
//			})
//		}
//	}
func withNoArgumentName(int) {}

// []
//
//	func Test_withSkipArgument(t *testing.T) {
//		type args struct {
//			in1 bool
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
//				withSkipArgument(tt.args.in1)
//			})
//		}
//	}
func withSkipArgument(_ bool) {}

// []
//
//	func Test_withSkippedAndNonSkippedArguments(t *testing.T) {
//		type args struct {
//			s 	bool
//			in1 string
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
//				withSkippedAndNonSkippedArguments(tt.args.s, tt.args.in1)
//			})
//		}
//	}
func withSkippedAndNonSkippedArguments(s bool, _ string) {}

func withGenericArgument[T any](s T) {}

func withMultipleGenericArguments[T any, U any](s T, t U) {}

func withReturn() string { return "" }

func withMultipleReturns() (string, int) { return "", 0 }

func withNamedReturns() (s string, i int) { return }

func withGenericReturn[T any]() *T { return nil }
