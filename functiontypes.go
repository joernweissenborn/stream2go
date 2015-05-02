package stream2go

type Suscriber func(interface {})
type Transformer func(interface {}) interface {}
type Filter func(interface {}) bool
