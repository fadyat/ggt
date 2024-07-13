package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseFlags(t *testing.T) {
	type want struct {
		want    *Flags
		wantErr error
	}

	testcases := []struct {
		name string
		want want
	}{
		{
			name: "failed_to_parse_flags",
			want: want{
				want:    nil,
				wantErr: fmt.Errorf("input file must have .go extension"),
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ParseFlags()

			require.Equal(t, tt.want.want, got)
			require.Equal(t, tt.want.wantErr, gotErr)
		})
	}
}
