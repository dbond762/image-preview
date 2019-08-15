package main

import (
	"image"
	"net/http"
)

func loadImage(urls []string, images chan<- *OnlineImage, errors chan<- error) {
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

		images <- &OnlineImage{
			Image: img,
			URL:   url,
		}
	}

	close(images)
	close(errors)
}
