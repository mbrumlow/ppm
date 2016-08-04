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

	// Check magic.
	if err := decodeMagic(r); err != nil {
		return nil, fmt.Errorf("Failed to decode magic: %v", err)
	}

	// Read comments and white space.
	if err := decodeComments(r); err != nil {
		return nil, fmt.Errorf("Failed to decode comments: %v", err)
	}

	// Get the image size.
	width, height, err := decodeWidthHeight(r)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode width/height: %v", err)
	}

	// Get image color format.
	if _, err := decodeFormat(r); err != nil {
		return nil, fmt.Errorf("Failed to decode format; %v", err)
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

func decodeMagic(r *bufio.Reader) error {

	s, err := r.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Failed to read magic: %v", err.Error())
	}

	if s != "P6\n" {
		return fmt.Errorf("Invalid magic.")
	}

	return nil

}

func decodeComments(r *bufio.Reader) error {

	for {
		if b, err := r.Peek(1); err != nil {
			return err
		} else if b[0] == '#' {
			if _, err := r.ReadBytes('\n'); err != nil {
				return fmt.Errorf("Failed to read white space: %v", err.Error())
			}
		} else {
			break
		}
	}

	return nil
}

func decodeWidthHeight(r *bufio.Reader) (int, int, error) {

	var width, height int
	if n, err := fmt.Fscanf(r, "%d %d\n", &width, &height); err != nil {
		return width, height, fmt.Errorf("Failed to read width/height: %v", err.Error())
	} else if n != 2 {
		return width, height, fmt.Errorf("Failed to read image width/height.")
	}

	if width < 0 || height < 0 {
		return 0, 0, fmt.Errorf("out of range.")
	}

	return width, height, nil

}

func decodeFormat(r *bufio.Reader) (int, error) {

	var rgb_color int
	if n, err := fmt.Fscanf(r, "%d\n", &rgb_color); err != nil {
		return rgb_color, fmt.Errorf("Failed to read format: %v", err.Error())
	} else if n != 1 {
		return rgb_color, fmt.Errorf("Failed to read format.")
	}

	if rgb_color != 255 {
		return rgb_color, fmt.Errorf("Invalid or unsupported color format.")
	}

	return rgb_color, nil
}
