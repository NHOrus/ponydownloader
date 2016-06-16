//Command ponydownloader uses Derpibooru.org API to download pony images
//by ID or by tags, with some client-side filtration ability
package main

import (
	"fmt"
	"os"
	"os/signal"
)

//Default global variables
var (
	scheme         = "https:"
	interruptParse = make(chan os.Signal, 1)
	interruptDL    = make(chan os.Signal, 1)
)

func init() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		interruptParse <- os.Interrupt
		close(interruptDL)
		<-sig
		lDone("Program interrupted by user's command")
	}()
}

func main() {
	fmt.Println("Derpibooru.org Downloader version 0.9.0")

	opts, lostArgs := getOptions()

	lInfo("Program start")
	// Checking for extra arguments we got no idea what to do with
	if len(lostArgs) != 0 {
		lErr("Too many arguments, skipping following:", lostArgs)
	}
	//If no arguments after flags and empty/unchanged tag, what we should download? Sane end of line.
	if len(opts.Args.IDs) == 0 && opts.Tag == "" {
		lDone("Nothing to download, bye!")
	}

	if opts.UnsafeHTTPS {
		makeHTTPSUnsafe()
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
			lInfo("Processing image №", opts.Args.IDs[0])
		} else {
			lInfo("Processing images №", debracket(opts.Args.IDs))
		}
		go imgdat.ParseImg(opts.Args.IDs, opts.Key) // Sending Image ID to parser. Here validity is our problem

	} else {

		// And here we send tags to getter/parser. Query and JSON validity is mostly server problem
		// Server response validity is ours
		lInfo("Processing tags", opts.Tag)
		go imgdat.ParseTag(opts.TagOpts, opts.Key)
	}

	lInfo("Starting worker") //It would be funny if worker goroutine does not start

	filterInit(opts.FiltOpts, bool(opts.Config.LogFilters)) //Initiating filters based on our given flags
	filtimgdat := FilterChannel(imgdat)                     //Actual filtration

	filtimgdat.interrupt().downloadImages(opts.Config) // Now that we got asynchronous list of images we want to get done, we can get them.

	lDone("Finished")
	//And we are done here! Hooray!
}

func (imgchan ImageCh) interrupt() (outch ImageCh) {
	outch = make(ImageCh)
	go func() {
		for {
			select {
			case <-interruptDL:
				close(outch)
				return
			case img, ok := <-imgchan:
				if !ok {
					close(outch)
					imgchan = nil
					return
				}
				outch <- img
			}
		}
	}()

	return outch
}
