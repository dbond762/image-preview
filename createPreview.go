package main

import (
	"github.com/nfnt/resize"
)

func createPreview(images <-chan *OnlineImage, previews chan<- *OnlineImage) {
	for img := range images {
		previews <- &OnlineImage{
			Image: resize.Thumbnail(100, 100, img.Image, resize.Bicubic),
			URL:   img.URL,
		}
	}

	close(previews)
}
