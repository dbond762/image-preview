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

	log.Printf("%v", urls)

	previewUrls := make([]string, 0, len(urls))
	for _, url := range urls {
		img, err := loadImage(url)
		if err != nil {
			log.Printf("Error on load image from %s", url)
			continue
		}

		preview, err := createPreview(img)
		if err != nil {
			log.Print(err)
			continue
		}

		previewFile, err := saveImage(preview, []byte(url))
		if err != nil {
			log.Print(err)
			continue
		}

		previewUrls = append(previewUrls, staticPrefix+previewFile)
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
	if err := os.MkdirAll(staticDir, 666); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/preview", preview)

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle(staticPrefix, http.StripPrefix(staticPrefix, fs))

	port := 8000
	log.Printf("Server run at localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
