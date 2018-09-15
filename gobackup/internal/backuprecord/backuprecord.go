package backuprecord

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
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
func New(path string) Backuprecord {
	abspath, _ := filepath.Abs(path)
	finfo, err := os.Lstat(abspath)
	if err != nil {
		log.Fatal("Error getting file info", err)
	}
	checksum, err := getChecksum(abspath)
	if err != nil {
		log.Fatal("Error calculating checksum", err)
	}
	result := Backuprecord{Path: abspath, Checksum: checksum, LastModified: finfo.ModTime(), Size: finfo.Size()}
	return result
}

func getChecksum(path string) (string, error) {
	file, err := os.Open(path)
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
