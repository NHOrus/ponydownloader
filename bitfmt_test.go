package main

import "testing"

func TestBytefmtNegative(t *testing.T) {
	defer func() {
		if err := recover(); err != "Natural number is less than zero. Stuff is wrong" {
			t.Error("Not panicing when we really, really should")
		}
	}()
	fmtbytes(-3)
}

func TestBytefmtAll(t *testing.T) {
	a := 192.0
	if fmtbytes(a) != "192 B" {
		t.Error("Default formatting error, wanted 192 B, got ", fmtbytes(a))
	}
	a = 8630
	if fmtbytes(a) != "8.43 KiB" {
		t.Error("Kilobyte formatting error, wanted 8.43 KiB, got ", fmtbytes(a))
	}
	a = 8837120
	if fmtbytes(a) != "8.43 MiB" {
		t.Error("Kilobyte formatting error, wanted 8.43 MiB, got ", fmtbytes(a))
	}
	a = 9049210880
	if fmtbytes(a) != "8.43 GiB" {
		t.Error("Kilobyte formatting error, wanted 8.43 GiB, got ", fmtbytes(a))
	}
	a = 9049210880 * 1024
	if fmtbytes(a) != "8.43 TiB" {
		t.Error("Kilobyte formatting error, wanted 8.43 TiB, got ", fmtbytes(a))
	}
	a = 9049210880 * 1024 * 1024
	if fmtbytes(a) != "way too many B" {
		t.Error("Kilobyte formatting error, wanted 8.43 GiB, got ", fmtbytes(a))
	}
}
