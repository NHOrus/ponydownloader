package main

import "testing"

func TestDebracketEmpty(t *testing.T) {
	a := []int{}
	b := debracket(a)
	if b != "" {
		t.Error("String that should be empty is ", b)
	}
}

func TestDebracketOne(t *testing.T) {
	a := []int{1}
	b := debracket(a)
	if b != "1" {
		t.Error("Single value debracketed wrong, instead of 1 got ", b)
	}
}

func TestDebracketMulti(t *testing.T) {
	a := []int{1, 2, 4, 42}
	b := debracket(a)
	if b != "1, 2, 4, 42" {
		t.Error("Multiple values debracketed wrong, instead of 1, 2, 4, 42 got ", b)
	}
}
