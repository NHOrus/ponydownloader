package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/NHOrus/ponydownloader/derpiapi" //Things we do with images and stuff
	"github.com/NHOrus/ponydownloader/settings" //Here we are working with setting things up or down, depending.
)

//Default hardcoded variables
var (
	QDEPTH     int        = 20    //Depth of the queue buffer - how many images are enqueued
	IMGDIR                = "img" //Default download directory
	TAG        string             //Default tag string is empty, it should be extracted from command line and only command line
	STARTPAGE  = 1                //Default start page, derpiboo.ru 1-indexed
	STOPPAGE   = 0                //Default stop page, would stop parsing json when stop page is reached or site reaches the end of search
	elog       log.Logger         //The logger for errors
	KEY        string     = ""    //Default identification key. Get your own and place it in configuration, people
	SCRFILTER  int                //So we can ignore things with limited
	FILTERFLAG = false            //Gah, not sure how to make it better.
)

func init() {

	Set := settings.Settings{QDepth: QDEPTH, ImgDir: IMGDIR, Key: KEY}

	Set.GetConfig(elog)

	QDEPTH = Set.QDepth
	KEY = Set.Key
	IMGDIR = Set.ImgDir

	//Here we are parsing all the flags. Command line argument hold priority to config.
	flag.StringVar(&TAG, "t", TAG, "Tags to download")
	flag.IntVar(&STARTPAGE, "p", STARTPAGE, "Starting page for search")
	flag.IntVar(&STOPPAGE, "np", STOPPAGE, "Stopping page for search, 0 - parse all all search pages")
	flag.StringVar(&KEY, "k", KEY, "Your key to derpibooru API")
	flag.IntVar(&SCRFILTER, "scr", SCRFILTER, "Minimal score of image for it to be downloaded")
	flag.BoolVar(&FILTERFLAG, "filter", FILTERFLAG, "If set (to true), enables client-side filtration of downloaded images")

	flag.Parse()

}

func main() {

	fmt.Println("Derpiboo.ru Downloader version 0.2.0")

	elog, logfile := settings.SetLog() //Setting up logging of errors

	defer logfile.Close() //Almost forgot. Always close the file in the end.

	if flag.NArg() == 0 && TAG == "" { //If no arguments after flags and empty/unchanged tag, what we should download? Sane end of line.

		log.SetPrefix("Done at ")                //We can not do this with elog!
		log.Println("Nothing to download, bye!") //Need to reshuffle flow: now it could end before it starts.
		os.Exit(0)
	}

	//Creating directory for downloads if it does not yet exist
	err := os.MkdirAll(IMGDIR, 0755)

	if err != nil { //Execute bit means different thing for directories that for files. And I was stupid.
		elog.Fatalln(err) //We can not create folder for images, end of line.
	}

	//	Creating channels to pass info to downloader and to signal job well done
	imgdat := make(chan derpiapi.Image, QDEPTH) //Better leave default queue depth. Experiment shown that depth about 20 provides optimal perfomance on my system
	done := make(chan bool)

	if TAG == "" { //Because we can put imgid with flags. Why not?

		//	Checking argument for being a number and then getting image data

		imgid := flag.Arg(0) //0-indexed, unlike os.Args. os.Args[0] is path to program. It needs to be used later, when we are searching for what directory we are writing in
		_, err := strconv.Atoi(imgid)

		if err != nil {
			elog.Fatalln("Wrong input: can not parse ", imgid, "as a number")
		}

		log.Println("Processing image No", imgid)

		go derpiapi.ParseImg(imgdat, imgid, KEY, elog) // Sending imgid to parser. Here validity is our problem

	} else {

		//	and here we send tags to getter/parser. Validity is server problem, mostly

		log.Println("Processing tags", TAG)
		go derpiapi.ParseTag(imgdat, TAG, KEY, STARTPAGE, STOPPAGE, elog)
	}

	log.Println("Starting worker") //It would be funny if worker goroutine does not start

	filtimgdat := make(chan derpiapi.Image)
	fflag := derpiapi.FilterSet{Scrfilter: SCRFILTER, Filterflag: FILTERFLAG}

	go derpiapi.FilterChannel(imgdat, filtimgdat, fflag) //see to move it into filter.Filter(inchan, outchan) where all filtration is done

	go derpiapi.DlImg(filtimgdat, done, elog, IMGDIR)

	<-done
	log.SetPrefix("Done at ")
	log.Println("Finished")
	//And we are done here! Hooray!
	return
}
