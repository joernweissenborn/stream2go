package stream2go

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
	sr Suscribor
	c chan Suscription
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

func (s Stream) Listen(sr Suscribor) (ss Suscription) {
	c := make(chan Suscription)
	s.addsus <- addsuscrption{sr,c}
	return <-c
}

func (s Stream) Transform(t Transformer) (ts Stream){
	return
}

func (s Stream) Where(f Filter) (fs Stream) {
	return
}
