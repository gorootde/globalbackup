package mylog

import (
	"os"

	"github.com/op/go-logging"
)

var format1 = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{color:reset} %{message}`,
)
var format2 = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func Init() {
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	//backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format1)
	//backend2Formatter := logging.NewBackendFormatter(backend2, format2)

	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(logging.ERROR, "volume")
	logging.SetBackend(backend1Formatter)

}
