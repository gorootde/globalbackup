package globalsettings

import (
	"log"
	"os"
	"os/user"
)

//VolumeExtension defines the Filextension for backup folumes
const VolumeExtension = ".gbv"
const ManifestExtension = ".gmf"
const mkdirpermissions os.FileMode = 0700
const appname string = "gobackup"

//TempDir returns the directory for temporary files
func TempDir() string {
	path := HomeDir() + "/tmp"
	os.MkdirAll(path, mkdirpermissions)
	return path
}

//HomeDir returns the home directory of this application (~/.<appname>)
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := usr.HomeDir + "/." + appname
	os.MkdirAll(path, mkdirpermissions)
	return path
}
