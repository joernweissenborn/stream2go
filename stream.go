package stream2go

import "github.com/joernweissenborn/future2go"

func NewStream() (s Stream) {
	s.in = make(chan interface {})
	s.addsus = make(chan addsuscrption)
	s.rmsusc = make(chan int)
	go s.run()
	return
}

type Stream struct {

	in chan interface {}

	addsus chan addsuscrption
	rmsusc chan int

}

type addsuscrption struct {
	sr Suscriber
	c chan Suscription
}

func(s Stream) Close(){
	close(s.addsus)
	close(s.rmsusc)
	close(s.in)
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
}

func (s Stream) add(d interface {}) {
	s.in <- d
	return
}

func (s Stream) Listen(sr Suscriber) (ss Suscription) {
	c := make(chan Suscription)
	s.addsus <- addsuscrption{sr,c}
	return <-c
}

func (s Stream) Transform(t Transformer) (ts Stream){
	ts = NewStream()
	s.Listen(func(d interface {}){
			ts.add(t(d))
	})
	return
}

func (s Stream) Where(f Filter) (fs Stream) {
	fs = NewStream()
	s.Listen(func(d interface {}){
		if f(d) {
			fs.add(d)
		}
	})

	return
}

func (s Stream) First() (f *future2go.Future) {
	f = future2go.New()
	s.Listen(func(d interface {}){f.Complete(d)})
	return
}
