package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image/png"
	"os"
)

func saveImage(images <-chan *OnlineImage, names chan<- string, errors chan<- error) {
	for img := range images {
		md5Hash := md5.Sum([]byte(img.URL))
		hash := hex.EncodeToString(md5Hash[:])

		name := fmt.Sprintf("%s/%s.png", staticDir, hash)

		f, err := os.Create(name)
		if err != nil {
			errors <- err
			continue
		}

		if err := png.Encode(f, img); err != nil {
			if err := f.Close(); err != nil {
				errors <- err
			}
			errors <- err
			continue
		}

		if err := f.Close(); err != nil {
			errors <- err
			continue
		}

		names <- fmt.Sprintf("%s.png", hash)
	}

	close(names)
	close(errors)
}
