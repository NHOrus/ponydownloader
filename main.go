package main

import (
	"fmt"
	"os"
	"time"

	"github.com/inconshreveable/mousetrap"
)

//Default global variables
var (
	prefix = "https:"
)

func init() {
	if mousetrap.StartedByExplorer() {
		fmt.Println("Don't double-click ponydownloader")
		fmt.Println("You need to open cmd.exe and run it from the command line!")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Derpibooru.org Downloader version 0.6.0")

	SetLog() //Setting up logging of errors

	opts, lostArgs := doOptions()

	lInfo("Program start")
	// Checking for extra arguments we got no idea what to do with
	if len(lostArgs) != 0 {
		lErr("Too many arguments, skipping following:", lostArgs)
	}
	//If no arguments after flags and empty/unchanged tag, what we should download? Sane end of line.
	if len(opts.Args.IDs) == 0 && opts.Tag == "" {
		lDone("Nothing to download, bye!")
	}

	if opts.NoHTTPS {
		prefix = "http:" //Horrible kludge that must be removed in favor of url.URL.Scheme
	}

	//Creating directory for downloads if it does not yet exist
	err := os.MkdirAll(opts.ImageDir, 0755)

	if err != nil { //Execute bit means different thing for directories that for files. And I was stupid.
		lFatal(err) //We can not create folder for images, end of line.
	}

	//	Creating channels to pass info to downloader and to signal job well done
	imgdat := make(ImageCh, opts.QDepth) //Better leave default queue depth. Experiment shown that depth about 20 provides optimal performance on my system

	if opts.Tag == "" { //Because we can put Image ID with flags. Why not?

		if len(opts.Args.IDs) == 1 {
			lInfo("Processing image №", opts.Args.IDs)
		} else {
			lInfo("Processing images №", opts.Args.IDs)
		}
		go imgdat.ParseImg(opts.Args.IDs, opts.Key, opts.Unsafe) // Sending Image ID to parser. Here validity is our problem

	} else {

		// And here we send tags to getter/parser. Query and JSON validity is mostly server problem
		// Server response validity is ours
		lInfo("Processing tags", opts.Tag)
		go imgdat.ParseTag(opts)
	}

	lInfo("Starting worker") //It would be funny if worker goroutine does not start

	filterInit(opts)                    //Initiating filters based on our given flags
	filtimgdat := FilterChannel(imgdat) //Actual filtration

	filtimgdat.DlImg(opts) // Now that we got asynchronous list of images we want to get done, we can get them.

	lDone("Finished")
	//And we are done here! Hooray!
}
