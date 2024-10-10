package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	InfoLog  *log.Logger
	WarnLog  *log.Logger
	ErrorLog *log.Logger
)

func init() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		InfoLog.Output(2, fmt.Sprintf("Failed to open log file: %v", err))
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	InfoLog = log.New(multiWriter, "INFO:\t", log.Ldate|log.Ltime)
	WarnLog = log.New(multiWriter, "WARN:\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(multiWriter, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
}
