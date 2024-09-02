package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseFlags(t *testing.T) {
	type args struct {
		args []string
	}

	type want struct {
		want    *Flags
		wantErr require.ErrorAssertionFunc
	}

	testcases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "invalid flags",
			args: args{args: []string{"--invalid"}},
			want: want{
				want: nil,
				wantErr: func(t require.TestingT, err error, _ ...interface{}) {
					require.Equal(t, ErrParseFlags, err)
				},
			},
		},
		{
			name: "invalid input file",
			args: args{args: []string{"-i", "input"}},
			want: want{
				want: nil,
				wantErr: func(t require.TestingT, err error, _ ...interface{}) {
					require.Contains(t, err.Error(), `input file "input" expected format: *.go`)
				},
			},
		},
		{
			name: "invalid output file",
			args: args{args: []string{"-i", "input.go", "-o", "output"}},
			want: want{
				want: nil,
				wantErr: func(t require.TestingT, err error, _ ...interface{}) {
					require.Contains(t, err.Error(), `output file "output" expected format: *_test.go`)
				},
			},
		},
		{
			name: "valid input and output file",
			args: args{args: []string{"input.go", "output_test.go"}},
			want: want{
				want: &Flags{
					InputFile:  "input.go",
					OutputFile: "output_test.go",
				},
				wantErr: require.NoError,
			},
		},
		{
			name: "output file autocomplete",
			args: args{args: []string{"input.go"}},
			want: want{
				want: &Flags{
					InputFile:  "input.go",
					OutputFile: "input_test.go",
				},
				wantErr: require.NoError,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ParseFlags(tt.args.args)
			require.Equal(t, tt.want.want, got)
			tt.want.wantErr(t, gotErr)
		})
	}
}
