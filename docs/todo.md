### MAJOR: import aliases

### MAJOR: mocks support, generate automatic prepare function + expectations

> - independent of the test framework

### MAJOR: user documentation

### MAJOR: need to overwrite variables to another names, if they are conflicting with a testing one (t, tt, etc.)

### MAJOR: adding imports, if test file already exists, but new imports are needed

> - need to parse the already existing tests file
> - get imports + add new once
> - append already existing tests without modifications
> - append new tests

### MAJOR: installation guidelines, brew, go install, from binaries, etc.

### MAJOR: generate tests only for specified functions (regex, name, etc.)

### MAJOR: can generate only input generics for input arguments, output generics for output arguments

### MAJOR: tests generation verification framework

> - verify that generated tests are correct
> - using the special test cases, which have defined input and output

```go
//	func Test_letsReturnError(t *testing.T) {
//		type args struct {
//			msg string
//		}
//		type want struct {
//			wantErr require.ErrorAssertionFunc
//		}
//
//		testcases := []struct {
//			name string
//			args args
//			want want
//		}{
//			{},
//		}
//
//		for _, tt := range testcases {
//			t.Run(tt.name, func(t *testing.T) {
//				gotErr := letsReturnError(tt.args.msg)
//
//				tt.want.wantErr(t, gotErr)
//			})
//		}
//	}
func letsReturnError(msg string) error {
	return errors.New(msg)
}
```

### MAJOR: exported functions only

### MAJOR: adding new missing dependencies to a test function (like, new field was added to the struct)

### MINOR: header in the generated file, that it's generated using the tool

> - write some text, which tells that the file is generated using the tool
    > with tool version and reference to the tool
> - example: https://github.com/vektra/mockery/blob/fb63d008ca3ec7539eef0ab366cb993555c4ca80/pkg/generator.go#L433
> - if different versions of the tool are used, header should be regenerated

### MINOR: logging some general information

> - initializing the slog package
> - choose places, which are important to log
> - add logging to the places
> - example: generating the file, adding missing imports, adding new tests, some tests already exists

### MINOR: default `context.TODO()` for all contexts??

### MINOR: if it is not a struct, need to make it as an argument

```golang
type b string

func (b) String() {}
```

### MINOR: count minimal number for function by parsing return statements