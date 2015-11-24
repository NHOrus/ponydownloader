package main

import "testing"

func TestFilterNop(t *testing.T) {
	filterInit(&FiltOpts{}, false)
	in := make(ImageCh, 1)
	out := FilterChannel(in)
	in <- Image{}
	tval, ok := <-out
	if !ok {
		t.Fatal("No-op filter closed unexpectedly")
	}
	if tval.Imgid != 0 ||
		tval.URL != "" ||
		tval.Filename != "" ||
		tval.Score != 0 ||
		tval.Faves != 0 {
		t.Error("No-op filter mangles passing image")
	}
}

func TestFilterGeneratedTrue(t *testing.T) {
	in := make(ImageCh, 1)
	filter := filterGenerator(func(Image) bool { return true }, false)
	out := filter(in)
	in <- Image{}
	tval, ok := <-out

	if !ok {
		t.Fatal("Generated filter closed unexpectedly")
	}
	if tval.Imgid != 0 ||
		tval.URL != "" ||
		tval.Filename != "" ||
		tval.Score != 0 ||
		tval.Faves != 0 {
		t.Error("Generated filter mangles passing image")
	}
	close(in)
	_, ok = <-out
	if ok {
		t.Error("Generated filter remains unexpectedly open")
	}
}

func TestFilterGeneratedFalse(t *testing.T) {
	in := make(ImageCh, 1)
	filter := filterGenerator(func(Image) bool { return false }, false)
	out := filter(in)
	in <- Image{}

	close(in)
	_, ok := <-out
	if ok {
		t.Error("Generated filter remains unexpectedly open")
	}
}
