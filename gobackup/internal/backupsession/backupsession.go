package backupsession

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"../backuprecord"
	"../constants"
	"../manifest"
	"../volume"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("backupsession")

//Backupsession controlls the process of backing up
type Backupsession struct {
	currentVolume      *volume.Volume
	manifest           manifest.Manifest
	maxVolSize         int64
	id                 string
	currentVolumeIndex int64
}

//New returns a backupsession instace. Parameter defines the max volume size.
func New(maxVolumeSize int64) Backupsession {
	id := time.Now().UTC().Format(constants.Dateformat)

	result := Backupsession{maxVolSize: maxVolumeSize, manifest: manifest.New(id), id: id}
	log.Infof("Backup Session created with id %v", id)
	return result
}

func (session *Backupsession) nextVolume() error {
	if session.currentVolume != nil {
		session.currentVolume.Close()
		session.currentVolume = nil
		session.currentVolumeIndex = session.currentVolumeIndex + 1
	}

	volumeid := volume.VolumeId(fmt.Sprintf("%s-%d", session.id, session.currentVolumeIndex))
	session.currentVolume = volume.New("verysecret", volumeid)
	return nil
}

//Process processes the given directory on harddisk recursively
func (session *Backupsession) Process(directory string) error {

	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			isSymlink := info.Mode()&os.ModeSymlink != 0
			if !info.IsDir() && !isSymlink {

				start := time.Now()
				session.nextVolume()
				processerr := session.processFile(path)
				session.manifest.Persist()
				elapsed := time.Since(start)
				if processerr != nil {
					log.Errorf("%s Error: %v\n", path, err)
				} else {
					log.Infof("%v [%v]\n", path, elapsed)
				}

			}

			return nil
		})
	if err != nil {
		log.Error(err)
	}
	session.currentVolume.Close()
	return nil
}

func (session *Backupsession) processFile(path string) error {
	record := backuprecord.New(path)
	_, err := session.currentVolume.Add(path)
	if err != nil {
		return fmt.Errorf("Unable to add file to volume: %v", err)
	}
	record.AddVolume(session.currentVolume.GetID())
	session.manifest.Add(record)
	return nil
}
