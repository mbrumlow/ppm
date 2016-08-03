package ppm

import (
	"os"
	"testing"
)

// Test PPM image decoding with out comments in the header.
func TestDecodeNoComments(t *testing.T) {

	file, err := os.Open("images/im1.ppm")
	if err != nil {
		t.Errorf("Failed to open test image: %v\n", err.Error())
	}

	if _, err := Decode(file); err != nil {
		t.Errorf("Decode failed: %v\n", err.Error())
	}

}

// Test PPM image decoding with comments in the header.
func TestDecodeComments(t *testing.T) {

	file, err := os.Open("images/im2.ppm")
	if err != nil {
		t.Errorf("Failed to open test image: %v\n", err.Error())
	}

	if _, err := Decode(file); err != nil {
		t.Errorf("Decode failed: %v\n", err.Error())
	}

}
