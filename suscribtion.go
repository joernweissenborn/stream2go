package stream2go

func NewSuscription(index int, close chan int, sr Suscribor) (s Suscription) {
	s.in = make(chan interface {})
	s.add = make(streamchannel)
	go s.add.pipe(s.in)
	s.index = index
	s.close = close
	s.sr = sr
	go s.run()
	return
}
type Suscription struct {

	in chan interface {}
	add streamchannel
	index int
	close chan int

	sr Suscribor
}

func (s Suscription) run(){
	for d := range s.in {
		s.sr(d)
	}
}


func (s Suscription) Close() {
	s.close <- s.index
}



