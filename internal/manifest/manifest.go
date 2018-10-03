package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"../backuprecord"
	"../constants"
	"../globalsettings"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("manifest")

//Manifest describing the content of a backup session
type Manifest struct {
	Source    string                      //Identifier of the backup source
	Timestamp int32                       //Timestamp of the backup operation
	Records   []backuprecord.Backuprecord //List of backup records
	path      string
}

func New(id string) Manifest {
	path := fmt.Sprintf("%s/gobackup-manifest-%s.json", globalsettings.TempDir(), id)

	log.Debugf("Manifest '%v' created", path)
	hostname, _ := os.Hostname()
	timestamp := int32(time.Now().Unix())
	return Manifest{path: path, Source: hostname, Timestamp: timestamp}
}

func (manifest *Manifest) Add(record backuprecord.Backuprecord) {
	manifest.Records = append(manifest.Records, record)

}

func (manifest Manifest) Persist() error {
	manifestfile, err := os.OpenFile(manifest.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, constants.Filepermissions)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	enc := json.NewEncoder(manifestfile)
	if err := enc.Encode(manifest); err != nil {
		return err
	}
	return nil
}

func (mainifest Manifest) Size() int {
	return len(mainifest.Records)
}
