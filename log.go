package gojob

import (
  "log"
	"io"
	"io/ioutil"
	"os"
)

var (
	Trace *log.Logger
  Info *log.Logger
	Error *log.Logger
)

func InitLog(traceHandle io.Writer,
						 infoHandle io.Writer,
						 errorHandle io.Writer) {

  Trace = log.New(traceHandle,
    "INFO: ",
    log.Ldate|log.Ltime)
  Info = log.New(infoHandle,
    "INFO: ",
    log.Ldate|log.Ltime)
	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func InitLogInfo() {
  InitLog(ioutil.Discard,
					os.Stdout,
					os.Stderr)
}

func InitLogTrace() {
  InitLog(os.Stdout,
					os.Stdout,
					os.Stderr)
}
