package fuzz

import (
	"testing"
)

func TestMessageIDFuzz(t *testing.T) {
	input := []byte{0x12, 0x00, 0x32, 0x00}

	MessageIDFuzz(input)
}