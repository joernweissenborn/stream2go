package stream2go

type Suscribor func(interface {})
type Transformer func(interface {}) interface {}
type Filter func(interface {}) bool
