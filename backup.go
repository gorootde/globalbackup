package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type backuprecord struct {
	Path         string
	Checksum     string
	LastModified time.Time
	Size         int64
}

func main() {

	//listrec("..")
	//fmt.Println(getfilechecksum("backup.go"))

	var entries []backuprecord
	err := filepath.Walk("..",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				abspath, _ := filepath.Abs(path)
				checksum := getFileChecksum(abspath)
				fmt.Println(path, checksum)
				record := backuprecord{Path: abspath, LastModified: info.ModTime(), Checksum: checksum, Size: info.Size()}
				entries = append(entries, record)

			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Hello World")
	fmt.Println(len(entries))
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(entries); err != nil {
		log.Println(err)
	}
}

func getFileChecksum(path string) string {
	hasher := sha512.New()
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
		return ""
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
