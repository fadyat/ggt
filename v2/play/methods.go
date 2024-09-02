package play

import (
	"fmt"
	"go/ast"
)

type ggt struct{ a, b int }

func (g *ggt) PublicMethod() {}

func (g *ggt) privateMethod() {}

type ggtInterface struct{ stringer fmt.Stringer }

func (g *ggtInterface) InterfaceMethod() string { return g.stringer.String() }

type ggtNested struct{ nested struct{ a int } }

func (g *ggtNested) NestedMethod() {}

type ggtEmbeddedOuterPackage struct{ ast.Field }

func (g *ggtEmbeddedOuterPackage) EmbeddedMethodOuter() {}

type ggtEmbedded struct{ *ggt }

func (g *ggtEmbedded) EmbeddedMethod() {}

type stringTypeUnderlying string

func (stringTypeUnderlying) UnderlyingMethod() string { return "" }

type ggtGeneric[T any] struct{ t T }

func (g *ggtGeneric[T]) GenericMethod() {}

type ggtNoPointer struct{ a int }

func (g ggtNoPointer) NoPointerMethod() {}

type ggtMultipleGenerics[T1 comparable, T2 fmt.Stringer] struct {
	t1 T1
	t2 T2
	ch chan<- T1
	mp map[T1]T2
	f  ast.Field
}

func (g ggtMultipleGenerics[T1, T2]) MultipleGenericsMethod() {}
