package globalsettings

import (
	"log"
	"os"
	"os/user"

	"../constants"
)

//TempDir returns the directory for temporary files
func TempDir() string {
	path := HomeDir() + "/tmp"
	os.MkdirAll(path, constants.Filepermissions)
	return path
}

//HomeDir returns the home directory of this application (~/.<appname>)
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := usr.HomeDir + "/." + constants.Appname
	os.MkdirAll(path, constants.Filepermissions)
	return path
}
