package stream2go

import (
	"testing"
	"github.com/joernweissenborn/future2go"
	"time"
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

func TestStreamDeliverAfterClose(t *testing.T){
	var sc StreamController
	sc = New()
	c := make(chan interface {})
	sc.Listen(testlistener(c))
	sc.Add("test")
	sc.Close()

	select {
	case <-time.After(1 * time.Second):
		t.Error("no response")
	case data := <-c:
		if data.(string) != "test" {
			t.Error("got wrong data")
		}
	}
}

func TestStreamClose(t *testing.T){
	var sc StreamController
	sc = New()
	sc.Close()
	if !sc.Closed.IsComplete() {
		t.Error("channel didnt close")
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

func TestStreamSplit(t *testing.T){
	var sc StreamController
	sc = New()
	defer sc.Close()
	c1 := make(chan interface {})
	c2 := make(chan interface {})
	y, n := sc.Split(func(d interface {})bool {return d.(int) != 2})
	y.Listen(testlistener(c1))
	n.Listen(testlistener(c2))
	sc.Add(1)
	sc.Add(2)
	sc.Add(2)
	sc.Add(1)
	sc.Add(2)
	sc.Add(5)
	for i := 0; i < 3; i++ {
		if (<-c1).(int) == 2 {
			t.Error("got 2")
		}
	}
	for i := 0; i < 3; i++ {
		if (<-c2).(int) != 2 {
			t.Error("didnt got 2")
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
