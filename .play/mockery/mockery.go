//go:generate ggt -input mockery.go

package mockery

import (
	"fmt"
	astalias "go/ast"
)

//go:generate mockery --name=InPackageTools --inpackage --case=underscore --with-expecter
type InPackageTools interface{ fmt.Stringer }

//go:generate mockery --name=SeparatePackageTools --case=underscore --with-expecter --output=mocks
type SeparatePackageTools interface{ fmt.Stringer }

type someToolsManager struct {
	inPackageTools       InPackageTools
	separatePackageTools SeparatePackageTools
}

func (m *someToolsManager) UseInPackageTools(_ astalias.File) string {
	return m.inPackageTools.String()
}

func (m *someToolsManager) UseSeparatePackageTools() string {
	return m.separatePackageTools.String()
}
