package main

import (
	"flag"

	"./internal/backupsession"
	"./internal/mylog"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

const defaultMaxVolumeSize = 1024 * 1024 //1M
const datadirectory = "../testdata"

func main() {
	mylog.Init()
	var maxVolSize = flag.Int("maxvolsize", defaultMaxVolumeSize, "Maximum volume size. If the size of a volume exceeds this number, a new one will be created")

	flag.Parse()

	session := backupsession.New(int64(*maxVolSize))

	log.Infof("Starting backup for directory %s", datadirectory)
	session.Process(datadirectory)

}
