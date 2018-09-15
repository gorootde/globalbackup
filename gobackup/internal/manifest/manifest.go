package manifest

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"../backuprecord"
	"../globalsettings"
	"../mylog"
	"github.com/tink-ab/tempfile"
)

type Manifest struct {
	records []backuprecord.Backuprecord
	path    string
}

func New() Manifest {
	tmpfile, err := tempfile.TempFile(globalsettings.TempDir(), "manifest", globalsettings.ManifestExtension)
	if err != nil {
		log.Fatal(err)
	}
	mylog.Log.Debugf("Manifest '%v' created", tmpfile.Name())
	return Manifest{path: tmpfile.Name()}
}

func (manifest *Manifest) Add(record backuprecord.Backuprecord) {
	manifest.records = append(manifest.records, record)

}

func (manifest Manifest) Persist() error {
	manifestfile, err := os.OpenFile(manifest.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	enc := json.NewEncoder(manifestfile)
	if err := enc.Encode(manifest.records); err != nil {
		return err
	}
	return nil
}

func (mainifest Manifest) Size() int {
	return len(mainifest.records)
}
