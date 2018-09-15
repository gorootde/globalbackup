package backupsession

import (
	"fmt"

	"../backuprecord"
	"../manifest"
	"../mylog"
	"../volume"
)

//Backupsession controlls the process of backing up
type Backupsession struct {
	currentVolume   volume.Volume
	currentManifest manifest.Manifest
	maxVolSize      int64
}

//New returns a backupsession instace. Parameter defines the max volume size.
func New(maxVolumeSize int64) Backupsession {
	result := Backupsession{maxVolSize: maxVolumeSize, currentVolume: volume.New(), currentManifest: manifest.New()}
	return result
}

//Process processes the given file on harddisk (adds it to the backup)
func (session *Backupsession) Process(path string) error {

	if session.currentVolume.Size() > session.maxVolSize {
		mylog.Log.Debugf("Current volume size is %v. Max volsize is %v. Creating new volume.", session.currentVolume.Size(), session.maxVolSize)

		session.currentManifest.Persist()
		session.currentVolume.Close()
		session.currentVolume = volume.New()
		session.currentManifest = manifest.New()
	}
	record := backuprecord.New(path)
	_, err := session.currentVolume.Add(path)
	if err != nil {
		return fmt.Errorf("Unable to add file to volume: %v", err)
	}
	session.currentManifest.Add(record)
	session.currentManifest.Persist()
	return nil
}
