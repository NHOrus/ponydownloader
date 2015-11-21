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
	in <- Image{}
	select {
	case tval, ok := <-out:
		if !ok {
			t.Fatal("Out channel is closed prematurely")
		}
		if !tval.isDeepEqual(Image{}) {
			t.Error("Pass through dispatcher mangles image")
		}
	default:
		t.Fatal("Dispatcher blocks")
	}
}
