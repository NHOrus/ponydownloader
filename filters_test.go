package main

import "testing"

func TestFilterNop(t *testing.T) {
	filterInit(&FiltOpts{}, false)
	in := make(ImageCh, 1)
	out := FilterChannel(in)
	in <- Image{Score: 1}
	tval, ok := <-out
	if !ok {
		t.Fatal("No-op filter closed unexpectedly")
	}
	if (tval != Image{Score: 1}) {
		t.Error("No-op filter mangles passing image")
	}
}

func TestFilterAlwaysTrue(t *testing.T) {
	in := make(ImageCh, 1)
	filter := filterGenerator(func(Image) bool { return true }, false)
	out := filter(in)
	in <- Image{Score: 1}
	tval, ok := <-out

	if !ok {
		t.Fatal("Filter closed unexpectedly")
	}
	if (tval != Image{Score: 1}) {
		t.Error("Filter mangles passing image")
	}
	close(in)
	_, ok = <-out
	if ok {
		t.Error("Filter remains unexpectedly open")
	}
}

func TestFilterAlwaysFalse(t *testing.T) {
	in := make(ImageCh, 1)
	filter := filterGenerator(func(Image) bool { return false }, false)
	out := filter(in)
	in <- Image{Score: 1}

	close(in)
	pass, ok := <-out

	if (pass == Image{Score: 1}) {
		t.Error("Filter passes through image it shouldn't")
	}

	if ok {
		t.Error("Filter remains unexpectedly open")
	}
}
