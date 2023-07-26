/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package logger

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Fatal *log.Logger

	logFile *os.File
)

func Open(file string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}

	flags := log.Ldate | log.Ltime | log.Lshortfile
	Info = log.New(f, "INFO: ", flags)
	Warn = log.New(f, "WARN: ", flags)
	Error = log.New(f, "ERROR: ", flags)
	Fatal = log.New(f, "FATAL: ", flags)

	logFile = f
}

func Close() {
	logFile.Close()
}

// func Write() {

// }
