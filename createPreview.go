package main

import (
	"github.com/dbond762/image-preview/onlineImage"
	"github.com/nfnt/resize"
)

func createPreview(images <-chan *onlineImage.OnlineImage, previews chan<- *onlineImage.OnlineImage) {
	for img := range images {
		previews <- &onlineImage.OnlineImage{
			Image: resize.Thumbnail(100, 100, img.Image, resize.Bicubic),
			URL:   img.URL,
		}
	}

	close(previews)
}
