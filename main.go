package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	staticDir    = "imgs/"
	staticPrefix = "/imgs/"

	// todo move port to flag
	port       = 8000
	timeout    = 20 * time.Second
	numOfFiles = 4
)

func main() {
	if err := os.RemoveAll(staticDir); err != nil {
		log.Print(err)
	}

	if err := os.MkdirAll(staticDir, 644); err != nil {
		log.Fatal(err)
	}

	preview := &PreviewHandler{
		fileManager: NewFileManager(timeout, numOfFiles),
		numOfFiles:  numOfFiles,
	}

	http.Handle("/preview", preview)

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle(staticPrefix, http.StripPrefix(staticPrefix, fs))

	log.Printf("Server run at localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
