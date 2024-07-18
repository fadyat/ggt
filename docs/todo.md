- adding imports, if test file already exists, but new imports are needed
- header in the generated file, that it's generated using the tool
- logging some general information
- user documentation
- if it is not a struct, need to make it as an argument
    ```text
    type b string
    
    func (b) String() {}
    ```

- mocks support, generate automatic prepare function + expectations, independent of the test framework
- error assertion `wantErr assert.ErrorFunc / require.ErrorFunc`
- generate want inside the testcase struct, not in want struct
- default context.TODO() for all contexts??
- brew?
