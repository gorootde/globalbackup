package backuprecord

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/silvasur/golibrsync/librsync"
)

//Backuprecord is a record describing a single file in a backup set
type Backuprecord struct {
	Path         string
	Checksum     string
	LastModified time.Time
	Size         int64
}

//New creates a new backup record based on a given file on harddisk
func New(filepath string, info os.FileInfo, chunksize int64) Backuprecord {
	fchecksum := getFileChecksum(filepath)
	chunksums := getFileSliceChecksums(filepath, chunksize)
	getRsyncSignature(filepath)
	result := Backuprecord{Path: filepath, Checksum: fchecksum, ChunkSums: chunksums, LastModified: info.ModTime(), Size: info.Size()}
	return result
}

func getChecksum(path string) (string, error) {
	abspath, _ := filepath.Abs(path)
	file, err := os.Open(abspath)
	if err != nil {
		return "", fmt.Errorf("Error opening file: %v", err)
	}
	signature := new(bytes.Buffer)
	err = librsync.CreateSignature(file, bufio.NewWriter(signature))
	if err != nil {
		return "", fmt.Errorf("Error creading signature: %v", err)
	}

	ssig := hex.EncodeToString(signature.Bytes())
	return ssig, nil
}

func getFileSliceChecksums(path string, chunksize int64) []string {

	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer f.Close()

	var result []string
	for {

		chunkReader := io.LimitReader(bufio.NewReader(f), chunksize)

		hasher := sha512.New()
		numbytes, err := io.Copy(hasher, chunkReader)

		if err != nil {
			log.Fatal(err)
			break
		}
		if numbytes == 0 {
			break
		}
		hash := hex.EncodeToString(hasher.Sum(nil))
		result = append(result, hash)
	}
	return result
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
