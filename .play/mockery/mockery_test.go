package mockery

import (
	astalias "go/ast"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fadyat/ggt/play/mockery/mocks"
)

func Test_someToolsManager_UseInPackageTools(t *testing.T) {
	type fields struct {
		inPackageTools       InPackageTools
		separatePackageTools *mocks.SeparatePackageTools
	}
	type args struct {
		want astalias.File
	}
	type want struct {
		want string
	}

	testcases := []struct {
		name   string
		fields fields
		args   args
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

			got := m.UseInPackageTools(tt.args.want)
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
