package main

import (
	"fmt"
	"github.com/NHOrus/ponydownloader/derpiapi" //Things we do with images and stuff
	flag "github.com/jessevdk/go-flags"
	"log"
	"os"
	"strconv"
)

//Default hardcoded variables
var (
	QDEPTH     = 20       //Depth of the queue buffer - how many images are enqueued
	IMGDIR     = "img"    //Default download directory
	TAG        string     //Default tag string is empty, it should be extracted from command line and only command line
	STARTPAGE  = 1        //Default start page, derpiboo.ru 1-indexed
	STOPPAGE   = 0        //Default stop page, would stop parsing json when stop page is reached or site reaches the end of search
	elog       log.Logger //The logger for errors
	KEY        string     //Default identification key. Get your own and place it in configuration, people
	SCRFILTER  int        //So we can ignore things by the score
	FILTERFLAG = false    //Gah, not sure how to make it better.
)

func main() {
	fmt.Println("Derpiboo.ru Downloader version 0.4.0")

	Set := Settings{QDepth: QDEPTH, ImgDir: IMGDIR, Key: KEY}

	err := flag.IniParse("config.ini", &opts)
	if err != nil {
		panic(err)
	}

	_, err = flag.Parse(&opts)
	if err != nil {
		flagError := err.(*flag.Error)
		if flagError.Type == flag.ErrHelp {
			return
		}
		if flagError.Type == flag.ErrUnknownFlag {
			fmt.Println("Use --help to view all available options")
			return
		}

		fmt.Printf("Error parsing flags: %s\n", err)
		return
	}
	Set.GetConfig(elog)

	elog, logfile := SetLog() //Setting up logging of errors

	defer logfile.Close() //Almost forgot. Always close the file in the end.

	if len(os.Args) == 1 && TAG == "" { //If no arguments after flags and empty/unchanged tag, what we should download? Sane end of line.

		log.SetPrefix("Done at ")                //We can not do this with elog!
		log.Println("Nothing to download, bye!") //Need to reshuffle flow: now it could end before it starts.
		os.Exit(0)
	}

	//Creating directory for downloads if it does not yet exist
	err = os.MkdirAll(IMGDIR, 0755)

	if err != nil { //Execute bit means different thing for directories that for files. And I was stupid.
		elog.Fatalln(err) //We can not create folder for images, end of line.
	}

	//	Creating channels to pass info to downloader and to signal job well done
	imgdat := make(derpiapi.ImageCh, QDEPTH) //Better leave default queue depth. Experiment shown that depth about 20 provides optimal perfomance on my system
	done := make(chan bool)

	if TAG == "" { //Because we can put imgid with flags. Why not?

		//	Checking argument for being a number and then getting image data

		imgid := os.Args[1]
		_, err := strconv.Atoi(imgid)

		if err != nil {
			elog.Fatalln("Wrong input: can not parse ", imgid, "as a number")
		}

		log.Println("Processing image No", imgid)

		go imgdat.ParseImg(imgid, KEY, elog) // Sending imgid to parser. Here validity is our problem

	} else {

		//	and here we send tags to getter/parser. Validity is server problem, mostly

		log.Println("Processing tags", TAG)
		go imgdat.ParseTag(TAG, KEY, STARTPAGE, STOPPAGE, elog)
	}

	log.Println("Starting worker") //It would be funny if worker goroutine does not start

	filtimgdat := make(derpiapi.ImageCh)
	fflag := derpiapi.FilterSet{Scrfilter: SCRFILTER, Filterflag: FILTERFLAG}

	go derpiapi.FilterChannel(imgdat, filtimgdat, fflag) //see to move it into filter.Filter(inchan, outchan) where all filtration is done

	go filtimgdat.DlImg(done, elog, IMGDIR)

	<-done
	log.SetPrefix("Done at ")
	log.Println("Finished")
	//And we are done here! Hooray!
	return
}
