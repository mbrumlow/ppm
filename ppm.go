package ppm

import (
	"bufio"
	"fmt"
	"image"
	"io"
)

func Decode(r io.Reader) (image.Image, error) {
	b := bufio.NewReader(r)
	return decode(b)
}

func decode(r *bufio.Reader) (image.Image, error) {

	// Check format header
	for _, t := range []byte{'P', '6'} {
		if c, err := r.ReadByte(); err != nil {
			return nil, fmt.Errorf("Failed to read header: %v", err.Error())
		} else if c != t {
			return nil, fmt.Errorf("Invalid image format.")
		}
	}

	// Read comments and white space.
	for {
		if c, err := r.ReadByte(); err != nil {
			return nil, fmt.Errorf("Failed to read comments: %v", err.Error())
		} else if c == '#' {
			if _, err := r.ReadBytes('\n'); err != nil {
				return nil, fmt.Errorf("Failed to read white space: %v", err.Error())
			}
		} else if c != '\n' {
			break
		}
	}

	// Get the image size.
	var width, height int
	if n, err := fmt.Fscanf(r, "%d %d\n", &width, &height); err != nil {
		return nil, fmt.Errorf("Failed to read size and format: %v", err.Error())
	} else if n != 2 {
		return nil, fmt.Errorf("Failed to read image size.")
	}

	// Get image color format.
	var rgb_color int
	if n, err := fmt.Fscanf(r, "%d\n", &rgb_color); err != nil {
		return nil, fmt.Errorf("Failed to read color format: %v", err.Error())
	} else if n != 1 {
		return nil, fmt.Errorf("Failed to read color format.")
	}

	// Read the image data.
	b := make([]byte, 3)
	i := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width*height; x++ {

		if _, err := io.ReadFull(r, b); err != nil {
			return nil, fmt.Errorf("Failed to read image data: %v", err.Error())
		}

		i.Pix[x+0] = b[0]
		i.Pix[x+1] = b[1]
		i.Pix[x+2] = b[2]
		i.Pix[x+0] = 0x00

	}

	return i, nil
}
