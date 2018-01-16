package main

import "testing"

func TestFilterInit(t *testing.T) {
	filterInit(&FiltOpts{ScoreF: true, FavesF: true}, false)
	if len(filters) != 2 {
		t.Error("Filter initialization doesn't work as expected")
	}
}

func TestFilterNone(t *testing.T) {
	in := make(chan Image, 1)
	out := FilterChannel(in)
	in <- Image{}
	_, ok := <-out
	if !ok {
		t.Fatal("No-op filter closed unexpectedly")
	}
}

func TestFilterAlwaysTrue(t *testing.T) {
	in := make(chan Image, 1)
	filter := filterGenerator(func(Image) bool { return true }, false)
	out := filter(in)
	in <- Image{}
	_, ok := <-out

	if !ok {
		t.Fatal("Filter closed unexpectedly")
	}
	close(in)
	_, ok = <-out
	if ok {
		t.Error("Filter remains unexpectedly open")
	}
}

func TestFilterAlwaysFalse(t *testing.T) {
	in := make(chan Image, 1)
	filter := filterGenerator(func(Image) bool { return false }, false)
	out := filter(in)
	in <- Image{}

	close(in)
	_, ok := <-out

	if ok {
		t.Error("Filter remains unexpectedly open")
	}
}

func TestFilterComplex(t *testing.T) {
	filterInit(&FiltOpts{ScoreF: true, FavesF: true}, false)

	in := make(chan Image, 3)
	in <- Image{Score: -1}
	in <- Image{Faves: -1}
	in <- Image{Imgid: "1"}

	out := FilterChannel(in)
	close(in)
	pass := <-out

	if (pass != Image{Imgid: "1"}) {
		t.Error("Incorrect work of the filter, passed ", pass, "instead of ", Image{Imgid: "1"})
	}

	pass, ok := <-out

	if ok && (pass != Image{}) {
		t.Error("Wat is going on with FilterChannel?")
	}
}
