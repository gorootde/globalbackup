package volume

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"../encryptionwriter"
	"../globalsettings"
	"../splitfilewriter"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("volume")

type VolumeId string

//Volume is a set of bytes containing multiple files
type Volume struct {
	splitfilewriter *splitfilewriter.Splitfilewriter
	encwriter       *encryptionwriter.EncryptionWriter
	tarwriter       *tar.Writer
	id              VolumeId
}

//New creates a new volume. Secret is the passphrase for encrytion. volid is the id of the current backupsession
func New(secret string, volid VolumeId, volumechunksize int64) *Volume {

	archivepath := fmt.Sprintf("%s/gobackup-volume-%s.tar.gpg.part%%v", globalsettings.TempDir(), volid)

	splitwriter := splitfilewriter.New(archivepath, volumechunksize)
	encwriter := encryptionwriter.New(splitwriter, secret)
	tarwriter := tar.NewWriter(encwriter)

	return &Volume{splitwriter, encwriter, tarwriter, volid}
}

//GetID returns the id of the current volume
func (volume Volume) GetID() VolumeId {
	return volume.id
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
	if err := volume.tarwriter.WriteHeader(hdr); err != nil {
		return 0, fmt.Errorf("Error writing header: %v", err)
	}

	bytes, err := io.Copy(volume.tarwriter, datafile)
	if err != nil {
		return 0, fmt.Errorf("Error writing data: %v", err)
	}

	volume.tarwriter.Flush()
	volume.encwriter.Flush()
	return bytes, nil
}

func (volume Volume) Close() {
	if err := volume.tarwriter.Close(); err != nil {
		log.Fatalf("Error closing volume tw: %v", err)
	}

	if err := volume.encwriter.Close(); err != nil {
		log.Fatalf("Error closing volume gzw: %v", err)
	}
	defer volume.splitfilewriter.Close()

}
