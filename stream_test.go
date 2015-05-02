package stream2go

import (
	"testing"
	"github.com/joernweissenborn/future2go"
)

func TestStreamBasics(t *testing.T){
	var sc StreamController
	sc = New()
	defer sc.Close()
	c := make(chan interface {})
	sc.Listen(testlistener(c))
	sc.Add("test")
	if (<-c).(string) != "test" {
		t.Error("got wrong data")
	}
}

func TestStreamFirst(t *testing.T){
	var sc StreamController
	sc = New()
	defer sc.Close()
	c := make(chan interface {})
	sc.First().Then(testcompleter(c))
	sc.Add("test")
	if (<-c).(string) != "test" {
		t.Error("got wrong data")
	}
}


func TestStreamFilter(t *testing.T){
	var sc StreamController
	sc = New()
	defer sc.Close()
	c := make(chan interface {})
	sc.Where(func(d interface {})bool {return d.(int) != 2}).Listen(testlistener(c))
	sc.Add(1)
	sc.Add(2)
	sc.Add(2)
	sc.Add(1)
	sc.Add(2)
	sc.Add(5)
	for i := 0; i < 3; i++ {
	if (<-c).(int) == 2 {
		t.Error("got 2")
	}
	}
}
func TestStreamTransformer(t *testing.T){
	var sc StreamController
	sc = New()
	defer sc.Close()
	c := make(chan interface {})
	sc.Transform(func(d interface {})interface {}{return d.(int)*2}).Listen(testlistener(c))
	sc.Add(5)
	if (<-c).(int) != 10 {
		t.Error("got wrong data")
	}
}

func TestStreamMultiplex(t *testing.T){
	var sc StreamController
	sc = New()
	defer sc.Close()
	c1 := make(chan interface {})
	sc.Listen(testlistener(c1))
	c2 := make(chan interface {})
	sc.Listen(testlistener(c2))
	sc.Add("test")
	if (<-c1).(string) != "test" {
		t.Error("got wrong data")
	}
	if (<-c2).(string) != "test" {
		t.Error("got wrong data")
	}
}

func testlistener(c chan interface {}) Suscriber {
	return func(d interface {}){
		c<-d
	}
}

func testcompleter(c chan interface {}) future2go.CompletionFunc {
	return func(d interface {})interface {}{
		c<-d
		return nil
	}
}
