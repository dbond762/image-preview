package main

func saveImage(fm *FileManager, images <-chan *OnlineImage, names chan<- string, errors chan<- error) {
	for img := range images {
		fileName, err := fm.Add(img)
		if err != nil {
			errors <- err
		}

		names <- fileName[len(staticDir):]
	}

	close(names)
	close(errors)
}
