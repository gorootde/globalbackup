package constants

import "os"

//Filepermissions defines the default chmod permissions for new created files / directories
const Filepermissions os.FileMode = 0700

//Appname is the name of this application used e.g. in filenames
const Appname string = "gobackup"

//Dateformat for filenames
const Dateformat string = "20060102T150405Z"
