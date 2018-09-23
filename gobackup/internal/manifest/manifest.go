package manifest

import (
	"encoding/json"
	"fmt"
	"os"

	"../backuprecord"
	"../constants"
	"../globalsettings"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("manifest")

type Manifest struct {
	records []backuprecord.Backuprecord
	path    string
}

func New(id string) Manifest {
	path := fmt.Sprintf("%s/gobackup-manifest-%s.json", globalsettings.TempDir(), id)

	log.Debugf("Manifest '%v' created", path)
	return Manifest{path: path}
}

func (manifest *Manifest) Add(record backuprecord.Backuprecord) {
	manifest.records = append(manifest.records, record)

}

func (manifest Manifest) Persist() error {
	manifestfile, err := os.OpenFile(manifest.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, constants.Filepermissions)
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
