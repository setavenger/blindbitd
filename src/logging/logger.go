package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

func LoadLoggers(directoryPath string) {

	file, err := os.OpenFile(fmt.Sprintf("%s/logs.log", directoryPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fileDebug, err := os.OpenFile(fmt.Sprintf("%s/logs-debug.log", directoryPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multi := io.MultiWriter(file, os.Stdout)
	//multiDebug := io.MultiWriter(fileDebug, os.Stdout)

	DebugLogger = log.New(fileDebug, "[DEBUG] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	InfoLogger = log.New(multi, "[INFO] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	WarningLogger = log.New(multi, "[WARNING] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	ErrorLogger = log.New(multi, "[ERROR] ", log.Ldate|log.Lmicroseconds|log.Llongfile|log.Lmsgprefix)
}

func LoadLoggersMock() {
	DebugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	InfoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	WarningLogger = log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	ErrorLogger = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Lmicroseconds|log.Llongfile|log.Lmsgprefix)
}

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)
