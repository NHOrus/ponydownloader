package main

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	doneLogger *log.Logger
	infoLogger *log.Logger
	errLogger  *log.Logger
	warnLogger *log.Logger
)

//SetLog sets up logfile as I want it to: Copy to event.log, copy to commandline
//Sometimes you just looks at wider net and feels that you must roll out your own solution
func SetLog() {

	logfile := &lumberjack.Logger{
		Filename:   "event.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}

	infoLog := io.MultiWriter(logfile, os.Stdout)
	errLog := io.MultiWriter(logfile, os.Stderr)

	doneLogger = log.New(infoLog, "Done at ", log.LstdFlags)
	infoLogger = log.New(infoLog, "Happened at ", log.LstdFlags)
	warnLogger = log.New(errLog, "Warning at ", log.LstdFlags|log.Lshortfile)
	errLogger = log.New(errLog, "Error at ", log.LstdFlags|log.Lshortfile) //Setting stuff for our logging: both errors and events.
}

//Wrappers for loggers to simplify invocation and don't suffer premade packages
func lInfo(v ...interface{}) {
	infoLogger.Println(v...)
}

func lInfof(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

func lDone(v ...interface{}) {
	doneLogger.Println(v...)
	os.Exit(0)
}

func lErr(v ...interface{}) {
	errLogger.Println(v...)
}

func lFatal(v ...interface{}) {
	errLogger.Fatalln(v...)
}

func lWarn(v ...interface{}) {
	warnLogger.Println(v...)
}
