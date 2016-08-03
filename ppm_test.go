package ppm

import (
	"bufio"
	"bytes"
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

func TestDecodeMagic(t *testing.T) {

	var magictest = []struct {
		in    []byte
		valid bool
	}{
		{[]byte{'P', '6'}, true},
		{[]byte{' ', ' '}, false},
		{[]byte{'6', 'P'}, false},
		{[]byte{'p', '6'}, false},
		{[]byte{'#', 'L'}, false},
		{[]byte{}, false},
	}

	for _, tt := range magictest {
		bb := bytes.NewBuffer(tt.in)
		r := bufio.NewReader(bb)

		if err := decodeMagic(r); tt.valid && err != nil {
			t.Errorf("Failed to decode valid magic.")
		} else if !tt.valid && err == nil {
			t.Errorf("Decoded bad magic.")
		}

	}

}
