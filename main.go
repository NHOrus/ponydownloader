package main

import (
	"fmt"
	"os"
	"time"

	"github.com/inconshreveable/mousetrap"

	flag "github.com/jessevdk/go-flags"
)

//Default global variables
var (
	done   chan bool
	prefix = "https:"
)

func init() {
	if mousetrap.StartedByExplorer() {
		fmt.Println("Don't double-click ponydownloader")
		fmt.Println("You need to open cmd.exe and run it from the command line!")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	done = make(chan bool)
}

func main() {
	fmt.Println("Derpibooru.org Downloader version 0.5.3")

	SetLog() //Setting up logging of errors
	lInfo("Program start")

	args, iniopts := configSetup(&opts)

	WriteConfig(iniopts)

	if len(args) != 0 {
		lErr("Too many arguments, skipping following:", args)

	}

	if len(opts.Args.IDs) == 0 && opts.Tag == "" { //If no arguments after flags and empty/unchanged tag, what we should download? Sane end of line.
		lDone("Nothing to download, bye!") //Need to reshuffle flow: now it could end before it starts.
	}

	if opts.NoHTTPS {
		prefix = "http:" //Horrible cludge that must be removed in favor of url.URL.Scheme
	}

	//Creating directory for downloads if it does not yet exist
	err := os.MkdirAll(opts.ImageDir, 0755)

	if err != nil { //Execute bit means different thing for directories that for files. And I was stupid.
		lFatal(err) //We can not create folder for images, end of line.
	}

	//	Creating channels to pass info to downloader and to signal job well done
	imgdat := make(ImageCh, opts.QDepth) //Better leave default queue depth. Experiment shown that depth about 20 provides optimal perfomance on my system

	if opts.Tag == "" { //Because we can put imgid with flags. Why not?

		lInfo("Processing image No", opts.Args.IDs)
		go imgdat.ParseImg() // Sending imgid to parser. Here validity is our problem

	} else {

		// And here we send tags to getter/parser. Query and JSON validity is mostly server problem
		// Server response validity is ours
		lInfo("Processing tags", opts.Tag)
		go imgdat.ParseTag()
	}

	lInfo("Starting worker") //It would be funny if worker goroutine does not start

	filterInit(opts)                    //Ining filters based on our given flags
	filtimgdat := FilterChannel(imgdat) //see to move it into filter.Filter(inchan, outchan) where all filtration is done

	go filtimgdat.DlImg()

	<-done
	lDone("Finished")
	//And we are done here! Hooray!
}

func configSetup(*Options) ([]string, Options) {
	err := flag.IniParse("config.ini", &opts)
	if err != nil {
		switch err.(type) {
		default:
			panic(err)
		case *os.PathError:
			lWarn("config.ini not found, using defaults")
		}
	}
	var iniopts = opts

	args, err := flag.Parse(&opts)
	if err != nil {
		flagError := err.(*flag.Error)

		switch flagError.Type {
		case flag.ErrHelp:
			fallthrough
		case flag.ErrUnknownFlag:
			fmt.Println("Use --help to view all available options")
			os.Exit(0)
		default:
			lErr("Can't parse flags: %s\n", err)
			os.Exit(1)
		}
	}
	return args, iniopts
}
