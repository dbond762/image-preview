package main

import (
	"image"

	"github.com/nfnt/resize"
)

func createPreview(img image.Image) (image.Image, error) {
	preview := resize.Thumbnail(100, 100, img, resize.Bicubic)
	return preview, nil
}
