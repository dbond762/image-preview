package main

import (
	"fmt"
	"html/template"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
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

		previewFile, err := saveImage(preview)
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

func createPreview(img image.Image) (image.Image, error) {
	return img, nil
}

// todo gen names from urls
func saveImage(img image.Image) (string, error) {
	hash, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("%s/%s.png", staticDir, hash)
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
