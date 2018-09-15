package volume

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"../globalsettings"
	"../mylog"
	"github.com/tink-ab/tempfile"
)

//Volume is a set of bytes containing multiple files
type Volume struct {
	file   *os.File
	writer *tar.Writer
}

//New creates a new volume
func New() Volume {

	tmpfile, err := tempfile.TempFile(globalsettings.TempDir(), "volume", globalsettings.VolumeExtension)
	if err != nil {
		log.Fatal(err)
	}
	archivepath := tmpfile.Name()

	volfile, err := os.OpenFile(archivepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		log.Fatalf("Error opening volume file: %v", err)
	}

	writer := tar.NewWriter(volfile)
	mylog.Log.Debugf("Volume '%v' created", archivepath)
	return Volume{volfile, writer}
}

//List returns a list of all files contained in this volume
func (volume Volume) List() ([]string, error) {
	tr := tar.NewReader(volume.file)
	var result []string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return nil, fmt.Errorf("Error listing files: %v", err)
		}
		result = append(result, hdr.Name)
	}
	return result, nil
}

//Add a file in local filesystem to a volume. Returns bytes written.
func (volume Volume) Add(path string) (int64, error) {
	abspath, _ := filepath.Abs(path)
	finfo, err := os.Lstat(abspath)
	if err != nil {
		return 0, fmt.Errorf("Error getting file info: %v", err)
	}

	datafile, err := os.Open(abspath)
	if err != nil {
		return 0, fmt.Errorf("Error opening datafile: %v", err)
	}

	if err != nil {
		return 0, fmt.Errorf("Error opening volume: %v", err)
	}

	hdr, err := tar.FileInfoHeader(finfo, finfo.Name())
	if err != nil {
		return 0, fmt.Errorf("Error creating header: %v", err)
	}
	hdr.Name = abspath
	if err := volume.writer.WriteHeader(hdr); err != nil {
		return 0, fmt.Errorf("Error writing header: %v", err)
	}

	bytes, err := io.Copy(volume.writer, datafile)
	if err != nil {
		return 0, fmt.Errorf("Error writing data: %v", err)
	}

	volume.writer.Flush()
	return bytes, nil
}

func (volume Volume) Close() {
	if err := volume.writer.Close(); err != nil {
		log.Fatalf("Error closing volume: %v", err)
	}
	defer volume.file.Close()
	mylog.Log.Debugf("Volume '%v' closed", volume.file.Name())
}

//Size returns the size of this volume
func (volume Volume) Size() int64 {
	finfo, err := os.Lstat(volume.file.Name())
	if err != nil {
		log.Fatal(err)
	}
	return finfo.Size()
}
