package main

import (
	"fmt"
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

//Setting up logfile as I want it to: Copy to event.log, copy to command line
//Sometimes you just look at available packages and feel that you must roll out your own solution
func init() {

	logfile := &lumberjack.Logger{
		Filename:   "event.log",
		MaxSize:    1, // megabytes
		MaxBackups: 9,
		MaxAge:     28, //days
	}

	infoLog := io.MultiWriter(logfile, os.Stdout)
	errLog := io.MultiWriter(logfile, os.Stderr)

	doneLogger = log.New(infoLog, "Done at ", log.LstdFlags)
	infoLogger = log.New(infoLog, "Happened at ", log.LstdFlags)
	warnLogger = log.New(errLog, "Warning at ", log.LstdFlags|log.Lshortfile)
	errLogger = log.New(errLog, "Error at ", log.LstdFlags|log.Lshortfile)
}

//Wrappers for loggers to simplify invocation and don't suffer premade packages
//lInfo logs generic necessary program flow
func lInfo(v ...interface{}) {
	infoLogger.Println(v...)
}

//lCondInfo doesn't log when it's disalbed
func lCondInfo(on bool, v ...interface{}) {
	if on {
		infoLogger.Println(v...)
	}
}

//lInfof logs generic program flow with ability to format string beyond defaults
//Used only to note downloading speed and timing
func lInfof(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

//lDone notes that we are finished and there is nothing left to do, sane way
func lDone(v ...interface{}) {
	doneLogger.Println(v...)
}

//lErr notes non-fatal error and usually continues trying to crunch on
func lErr(v ...interface{}) {
	/* #nosec */
	_ = errLogger.Output(2, fmt.Sprintln(v...)) //Following log package, ignoring error value
}

//lFatal happens when suffer some kind of error and we can't recover
func lFatal(v ...interface{}) {
	/* #nosec */
	_ = errLogger.Output(2, fmt.Sprintln(v...)) //Following log package, ignoring error value
	os.Exit(1)
}

//lWarn is when there is no noticeable error, but something suspicious still happed
func lWarn(v ...interface{}) {
	/* #nosec */
	_ = warnLogger.Output(2, fmt.Sprintln(v...)) //Following log package, ignoring error value
}
