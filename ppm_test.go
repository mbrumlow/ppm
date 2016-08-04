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
func TestDecodeWithComments(t *testing.T) {

	file, err := os.Open("images/im2.ppm")
	if err != nil {
		t.Errorf("Failed to open test image: %v\n", err.Error())
	}

	if _, err := Decode(file); err != nil {
		t.Errorf("Decode failed: %v\n", err.Error())
	}

}

// Test reading the magic.
func TestDecodeMagic(t *testing.T) {

	var magictest = []struct {
		in    string
		valid bool
	}{
		{"P6\n", true},
		{"P6", false},
		{"P6 ", false},
		{"  ", false},
		{"6P", false},
		{"p6", false},
		{"#L", false},
		{"", false},
	}

	for _, tt := range magictest {

		bb := bytes.NewBuffer([]byte(tt.in))
		r := bufio.NewReader(bb)

		if err := decodeMagic(r); tt.valid && err != nil {
			t.Errorf("Failed to decode valid magic: %v", err)
		} else if !tt.valid && err == nil {
			t.Errorf("Decoded bad magic.")
		}

	}

}

// Tesst reading comments.
func TestDecodeComments(t *testing.T) {

	var commenttest = []struct {
		in    string
		valid bool
	}{
		{"", false},
		{"#\n500 500", true},
		{"#\n# comment\n500 500", true},
		{"# comment\n500 500", true},
		{"# comment 1\n#comment 2\n500 500", true},
	}

	for _, tt := range commenttest {

		bb := bytes.NewBuffer([]byte(tt.in))
		r := bufio.NewReader(bb)

		if err := decodeComments(r); tt.valid && err != nil {
			t.Errorf("Failed to decode comments: %v", err)
		} else if !tt.valid && err == nil {
			t.Errorf("Decoded bad comments.")
		}
	}

}

// Test reading width and height.
func TestDecodeWidthHeight(t *testing.T) {

	var widthheighttest = []struct {
		in     string
		width  int
		height int
		valid  bool
	}{
		{"", 0, 0, false},
		{"\n", 0, 0, false},
		{"800 600\n", 800, 600, true},
		{"1000 1000\n", 1000, 1000, true},
		{"-1 1000\n", -1, 1000, false},
		{"-1 -1\n", -1, -1, false},
	}

	for _, tt := range widthheighttest {

		bb := bytes.NewBuffer([]byte(tt.in))
		r := bufio.NewReader(bb)

		if width, height, err := decodeWidthHeight(r); tt.valid && err != nil {
			t.Errorf("Failed to decode valid width/height: %v\n", err)
		} else if !tt.valid && err == nil {
			t.Errorf("Decoded bad width/height.")
		} else if tt.valid && (width != tt.width || height != tt.height) {
			t.Errorf("Wrong width/height.")
		}
	}

}

// Test reading format.
func TestDecodeFormat(t *testing.T) {

	var formattest = []struct {
		in     string
		format int
		valid  bool
	}{
		{"255\n", 255, true},
		{"100\n", -1, false},
		{"-1\n", -1, false},
		{"", 0, false},
		{"\n", 0, false},
	}

	for i, tt := range formattest {

		bb := bytes.NewBuffer([]byte(tt.in))
		r := bufio.NewReader(bb)

		if format, err := decodeFormat(r); tt.valid && err != nil {
			t.Errorf("Failed to decode valid format: %v\n", err)
		} else if !tt.valid && err == nil {
			t.Errorf("Decoded bad format: %d.", i)
		} else if tt.valid && format != tt.format {
			t.Errorf("Wrong format.")
		}

	}

}
