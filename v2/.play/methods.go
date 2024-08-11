package play

import "fmt"

type ggt struct{ a, b int }

func (g *ggt) PublicMethod() {}

func (g *ggt) privateMethod() {}

type ggtInterface struct{ stringer fmt.Stringer }

func (g *ggtInterface) InterfaceMethod() string { return g.stringer.String() }

type ggtNested struct{ nested struct{ a int } }

func (g *ggtNested) NestedMethod() {}

type ggtEmbedded struct{ *ggt }

func (g *ggtEmbedded) EmbeddedMethod() {}

type stringTypeUnderlying string

func (stringTypeUnderlying) UnderlyingMethod() string { return "" }

type ggtGeneric[T any] struct{ t T }

func (g *ggtGeneric[T]) GenericMethod() {}

type ggtNoPointer struct{ a int }

func (g ggtNoPointer) NoPointerMethod() {}
