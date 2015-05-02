package stream2go

type StreamController struct {
	Stream
}

func (sc StreamController) Add(Data interface {}) {
	sc.Stream.add(Data)
}

func New() (sc StreamController){
	sc.Stream = NewStream()
	return
}
