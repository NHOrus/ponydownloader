package main

import (
	"syscall"
	"testing"
)

func TestParseInterrupt(t *testing.T) {
	if !isParseInterrupted() {
		t.Error("Parsing gets interrupted by default")
	}

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	if isParseInterrupted() {
		t.Error("Parsing will continue after user interrupt")
	}
}
