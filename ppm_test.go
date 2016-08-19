package ppm

import (
	"bufio"
	"bytes"
	"testing"
)

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

	for i, tt := range magictest {

		bb := bytes.NewBuffer([]byte(tt.in))
		r := bufio.NewReader(bb)

		if err := decodeMagic(r); tt.valid && err != nil {
			t.Errorf("[%v] Failed to decode valid magic: %v", i, err)
		} else if !tt.valid && err == nil {
			t.Errorf("[%v] Decoded bad magic.", i)
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

	for i, tt := range commenttest {

		bb := bytes.NewBuffer([]byte(tt.in))
		r := bufio.NewReader(bb)

		if err := decodeComments(r); tt.valid && err != nil {
			t.Errorf("[%v] Failed to decode comments: %v", i, err)
		} else if !tt.valid && err == nil {
			t.Errorf("[%v] Decoded bad comments.", i)
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

	for i, tt := range widthheighttest {

		bb := bytes.NewBuffer([]byte(tt.in))
		r := bufio.NewReader(bb)

		if width, height, err := decodeWidthHeight(r); tt.valid && err != nil {
			t.Errorf("[%v] Failed to decode valid width/height: %v\n", i, err)
		} else if !tt.valid && err == nil {
			t.Errorf("[%v] Decoded bad width/height.", i)
		} else if tt.valid && (width != tt.width || height != tt.height) {
			t.Errorf("[%v] Wrong width/height.", i)
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
			t.Errorf("[%v] Failed to decode valid format: %v\n", i, err)
		} else if !tt.valid && err == nil {
			t.Errorf("[%v] Decoded bad format.", i)
		} else if tt.valid && format != tt.format {
			t.Errorf("[%v] Wrong format.", i)
		}

	}

}

// Test reading image data.
func TestDecodeImage(t *testing.T) {

	var datatest = []struct {
		width  int
		height int
		in     string
		valid  bool
	}{
		{1, 1, "000", true},
		{1, 1, "00", false},
		{0, 0, "", true},
		{-1, 0, "", false},
		{0, -1, "", false},
		{-1, -1, "", false},
		{3, 1, "000111222", true},
		{1, 3, "000111222", true},
		{2, 2, "000111222333", true},
	}

	for i, tt := range datatest {

		bb := bytes.NewBuffer([]byte(tt.in))
		r := bufio.NewReader(bb)

		if img, err := decodeImage(r, tt.width, tt.height); tt.valid && err != nil {
			t.Errorf("[%v] Failed to decode valid format: %v\n", i, err)
		} else if !tt.valid && err == nil {
			t.Errorf("[%v] Decoded bad image.\n", i)
		} else if tt.valid && err == nil &&

			(img.Bounds().Min.X != 0 ||
				img.Bounds().Min.Y != 0 ||
				img.Bounds().Max.X != tt.width ||
				img.Bounds().Max.Y != tt.height) {

			t.Errorf("[%v] Wrong image size.\n", i)

		}
	}

}

// Test the full line up.
func TestFull(t *testing.T) {

	var fulltest = []struct {
		in     string
		width  int
		height int
		valid  bool
	}{
		{"P6\n1 1\n255\n000", 1, 1, true},
		{"P6\n2 2\n255\n000111222333", 2, 2, true},
		{"P6\n-1 -1\n255\n000111222333", -1, -1, false},
		{"P6\n#a\n#b\n2 2\n255\n000111222333", 2, 2, true},
		{"P6\n%a\n#b\n2 2\n255\n000111222333", 2, 2, false},
		{"\n\n\n\n", 0, 0, false},
		{"2kj324lkj2lij2g32i", 0, 0, false},
	}

	for i, tt := range fulltest {

		bb := bytes.NewBuffer([]byte(tt.in))
		r := bufio.NewReader(bb)

		if img, err := Decode(r); tt.valid && err != nil {
			t.Errorf("[%v] Failed to decode valid image: %v", i, err)
		} else if !tt.valid && err == nil {
			t.Errorf("[%v] Decoded bad image.", i)
		} else if tt.valid && err == nil &&

			(img.Bounds().Min.X != 0 ||
				img.Bounds().Min.Y != 0 ||
				img.Bounds().Max.X != tt.width ||
				img.Bounds().Max.Y != tt.height) {

			t.Errorf("[%v] Wrong image size.\n", i)

		}

	}

}
