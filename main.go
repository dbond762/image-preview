package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	images := make([][]byte, 0, len(urls))
	for _, url := range urls {
		image, err := loadImage(url)
		if err != nil {
			log.Printf("Error on get request to %s", url)
			continue
		}

		images = append(images, image)
	}

	log.Printf("%v", images)

	previews := make([][]byte, 0, len(images))
	for _, image := range images {
		preview, err := createPreview(image)
		if err != nil {
			log.Print(err)
			continue
		}

		previews = append(previews, preview)
	}

	previewUrls := make([]string, 0, len(previews))

	previewUrls = append(previewUrls, "https://golang.org/lib/godoc/images/go-logo-blue.svg")

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

func loadImage(url string) ([]byte, error) {
	client := &http.Client{}

	resp, err := client.Get("http://example.com")
	if err != nil {
		return nil, err
	}

	image, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return image, nil
}

func createPreview(image []byte) ([]byte, error) {
	return image, nil
}

func main() {
	imgdir := "imgs"
	if err := os.MkdirAll(imgdir, 666); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/preview", preview)

	port := 8000
	log.Printf("Server run at localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
