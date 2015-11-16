package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"text/tabwriter"

	"gopkg.in/natefinch/lumberjack.v2"
)

//Options provide program-wide options. At maximum, we got one persistent global and one short-living copy for writing in config file
type Options struct {
	ImageDir  string `long:"dir" description:"Target Directory" default:"img" ini-name:"downdir"`
	QDepth    int    `short:"q" long:"queue" description:"Length of the queue buffer" default:"20" ini-name:"queue_depth"`
	Tag       string `short:"t" long:"tag" description:"Tag to download"`
	Key       string `short:"k" long:"key" description:"Derpibooru API key" ini-name:"key"`
	StartPage int    `short:"p" long:"startpage" description:"Starting page for search" default:"1"`
	StopPage  int    `short:"n" long:"stoppage" description:"Stopping page for search, default - parse all search pages"`
	Filter    bool   `short:"f" long:"filter" description:"If set, enables client-side filtering of downloaded images"`
	Score     int    `long:"score" description:"Filter option, minimal score of image for it to be downloaded"`
	Unsafe    bool   `long:"unsafe" description:"If set, trusts in unknown authority"`
	NoHTTPS   bool   `long:"nohttps" description:"Disable HTTPS and try to download insecurely"`
	Args      struct {
		IDs []int `description:"Image IDs to download" optional:"yes"`
	} `positional-args:"yes"`
}

var opts Options
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

//WriteConfig writes default, presumably sensible configuration into file.
func WriteConfig(iniopts Options) {

	if opts.compareStatic(&iniopts) { //If nothing to write, no double-writing files
		return
	}

	config, err := os.OpenFile("config.ini", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	if err != nil {
		lFatal("Could  not create configuration file")
		//Need to check if log file is created before config file. Also, work around screams
		//in case if nor log, nor config file could be created.
	}

	defer func() {
		err = config.Close()
		if err != nil {
			lFatal("Could  not close configuration file")
		}
	}()

	tb := tabwriter.NewWriter(config, 10, 8, 0, ' ', 0) //Tabs! Elastic! Pretty!
	fmt.Fprintf(tb, "key \t= %s\n", opts.Key)
	fmt.Fprintf(tb, "queue_depth \t= %s\n", strconv.Itoa(opts.QDepth))
	fmt.Fprintf(tb, "downdir \t= %s\n", opts.ImageDir)

	err = tb.Flush()

	if err != nil {
		lFatal("Could  not write in configuration file")
	}
}

//compareStatic compares only options I want to preserve across launches.
func (a *Options) compareStatic(b *Options) bool {
	if a.ImageDir == b.ImageDir &&
		a.QDepth == b.QDepth &&
		a.Key == b.Key {
		return true
	}
	return false
}
