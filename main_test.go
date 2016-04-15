package main

import (
	"syscall"
	"testing"
)

func TestInterruptThrough(t *testing.T) {
	in := make(ImageCh)
	out := in.interrupt()

	in <- Image{Score: 1}
	select {
	case tval, ok := <-out:
		if !ok {
			t.Fatal("Out channel is closed prematurely")
		}
		if (tval != Image{Score: 1}) {
			t.Error("Pass through interrupter mangles image")
		}
	default:
		t.Fatal("Interrupter blocks")
	}
}

func TestInterruptSignal(t *testing.T) {
	in := make(ImageCh)
	out := in.interrupt()

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	_, ok := <-out
	if ok {
		t.Error("Channel open and passes data when it should be closed by interrupt")
	}
}

func TestInterruptClose(t *testing.T) {
	in := make(ImageCh)
	out := in.interrupt()

	close(in)
	_, ok := <-out
	if ok {
		t.Error("Channel open and passes data when it should be closed by end of input")
	}
}
