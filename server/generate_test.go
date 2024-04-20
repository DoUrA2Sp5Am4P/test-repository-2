package main

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	for i := 0; i < 5; i++ {
		if len(generateID()) <= 0 {
			t.Error("Error when test")
		}
	}
}
