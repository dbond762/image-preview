package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	staticDir    = "imgs"
	staticPrefix = "/imgs/"

	port = 8000
)

type Data struct {
	Previews []string
}

func preview(w http.ResponseWriter, r *http.Request) {
	urls, ok := r.URL.Query()["url"]
	if !ok {
		http.Error(w, "Haven`t url parameter in request", http.StatusBadRequest)
		return
	}

	log.Print(urls)

	var (
		images   = make(chan *OnlineImage)
		previews = make(chan *OnlineImage)
		names    = make(chan string)

		loadErrors = make(chan error)
		saveErrors = make(chan error)
	)

	go loadImage(urls, images, loadErrors)
	go createPreview(images, previews)
	go saveImage(previews, names, saveErrors)

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

	log.Printf("%v", previewUrls)

	data := Data{
		Previews: previewUrls,
	}

	t, err := template.ParseFiles("tmpl/previews.tmpl")
	if err != nil {
		log.Print(err)
	}
	t.Execute(w, data)
}

func main() {
	if err := os.MkdirAll(staticDir, 644); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/preview", preview)

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle(staticPrefix, http.StripPrefix(staticPrefix, fs))

	log.Printf("Server run at localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
