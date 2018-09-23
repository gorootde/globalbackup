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

	"../volume"

	"github.com/silvasur/golibrsync/librsync"
)

//Backuprecord is a record describing a single file in a backup set
type Backuprecord struct {
	Path         string
	Signature    string
	LastModified time.Time
	Size         int64
	Volumes      []volume.VolumeId
}

//New creates a new backup record based on a given file on harddisk
func New(path string) Backuprecord {
	abspath, _ := filepath.Abs(path)
	finfo, err := os.Lstat(abspath)
	if err != nil {
		log.Fatal("Error getting file info", err)
	}
	signature, err := getSignature(abspath)
	if err != nil {
		log.Fatal("Error calculating checksum", err)
	}
	result := Backuprecord{Path: abspath, Signature: signature, LastModified: finfo.ModTime(), Size: finfo.Size()}
	return result
}

func (record *Backuprecord) AddVolume(volid volume.VolumeId) {
	record.Volumes = append(record.Volumes, volid)
}

//getSignature calculates the librsync signature
func getSignature(path string) (string, error) {
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
