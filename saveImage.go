package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func saveImage(img image.Image, url []byte) (string, error) {
	md5Hash := md5.Sum(url)
	hash := hex.EncodeToString(md5Hash[:])

	name := fmt.Sprintf("%s/%s.png", staticDir, hash)

	log.Print(name)

	f, err := os.Create(name)
	if err != nil {
		return "", nil
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		return "", err
	}
	// todo handle err
	f.Close()

	return fmt.Sprintf("%s.png", hash), nil
}
