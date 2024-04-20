package main

import (
	"testing"
)

func TestGetConcurrency(t *testing.T) {
	if getConcurrency() <= 0 {
		t.Error("Concurrency error")
		t.Error("Concurrency =", getConcurrency())
	}
}

func TestAll(t *testing.T) {
	t.Run("TestGetConcurrency", TestGetConcurrency)
}
