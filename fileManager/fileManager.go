package fileManager

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/dbond762/image-preview/onlineImage"
	"github.com/dbond762/image-preview/timer"
)

var (
	ErrFileNotFound = errors.New("File not found")
)

type FileManager struct {
	Files      []File
	Timeout    time.Duration
	NumOfFiles int
	StaticDir  string
}

type File struct {
	onlineImage.OnlineImage
	Name         string
	Cancel       context.CancelFunc
	Done         chan struct{}
	ResetTimeout chan struct{}
}

func New(staticDir string, timeout time.Duration, numOfFiles int) *FileManager {
	fm := &FileManager{
		Files:      []File{},
		Timeout:    timeout,
		NumOfFiles: numOfFiles,
		StaticDir:  staticDir,
	}

	go fm.manage()

	return fm
}

func (fm *FileManager) manage() {
	for {
		for i := 0; i < len(fm.Files); i++ {
			select {
			case <-fm.Files[i].Done:
				errors := make(chan error)

				go fm.deleteFile(fm.Files[i].Name, errors)

				go func(errors chan error) {
					for err := range errors {
						log.Print(err)
					}
				}(errors)

				fm.Files = append(fm.Files[:i], fm.Files[i+1:]...)
				i = i - 1
			default:
				// todo maybe add throtle
				continue
			}
		}
	}
}

func (fm *FileManager) Add(img *onlineImage.OnlineImage) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())

	var (
		done         = make(chan struct{})
		resetTimeout = make(chan struct{})
	)

	md5Hash := md5.Sum([]byte(img.URL))
	hash := hex.EncodeToString(md5Hash[:])

	name := fmt.Sprintf("%s%s.png", fm.StaticDir, hash)

	file := File{
		OnlineImage:  *img,
		Name:         name,
		Cancel:       cancel,
		Done:         done,
		ResetTimeout: resetTimeout,
	}

	// todo control files count
	fm.Files = append(fm.Files, file)

	go func(ctx context.Context, resetTimeout <-chan struct{}, done chan<- struct{}) {
		timeIsOver := make(chan struct{})
		log.Print(name, "start")

		go timer.Start(fm.Timeout, resetTimeout, timeIsOver)

		select {
		case <-ctx.Done():
		case <-timeIsOver:
		}
		log.Print(name, "done")

		done <- struct{}{}
	}(ctx, resetTimeout, done)

	if err := fm.saveFile(&file); err != nil {
		return "", err
	}

	return file.Name, nil
}

func (fm *FileManager) FindFileAndResetTimeout(url string) (*File, error) {
	for _, file := range fm.Files {
		if file.URL == url {
			file.ResetTimeout <- struct{}{}
			return &file, nil
		}
	}

	return nil, ErrFileNotFound
}

func (fm *FileManager) saveFile(file *File) error {
	f, err := os.Create(file.Name)
	if err != nil {
		return err
	}

	if err := png.Encode(f, file.Image); err != nil {
		if err := f.Close(); err != nil {
			return err
		}
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func (fm *FileManager) deleteFile(fileName string, errors chan<- error) {
	if err := os.Remove(fileName); err != nil {
		errors <- err
	}

	close(errors)
}
