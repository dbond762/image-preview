package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/dbond762/image-preview/fileManager"
	"github.com/dbond762/image-preview/onlineImage"
)

type PreviewRequest struct {
	URLs []string `json:"urls"`
}

type TemplateData struct {
	Previews []string
}

type PreviewHandler struct {
	manager    *fileManager.FileManager
	numOfFiles int
}

func (ph *PreviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var previewRequest PreviewRequest
	if err := decoder.Decode(&previewRequest); err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	urls := previewRequest.URLs

	if len(urls) > ph.numOfFiles {
		http.Error(w, fmt.Sprintf("To many images. It shoud be %d", ph.numOfFiles), http.StatusBadRequest)
		return
	}

	log.Print(urls)

	var (
		images   = make(chan *onlineImage.OnlineImage)
		previews = make(chan *onlineImage.OnlineImage)
		names    = make(chan string)

		loadErrors = make(chan error)
		saveErrors = make(chan error)
	)

	for i := 0; i < len(urls); i++ {
		file, err := ph.manager.FindFileAndResetTimeout(urls[i])
		if err != nil || err == fileManager.ErrFileNotFound {
			log.Print(err)
			continue
		}

		log.Printf("%v", file)

		urls = append(urls[:i], urls[i+1:]...)
		i = i - 1

		names <- file.Name[len(staticDir):]
	}

	go loadImage(urls, images, loadErrors)
	go createPreview(images, previews)
	go saveImage(ph.manager, previews, names, saveErrors)

	errors := []chan error{
		loadErrors,
		saveErrors,
	}
	for _, errorChan := range errors {
		go func(errorChan <-chan error) {
			for err := range errorChan {
				log.Print(err)
			}
		}(errorChan)
	}

	previewUrls := make([]string, 0, len(urls))
	for name := range names {
		previewUrls = append(previewUrls, staticPrefix+name)
	}

	log.Print(previewUrls)

	data := TemplateData{
		Previews: previewUrls,
	}

	t, err := template.ParseFiles("tmpl/previews.tmpl")
	if err != nil {
		log.Print(err)
	}
	t.Execute(w, data)
}
