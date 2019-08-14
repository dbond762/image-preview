package main

import (
	"image"
	"net/http"
)

func loadImage(url string) (image.Image, error) {
	// todo make one client
	client := &http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return img, nil
}
