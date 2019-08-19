package main

import (
	"image"
)

type OnlineImage struct {
	image.Image
	URL string
}
