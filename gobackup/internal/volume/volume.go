package volume

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/tink-ab/tempfile"

	"../globalsettings"
)

const volsize = 1024 * 1024 * 10

//Volume is a set of bytes containing multiple files
type Volume struct {
	Path string
}

//New creates a new volume
func New() Volume {

	tmpfile, err := tempfile.TempFile(globalsettings.TempDir(), "volume", globalsettings.VolumeExtension)
	if err != nil {
		log.Fatal(err)
	}
	return Volume{tmpfile.Name()}
}

//FromFile creates volume from an existing volume file on disk
func FromFile(tarfile string) Volume {
	return Volume{tarfile}
}

//List returns a list of all files contained in this volume
func List(volume Volume) ([]string, error) {
	file, err := os.Open(volume.Path)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v", err)
	}
	tr := tar.NewReader(file)
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

//Add a file in local filesystem to a volume
func Add(volume Volume, path string) (int64, error) {
	abspath, _ := filepath.Abs(path)
	finfo, err := os.Lstat(abspath)
	if err != nil {
		return 0, fmt.Errorf("Error getting file info: %v", err)
	}

	datafile, err := os.Open(abspath)
	if err != nil {
		return 0, fmt.Errorf("Error opening datafile: %v", err)
	}

	tarfile, err := os.OpenFile(volume.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		return 0, fmt.Errorf("Error opening volume: %v", err)
	}
	defer tarfile.Close()

	tw := tar.NewWriter(tarfile)
	hdr, err := tar.FileInfoHeader(finfo, finfo.Name())
	if err != nil {
		return 0, fmt.Errorf("Error creating header: %v", err)
	}
	hdr.Name = abspath
	if err := tw.WriteHeader(hdr); err != nil {
		return 0, fmt.Errorf("Error writing header: %v", err)
	}

	bytes, err := io.Copy(tw, datafile)
	if err != nil {
		return 0, fmt.Errorf("Error writing data: %v", err)
	}

	if err := tw.Close(); err != nil {
		return 0, fmt.Errorf("Error closing volume: %v", err)
	}

	return bytes, nil
}
