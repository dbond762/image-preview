package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/dbond762/image-preview/onlineImage"
	_ "golang.org/x/image/webp"
)

func loadImage(urls []string, images chan<- *onlineImage.OnlineImage, errors chan<- error) {
	client := &http.Client{}

	for _, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			resp.Body.Close()
			errors <- err
			continue
		}

		img, _, err := image.Decode(resp.Body)
		resp.Body.Close()
		if err != nil {
			errors <- err
			continue
		}

		images <- &onlineImage.OnlineImage{
			Image: img,
			URL:   url,
		}
	}

	close(images)
	close(errors)
}
