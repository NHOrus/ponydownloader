package main

import "testing"

func TestFilterNop(t *testing.T) {
	filterInit(&FiltOpts{
		Filter: false,
	}, false)
	in := make(ImageCh, 1)
	out := FilterChannel(in)
	in <- Image{}
	tval, ok := <-out
	if !ok {
		t.Error("No-op filter closed unexpectedly")
	}
	if !ok && tval.Imgid == 0 &&
		tval.URL == "" &&
		tval.Filename == "" &&
		tval.Score == 0 &&
		tval.Faves == 0 {
		t.Error("Pass No-op filter mangles image")
	}
}
