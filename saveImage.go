package main

import (
	"github.com/dbond762/image-preview/fileManager"
	"github.com/dbond762/image-preview/onlineImage"
)

func saveImage(manager *fileManager.FileManager, images <-chan *onlineImage.OnlineImage, names chan<- string, errors chan<- error) {
	for img := range images {
		fileName, err := manager.Add(img)
		if err != nil {
			errors <- err
		}

		names <- fileName[len(staticDir):]
	}

	close(names)
	close(errors)
}
