package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	doneLogger *log.Logger
	infoLogger *log.Logger
	errLogger  *log.Logger
	warnLogger *log.Logger
)

//SetLog sets up logfile as I want it to: Copy to event.log, copy to command line
//Sometimes you just looks at available packages and feels that you must roll out your own solution
func init() {

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
//lInfo logs generic necessary program flow
func lInfo(v ...interface{}) {
	infoLogger.Println(v...)
}

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
	os.Exit(0)
}

//lErr notes non-fatal error and usually continues trying to crunch on
func lErr(v ...interface{}) {
	_ = errLogger.Output(2, fmt.Sprintln(v...)) //Following log package, ignoring error value
}

//lFatal happens when suffer some kind of error and we can't recover
func lFatal(v ...interface{}) {
	_ = errLogger.Output(2, fmt.Sprintln(v...)) //Following log package, ignoring error value
	os.Exit(1)
}

//lWarn is when there is no noticeable error, but something suspicious still happed

//prettifying return, so brackets will go away
func debracket(slice []int) string {
	stringSlice := make([]string, len(slice))
	for idx, num := range slice {
		stringSlice[idx] = strconv.Itoa(num)
	}
	return strings.Join(stringSlice, ", ")
}

func lWarn(v ...interface{}) {
	_ = warnLogger.Output(2, fmt.Sprintln(v...)) ////Following log package, ignoring error value
}
