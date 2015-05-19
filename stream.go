package stream2go

import "github.com/joernweissenborn/future2go"

func NewStream() (s Stream) {
	s.in = make(chan interface {})
	s.addsus = make(chan addsuscrption)
	s.rmsusc = make(chan int)
	s.Closed = future2go.New()
	go s.run()
	return
}

type Stream struct {

	in chan interface {}

	addsus chan addsuscrption
	rmsusc chan int

	Closed *future2go.Future
}

type addsuscrption struct {
	sr Suscriber
	c chan Suscription
}

func(s Stream) Close(){
	if s.Closed.IsComplete() {
		return
	}
	close(s.addsus)
	close(s.rmsusc)
	close(s.in)
	s.Closed.Complete(nil)
}

func (s Stream) run() {

	suscribtions := map[int]Suscription{}
	nextSusIndex := 0

	ok := true

	for ok {
		select {

		case sus, ok := <-s.addsus:
			if !ok  {return}
			sc := NewSuscription(nextSusIndex, s.rmsusc,sus.sr)
			suscribtions[nextSusIndex] = sc
		sus.c <- sc
			nextSusIndex++

		case index, ok := <-s.rmsusc:
			if !ok  {return}
			delete(suscribtions,index)

		case d, ok := <- s.in:
			if !ok  {return}
		for _, sc := range suscribtions {
			sc.add <- d
		}
		}
	}
	for _, sus := range suscribtions {
		sus.Close()
	}
}

func (s Stream) add(d interface {}) {
	if s.Closed == nil {
		panic("Add on noninitialized Stream")
	}
	if !s.Closed.IsComplete() {
		s.in <- d
	}
	return
}

func (s Stream) Listen(sr Suscriber) (ss Suscription) {
	c := make(chan Suscription)
	s.addsus <- addsuscrption{sr,c}
	return <-c
}

func (s Stream) Transform(t Transformer) (ts Stream){
	ts = NewStream()
	s.Closed.Then(func(d interface {})interface {}{
		ts.Close()
		return nil
	})
	s.Listen(func(d interface {}){
		ts.add(t(d))
	})
	return
}

func (s Stream) Where(f Filter) (fs Stream) {
	fs = NewStream()
	s.Closed.Then(func(interface {})interface{}{
		fs.Close()
		return nil
	})
	s.Listen(func(d interface {}){
		if f(d) {
			fs.add(d)
		}
	})

	return
}

func (s Stream) WhereNot(f Filter) (fs Stream) {
	fs = NewStream()
	s.Closed.Then(func(interface {})interface{}{
		fs.Close()
		return nil
	})
	s.Listen(func(d interface {}){
		if !f(d) {
			fs.add(d)
		}
	})

	return
}

func (s Stream) First() (f *future2go.Future) {
	f = future2go.New()
	f2 := future2go.New()
	f2.Complete(s.Listen(func(d interface {}){
		f.Complete(d)
		f2.Then(func(d interface {})interface {}{
			d.(Suscription).Close()
			return nil
		})
	}))

	return
}

func (s Stream) Split(f Filter) (y Stream, n Stream) {
	return s.Where(f), s.WhereNot(f)
}

