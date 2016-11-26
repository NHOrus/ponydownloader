package main

import (
	"syscall"
	"testing"
	"time"
)

func TestInitInterrupter(t *testing.T) {
	if isInterrupted() {
		t.Error("Parsing gets interrupted by default")
	}

	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	if err != nil {
		t.Skip("Can't get pid, skipping")
	}

	if isInterrupted() {
		t.Error("Parsing will continue after user interrupt")
	}
}

func TestInterruptThrough(t *testing.T) {
	in := make(ImageCh)
	interrupter = make(chan struct{}, 1)
	out := in.interrupt()

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	select {
	case in <- Image{Score: 1}:
	case <-timeout:
		t.Fatal("Can't image into channel")
	}

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

func TestInterruptSequence(t *testing.T) {
	in := make(ImageCh)
	interrupter = make(chan struct{}, 1)
	out := in.interrupt()

	in <- Image{Score: 1}
	<-out

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	select {
	case in <- Image{Faves: 1}:
	case <-timeout:
		t.Fatal("Can't push second image into channel")
	}

	select {
	case tval, ok := <-out:
		if !ok {
			t.Fatal("Out channel is closed prematurely")
		}
		if (tval != Image{Faves: 1}) {
			t.Error("Second pass through interrupter mangles image")
		}
	default:
		t.Fatal("Interrupter blocks on second image")
	}
}

func TestInterruptClose(t *testing.T) {
	in := make(ImageCh)
	interrupter = make(chan struct{}, 1)
	out := in.interrupt()

	close(in)
	_, ok := <-out
	if ok {
		t.Error("Channel open and passes data when it should be closed by end of input")
	}
}
