package stream2go

import (
	"log"
	"github.com/joernweissenborn/future2go"
)

type StreamController struct {
	Stream
}

func (sc StreamController) Add(Data interface {}) {
	if sc.Closed == nil {
		log.Println("Data")
		panic("Add on noninitialized StreamController")
	}
	if sc.Stream.Closed.IsComplete(){
		panic("Add on closed stream")
	}
	sc.Stream.add(Data)
}

func New() (sc StreamController){
	sc.Stream = NewStream()
	if sc.Stream.Closed == nil {
		panic("Stream Init failed")
	}
	return
}

func (sc StreamController) Join(s Stream) {
	if (s.Closed == nil){
		panic("Join noninitialized Stream")
	}
	if (sc.Closed == nil){
		panic("Join on noninitialized Streamcontroller")
	}
	ss := s.Listen(addJoined(sc))
	s.Closed.Then(closeSus(ss))
}

func closeSus(ss Suscription) future2go.CompletionFunc {
	return func(interface {})interface {}{
		ss.Close()
		return nil
	}
}
func addJoined(sc StreamController) Suscriber {
	return func(d interface {}){
		sc.Stream.add(d)
	}
}
