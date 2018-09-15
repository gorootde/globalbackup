package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"./internal/backupsession"
	"./internal/mylog"
)

const defaultMaxVolumeSize = 1024 * 1024 //1M
const datadirectory = "../testdata"

func main() {
	mylog.Init()
	var maxVolSize = flag.Int("maxvolsize", defaultMaxVolumeSize, "Maximum volume size. If the size of a volume exceeds this number, a new one will be created")

	flag.Parse()

	session := backupsession.New(int64(*maxVolSize))

	mylog.Log.Infof("Starting backup for directory %s", datadirectory)
	err := filepath.Walk(datadirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			isSymlink := info.Mode()&os.ModeSymlink != 0
			if !info.IsDir() && !isSymlink {
				start := time.Now()
				processerr := session.Process(path)
				elapsed := time.Since(start)

				if processerr != nil {
					mylog.Log.Errorf("%s Error: %v\n", path, err)
				} else {
					mylog.Log.Infof("%v OK\n", path)
				}
				mylog.Log.Debugf("Processing this file took %v", elapsed)
			}

			return nil
		})
	if err != nil {
		mylog.Log.Error(err)
	}

}
