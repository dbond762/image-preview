package main

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func getHash(url string) string {
	md5Hash := md5.Sum([]byte(url))
	return hex.EncodeToString(md5Hash[:])
}
