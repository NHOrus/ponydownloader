package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"ponydownloader/settings"
	"strconv"

	"github.com/vaughan0/go-ini"
)

//	default variables
var (
	QDEPTH    int64       = 20    //Depth of the queue buffer - how many images are enqueued
	IMGDIR    string      = "img" //Default download directory
	TAG       string      = ""    //Default tag string is empty, it should be extracted from command line and only command line
	STARTPAGE int         = 1     //Default start page, derpiboo.ru 1-indexed
	STOPPAGE  int         = 0     //Default stop page, would stop parsing json when stop page is reached or site reaches the end of search
	elog      *log.Logger         //The logger for errors
)

type Image struct {
	imgid    int
	url      string
	filename string
	//	hash     string
}

func main() {

	fmt.Println("Derpiboo.ru Downloader version 0.2.0")

	elog, logfile := settings.SetLog() //setting up logging of errors

	defer logfile.Close() //Almost forgot. Always close the file in the end.

	config, err := ini.LoadFile("config.ini") // Loading default config file and checking for various errors.

	if os.IsNotExist(err) {
		elog.Fatalln("Config.ini does not exist, create it") //We can not live without config. We could, in theory, but writing default config if none exist can wait
	}

	if err != nil {
		elog.Panicln(err) //Oh, something is broken beyond my understanding. Sorry.
	}

	//Getting stuff from config, overwriting hardwired defaults when needed

	key, ok := config.Get("main", "key")

	if !ok || key == "" {
		elog.Println("'key' variable missing from 'main' section. It is vital for server-side filtering") //Empty key or key does not exist. Derpibooru works with this, but default image filter filters too much. Use key to set your own!
	}

	Q_temp, _ := config.Get("main", "workers")

	if Q_temp != "" {
		QDEPTH, err = strconv.ParseInt(Q_temp, 10, 0)

		if err != nil {
			elog.Fatalln("Wrong configuration: Depth of the buffer queue is not a number")

		}
	}

	ID_temp, _ := config.Get("main", "downdir")

	if ID_temp != "" {
		IMGDIR = ID_temp
	}

	//Here we are parsing all the flags. Command line argument hold priority to config. Except for 'key'. API-key is config-only

	flag.StringVar(&TAG, "t", TAG, "Tags to download")
	flag.IntVar(&STARTPAGE, "p", STARTPAGE, "Starting page for search")
	flag.IntVar(&STOPPAGE, "sp", STOPPAGE, "Stopping page for search, 0 - parse all all search pages")

	flag.Parse()

	if flag.NArg() == 0 && TAG == "" { //If no arguments after flags and empty/unchanged tag, what we should download? Sane end of line.
		log.SetPrefix("Done at ") //We can not do this with elog!
		log.Println("Nothing to download, bye!")
		os.Exit(0)
	}

	//Creating directory for downloads if it does not yet exist
	if err := os.MkdirAll(IMGDIR, 0644); err != nil { //Execute? No need to execute any image. Also, all those other users can not do anything beyond enjoying our images.
		elog.Fatalln(err) //We can not create folder for images, end of line.
	}

	//	Creating channels to pass info to downloader and to signal job well done
	imgdat := make(chan Image, QDEPTH) //Better leave default queue depth. Experiment shown that depth about 20 provides optimal perfomance on my system
	done := make(chan bool)

	if TAG == "" { //Because we can put imgid with flags. Why not?

		//	Checking argument for being a number and then getting image data

		imgid := flag.Arg(0) //0-indexed, unlike os.Args. os.Args[0] is path to program. It needs to be used later, when we are searching for what directory we are writing in
		_, err = strconv.Atoi(imgid)

		if err != nil {
			elog.Fatalln("Wrong input: can not parse", imgid, "as a number")
		}

		log.Println("Processing image No", imgid)

		go ParseImg(imgdat, imgid, key) // Sending imgid to parser. Here validity is our problem

	} else {

		//	and here we send tags to getter/parser. Validity is server problem, mostly

		log.Println("Processing tags", TAG)
		go ParseTag(imgdat, TAG, key)
	}

	log.Println("Starting worker") //It would be funny if worker goroutine does not start
	go DlImg(imgdat, done)

	<-done
	log.SetPrefix("Done at ")
	log.Println("Finised")
	//And we are done here! Hooray!
	return
}

func ParseImg(imgchan chan<- Image, imgid string, key string) {

	source := "http://derpiboo.ru/images/" + imgid + ".json?nofav=&nocomments="
	if key != "" {
		source = source + "&key=" + key
	}

	fmt.Println("Getting image info at:", source)

	resp, err := http.Get(source) //Getting our nice http response. Needs checking for 404 and other responses that are... less expected
	if err != nil {
		elog.Println(err)
		return
	}

	defer resp.Body.Close() //and not forgetting to close it when it's done

	var dat map[string]interface{}

	body, err := ioutil.ReadAll(resp.Body) //stolen from official documentation
	if err != nil {
		elog.Println(err)
		return
	}

	if err := json.Unmarshal(body, &dat); //transforming json into native map

	err != nil {
		elog.Println(err)
		return
	}

	InfoToChannel(dat, imgchan)

	close(imgchan) //closing channel, we are done here

	return
}

func DlImg(imgchan <-chan Image, done chan bool) {

	fmt.Println("Worker started; reading channel") //nice notification that we are not forgotten

	for {

		imgdata, more := <-imgchan

		if more { //checking that there is an image in channel

			if imgdata.filename == "" {
				elog.Println("Empty filename. Oops?") //something somewhere had gone wrong, going to the next image
			} else {

				fmt.Println("Saving as", imgdata.filename)

				func() { // to not hold all the files open when there is no need

					output, err := os.Create(IMGDIR + string(os.PathSeparator) + imgdata.filename) //And now, THE FILE!
					if err != err {
						elog.Println("Error when creating file for image" + strconv.Itoa(imgdata.imgid))
						elog.Println(err)
						return
					}
					defer output.Close() //Not forgetting to deal with it after completing download

					response, err := http.Get(imgdata.url)
					if err != nil {
						elog.Println("Error when getting image", imgdata.imgid)
						elog.Println(err)
						return
					}
					defer response.Body.Close() //Same, we shall not listen to the void when we finished getting image

					io.Copy(output, response.Body) //	Writing things we got from Derpibooru into the file and into hasher

				}()
			}

			//fmt.Println("\n", hex.EncodeToString(hash.Sum(nil)), "\n", imgdata.hash )

		} else {
			done <- true //well, there is no images in channel, it means we got them all, so synchronization is kicking in and ending the process
			break        //Just in case, so it would not stupidly die when program finishes - it will die smartly

		}
	}
}

func ParseTag(imgchan chan<- Image, tag string, key string) {

	source := "http://derpiboo.ru/search.json?nofav=&nocomments=" //yay hardwiring url strings!

	if key != "" {
		source = source + "&key=" + key
	}

	fmt.Println("Searching as", source+"&q="+tag)
	var working bool = true
	i := STARTPAGE
	for working {
		func() {
			fmt.Println("Searching page", i)
			resp, err := http.Get(source + "&q=" + tag + "&page=" + strconv.Itoa(i)) //Getting our nice http response. Needs checking for 404 and other responses that are... less expected
			defer resp.Body.Close()                                                  //and not forgetting to close it when it's done. And before we panic and die horribly.
			if err != nil {
				elog.Println("Error while getting search page", i)
				elog.Println(err)
				return
			}

			var dats []map[string]interface{} //Because we got array incoming instead of single object, we using an slive of maps!

			//fmt.Println(resp)

			body, err := ioutil.ReadAll(resp.Body) //stolen from official documentation
			if err != nil {
				elog.Println("Error while reading search page", i)
				elog.Println(err)
				return
			}

			//fmt.Println(body)

			if err := json.Unmarshal(body, &dats); //transforming json into native slice of maps

			err != nil {
				elog.Println("Error while parsing search page", i)
				elog.Println(err)
				return

			}

			if len(dats) == 0 {
				fmt.Println("Pages are over")
				working = false
				return
			} //exit due to finishing all pages

			for _, dat := range dats {
				InfoToChannel(dat, imgchan)
			}
			if i == STOPPAGE {
				working = false
				return
			}
			i++

		}()
	}

	close(imgchan)
}

func InfoToChannel(dat map[string]interface{}, imgchan chan<- Image) {

	var imgdata Image

	imgdata.url = "http:" + dat["image"].(string)
	//	imgdata.hash = dat["sha512_hash"].(string)
	imgdata.filename = (strconv.FormatFloat(dat["id_number"].(float64), 'f', -1, 64) + "." + dat["file_name"].(string) + "." + dat["original_format"].(string))
	imgdata.imgid = int(dat["id_number"].(float64))

	//	for troubleshooting - possibly debug flag?
	//	fmt.Println(dat)
	//	fmt.Println(imgdata.url)
	//	fmt.Println(imgdata.hash)
	//	fmt.Println(imgdata.filename)

	imgchan <- imgdata
}
