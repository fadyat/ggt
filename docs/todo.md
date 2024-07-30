### MAJOR: generate want inside the testcase struct, not in want struct

### MAJOR: error assertion `wantErr assert.ErrorFunc / require.ErrorFunc`

> - via flags?
> - how to make it more flexible? like plugging for custom logics

### MAJOR: mocks support, generate automatic prepare function + expectations, independent of the test framework

### MAJOR: user documentation

### MAJOR: adding imports, if test file already exists, but new imports are needed

> - need to parse the already existing tests file
> - get imports + add new once
> - append already existing tests without modifications
> - append new tests

### MAJOR: installation guidelines, brew, go install, from binaries, etc.

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

    ```text
    type b string
    
    func (b) String() {}
    ```
