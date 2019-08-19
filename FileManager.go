package main

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
)

var (
	errFileNotFound = errors.New("File not found")
)

type FileManager struct {
	files      []File
	timeout    time.Duration
	numOfFiles int
}

type File struct {
	OnlineImage
	name         string
	cancel       context.CancelFunc
	done         chan struct{}
	resetTimeout chan struct{}
}

func NewFileManager(timeout time.Duration, numOfFiles int) *FileManager {
	fm := &FileManager{
		files:      []File{},
		timeout:    timeout,
		numOfFiles: numOfFiles,
	}

	go func(fm *FileManager) {
		for {
			for i := 0; i < len(fm.files); i++ {
				select {
				case <-fm.files[i].done:
					errors := make(chan error)

					go fm.DeleteFile(fm.files[i].name, errors)

					go func(errors chan error) {
						for err := range errors {
							log.Print(err)
						}
					}(errors)

					fm.files = append(fm.files[:i], fm.files[i+1:]...)
					i = i - 1
				default:
					continue
				}
			}
		}
	}(fm)

	return fm
}

func (fm *FileManager) Add(img *OnlineImage) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())

	var (
		done         = make(chan struct{})
		resetTimeout = make(chan struct{})
	)

	md5Hash := md5.Sum([]byte(img.URL))
	hash := hex.EncodeToString(md5Hash[:])

	name := fmt.Sprintf("%s%s.png", staticDir, hash)

	file := File{
		OnlineImage:  *img,
		name:         name,
		cancel:       cancel,
		done:         done,
		resetTimeout: resetTimeout,
	}

	// todo control files count
	fm.files = append(fm.files, file)

	go func(ctx context.Context, resetTimeout <-chan struct{}, done chan<- struct{}) {
		timeIsOver := make(chan struct{})

		go func(resetTimeout <-chan struct{}, timeISOver chan<- struct{}) {
			// todo reset timeout
			select {
			case <-resetTimeout:
			default:
			}
			time.Sleep(fm.timeout)
			timeISOver <- struct{}{}
		}(resetTimeout, timeIsOver)

		select {
		case <-ctx.Done():
		case <-timeIsOver:
		}

		done <- struct{}{}
	}(ctx, resetTimeout, done)

	if err := fm.SaveFile(&file); err != nil {
		return "", err
	}

	return file.name, nil
}

func (fm *FileManager) FindFileAndResetTimeout(url string) (*File, error) {
	for _, file := range fm.files {
		if file.URL == url {
			file.resetTimeout <- struct{}{}
			return &file, nil
		}
	}

	return nil, errFileNotFound
}

func (fm *FileManager) SaveFile(file *File) error {
	f, err := os.Create(file.name)
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

func (fm *FileManager) DeleteFile(fileName string, errors chan<- error) {
	if err := os.Remove(fileName); err != nil {
		errors <- err
	}

	close(errors)
}
