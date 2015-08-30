package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
)

var opts struct {
	ImageDir  string `long:"dir" description:"Target Directory" default:"img" ini-name:"downdir"`
	QDepth    int    `short:"q" long:"queue" description:"Length of the queue buffer" default:"20" ini-name:"queue_depth"`
	Tag       string `short:"t" long:"tag" description:"Tag to download"`
	Key       string `short:"k" long:"key" description:"Derpibooru API key" ini-name:"key"`
	StartPage int    `short:"p" long:"startpage" description:"Starting page for search" default:"1"`
	StopPage  int    `short:"n" long:"stoppage" description:"Stopping page for search, default - parse all search pages"`
	Filter    bool   `short:"f" long:"filter" description:"If set, enables client-side filtering of downloaded images"`
	Score     int    `long:"score" description:"Filter option, minimal score of image for it to be downloaded"`
	Args      struct {
		IDs []int `description:"Image IDs to download" optional:"yes"`
	} `positional-args:"yes"`
}

//SetLog sets up logfile as I want it to: Copy to event.log, copy to commandline
func SetLog() (retlog *log.Logger) {

	logfile, err := os.OpenFile("event.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644) //file for putting errors into

	if err != nil {
		panic(err)
	}
	//elog - recoverable errors. log - just things that happen
	retlog = log.New(io.MultiWriter(logfile, os.Stderr), "Errors at ", log.LstdFlags) //Setting stuff for our logging: both errors and events.

	log.SetPrefix("Happens at ")
	log.SetFlags(log.LstdFlags)
	log.SetOutput(io.MultiWriter(logfile, os.Stdout)) //we write in file and stdout
	log.Println("Program start")

	return
}

//WriteConfig writes default, presumably sensible configuration into file.
func WriteConfig() {
	config, err := os.OpenFile("config.ini", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer func() {
		err = config.Close()
		if err != nil {
			elog.Fatalln("Could  not close configuration file")
		}
	}()

	if err != nil {
		elog.Fatalln("Could  not create configuration file")
	}

	tb := tabwriter.NewWriter(config, 10, 8, 0, ' ', 0)
	fmt.Fprintf(tb, "key \t= %s\t\n", opts.Key)
	fmt.Fprintf(tb, "queue_depth \t= %s\t\n", strconv.Itoa(opts.QDepth))
	fmt.Fprintf(tb, "downdir \t= %s\t\n", opts.ImageDir)

	err = tb.Flush()

	if err != nil {
		elog.Fatalln("Could  not write in configuration file")
	}
}
