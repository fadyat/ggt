//go:generate ggt -input mockery.go

package mockery

import (
	"fmt"
)

//go:generate mockery --name=InPackageTools --inpackage --case=underscore --with-expecter
type InPackageTools interface{ fmt.Stringer }

//go:generate mockery --name=SeparatePackageTools --case=underscore --with-expecter --output=mocks
type SeparatePackageTools interface{ fmt.Stringer }

type someToolsManager struct {
	inPackageTools       InPackageTools
	separatePackageTools SeparatePackageTools
}

func (m *someToolsManager) UseInPackageTools() string {
	return m.inPackageTools.String()
}

func (m *someToolsManager) UseSeparatePackageTools() string {
	return m.separatePackageTools.String()
}
