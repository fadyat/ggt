package mockery

import (
	"github.com/fadyat/ggt/play/mockery/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_someToolsManager_UseInPackageTools(t *testing.T) {
	type fields struct {
		inPackageTools       InPackageTools
		separatePackageTools *mocks.SeparatePackageTools
	}
	type want struct {
		want string
	}

	testcases := []struct {
		name   string
		fields fields
		want   want
	}{
		{},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			m := someToolsManager{
				inPackageTools:       tt.fields.inPackageTools,
				separatePackageTools: tt.fields.separatePackageTools,
			}

			got := m.UseInPackageTools()
			require.Equal(t, tt.want.want, got)
		})
	}
}

func Test_someToolsManager_UseSeparatePackageTools(t *testing.T) {
	type fields struct {
		inPackageTools       InPackageTools
		separatePackageTools *mocks.SeparatePackageTools
	}
	type want struct {
		want string
	}

	testcases := []struct {
		name   string
		fields fields
		want   want
	}{
		{},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			m := someToolsManager{
				inPackageTools:       tt.fields.inPackageTools,
				separatePackageTools: tt.fields.separatePackageTools,
			}

			got := m.UseSeparatePackageTools()
			require.Equal(t, tt.want.want, got)
		})
	}
}
