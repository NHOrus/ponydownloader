package main

import (
	"os"
	"testing"
)

func TestDispatcherThrough(t *testing.T) {
	in := make(ImageCh)
	out := make(ImageCh)
	sig := make(chan os.Signal)

	go in.dispatcher(sig, out)
	in <- Image{Score: 1}
	select {
	case tval, ok := <-out:
		if !ok {
			t.Fatal("Out channel is closed prematurely")
		}
		if (tval != Image{Score: 1}) {
			t.Error("Pass through dispatcher mangles image")
		}
	default:
		t.Fatal("Dispatcher blocks")
	}
}

func TestDispatcherClose(t *testing.T) {
	in := make(ImageCh)
	out := make(ImageCh)
	sig := make(chan os.Signal)

	go in.dispatcher(sig, out)

	sig <- os.Interrupt
	_, ok := <-out
	if ok {
		t.Error("Channel open and passes data when it should be closed")
	}
}
