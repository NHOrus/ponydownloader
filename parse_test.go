package main

import (
	"syscall"
	"testing"
)

func TestParseInterrupt(t *testing.T) {
	if !isParseInterrupted() {
		t.Error("Parsing gets interrupted by default")
	}

	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	if err != nil {
		t.Skip("Can't get pid, skipping")
	}

	if isParseInterrupted() {
		t.Error("Parsing will continue after user interrupt")
	}
}
