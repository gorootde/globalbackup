package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"./internal/backupsession"
	"./internal/errorhandler"
	"./internal/mylog"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

const defaultMaxVolumeSize = 1024 * 1024 * 10 //10M

func main() {
	mylog.Init()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] SRC DST\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Environment variables:\n")
		fmt.Fprintf(os.Stderr, "  GOBACKUP_PASSWORD\n")
		fmt.Fprintf(os.Stderr, "        Variable containing the password for volume encryption\n")
		fmt.Fprintf(os.Stderr, "Parameters:\n")
		fmt.Fprintf(os.Stderr, "  SRC\n")
		fmt.Fprintf(os.Stderr, "        Directory to be backed up\n")
		fmt.Fprintf(os.Stderr, "  DST\n")
		fmt.Fprintf(os.Stderr, "        Destination of backup (URI)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	var maxVolSize = flag.Int("maxvolsize", defaultMaxVolumeSize, "Maximum volume size. If the size of a volume exceeds this number, a new one will be created")

	flag.Parse()

	password := os.Getenv("GOBACKUP_PASSWORD")
	params := flag.Args()
	if len(params) != 2 || len(password) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	src, err := filepath.Abs(params[0])
	errorhandler.ExitIfError(err, 1)
	dst, err := url.ParseRequestURI(params[1])
	errorhandler.ExitIfError(err, 1)

	log.Infof("Parameters:")
	log.Infof("  maxvolsize=%v", *maxVolSize)
	log.Infof("  src=%v", src)
	log.Infof("  dst=%v", dst)

	session := backupsession.New(int64(*maxVolSize), password)

	log.Infof("Starting backup for directory %s", src)
	session.Process(src)
	log.Infof("Finished")

}
