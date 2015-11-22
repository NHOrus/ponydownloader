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
		if tval.Imgid == 0 &&
		tval.URL == "" &&
		tval.Filename == "" &&
		tval.Score == 0 &&
		tval.Faves == 0 {
		} else {
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

	close(in)
	_, ok := <- out
	if ok {
		t.Error("Channel open and passes data when it should be closed. Or blocking")
	}
}