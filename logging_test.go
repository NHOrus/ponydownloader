package main

import "testing"

func TestDebracketEmpty(t *testing.T) {
	a := []int{}
	b := debracket(a)
	if b != "" {
		t.Error("Empty string not empty")
	}
}

func TestDebracketOne(t *testing.T) {
	a := []int{1}
	b := debracket(a)
	if b != "1" {
		t.Error("Single value debracketed wrong")
	}
}

func TestDebracketMulti(t *testing.T) {
	a := []int{1, 2, 4, 42}
	b := debracket(a)
	if b != "1, 2, 4, 42" {
		t.Error("Multiple values debracketed wrong")
	}
}
