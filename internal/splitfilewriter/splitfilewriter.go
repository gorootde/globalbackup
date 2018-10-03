package splitfilewriter

import (
	"fmt"
	"log"
	"os"

	"../constants"
)

type Splitfilewriter struct {
	bytesWritten     int64
	limit            int64
	currentFileIndex int
	filepattern      string
	currentFile      *os.File
}

func New(filenamepattern string, limit int64) *Splitfilewriter {
	result := &Splitfilewriter{bytesWritten: 0, limit: limit, currentFileIndex: 0, filepattern: filenamepattern, currentFile: nil}
	result.openFile()
	return result
}

func (spfw *Splitfilewriter) nextFile() {
	if spfw.currentFile != nil {
		spfw.Close()
		spfw.currentFileIndex++
		spfw.bytesWritten = 0
	}
	spfw.openFile()

}

func (spfw *Splitfilewriter) openFile() {
	archivepath := fmt.Sprintf(spfw.filepattern, spfw.currentFileIndex)
	file, err := os.OpenFile(archivepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, constants.Filepermissions)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	spfw.currentFile = file
}

func (spfw *Splitfilewriter) Close() {
	if spfw.currentFile != nil {
		spfw.currentFile.Close()
		spfw.currentFile = nil
	}
}

func (spfw *Splitfilewriter) Write(p []byte) (n int, err error) {
	bytesWritten, err := spfw.currentFile.Write(p)
	spfw.bytesWritten += int64(bytesWritten)
	if spfw.bytesWritten >= spfw.limit {
		spfw.nextFile()
	}
	return bytesWritten, err
}
