package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vaughan0/go-ini"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"log"
	"crypto/sha512"
	"encoding/hex"

	//	"errors"	
	//	"net"
)

//	defaults:
var (
	WORKERS   int64  = 10    //Number of workers
	IMGDIR    string = "img" //default download directory
	TAG       string = ""    //default string is empty, it can only ge extracted from command line
	STARTPAGE int    = 1     //default start page, derpiboo.ru 1-indexed
	STOPPAGE  int    = 0     //default stop page, would stop parsing json when stop page is reached or site reaches the end of search
)

func main() {
	
	fmt.Println("Derpiboo.ru Downloader version 0.1.3")
	
	logfile, err := os.OpenFile("event.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644 ) //file for putting errors into
		if err != nil { panic(err) }
	defer logfile.Close()	//Almost forgot. Always close the file in the end.
	
	elog := log.New(io.MultiWriter(logfile, os.Stderr), "Errors at ", log.LstdFlags) //setting stuff for our logging: both errors and events.
	log.SetPrefix("Happens at ")
	log.SetFlags(log.LstdFlags)
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	log.Println("Program start")
	
	config, err := ini.LoadFile("config.ini") // Loading default config file and checking for various errors.

	if os.IsNotExist(err) {
		elog.Fatalln("Config.ini does not exist, create it")
	}

	if err != nil {
		elog.Panicln(err)
	}

	//Getting stuff from config, overwriting defaults

	key, ok := config.Get("main", "key")

	if !ok {
		elog.Println("'key' variable missing from 'main' section")
	}

	W_temp, _ := config.Get("main", "workers")

	if W_temp != "" {
		WORKERS, err = strconv.ParseInt(W_temp, 10, 0)

		if err != nil {
			elog.Fatalln("Wrong configuration: Amount of workers is not a number")
			
		}
	}

	ID_temp, _ := config.Get("main", "downdir")

	if ID_temp != "" {
		IMGDIR = ID_temp
	}

	//Here we are parsing all the flags

	flag.StringVar(&TAG, "t", TAG, "Tags to download")
		flag.Parse()

	if flag.NArg() == 0 && TAG == "" {	//If no arguments after flags and empty/unchanged tag, what we should download? Sane end of line.
		elog.SetPrefix("Done at ")
//		elog.SetOutput(io.MultiWriter(logfile, os.Stdout))
		elog.Println("Nothing to download, bye!")
		os.Exit(0)
	}

	//	creating directory for downloads if not yet done
	if err := os.MkdirAll(IMGDIR, 0777); err != nil {
		elog.Fatalln(err)
	}

	//	creating channels to pass info to downloader and to signal job well done
	imgdat := make(chan Image, WORKERS)
	done := make(chan bool)

	if TAG == "" {

		//	checking argument for being a number and then getting image data

		imgid := flag.Arg(0)
		_, err = strconv.Atoi(imgid)

		if err != nil {
			elog.Fatalln("Wrong input: can not parse", imgid, "as a number")
		}

		log.Println("Processing image No", imgid)

		go parseImg(imgdat, imgid, key)

	} else {

		//	and here we send tags to getter/parser

		log.Println("Processing tags", TAG)
		go parseTag(imgdat, TAG, key)
	}

	log.Println("Starting worker")
	go dlimage(imgdat, done)

	<-done
	log.Println("Finised")

}

type Image struct {
	url      string
	filename string
	hash     string	
	}

func parseImg(imgchan chan<- Image, imgid string, key string) {

	source := "http://derpiboo.ru/" + imgid + ".json?nofav=&nocomments="
	if key != "" {
		source = source + "&key=" + key
	}

	fmt.Println("Getting image info at:", source)

	resp, err := http.Get(source) //Getting our nice http response. Needs checking for 404 and other responses that are... less expected
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close() //and not forgetting to close it when it's done

	var dat map[string]interface{}

	body, err := ioutil.ReadAll(resp.Body) //stolen from official documentation
	if err != nil {
		panic(err)
	}

	//fmt.Println(body)

	if err := json.Unmarshal(body, &dat); //transforming json into native map

	err != nil {
		panic(err)
	}

	InfoToChannel(dat, imgchan)

	close(imgchan) //closing channel, we are done here

	return
}

func dlimage(imgchan <-chan Image, done chan bool) {

	fmt.Println("Worker started; reading channel") //nice notification that we are not forgotten

	for {

		imgdata, more := <-imgchan

		if more { //checking that there is an image in channel

			if imgdata.filename == "" {
				fmt.Println("Empty filename. Oops?") //something somewhere had gone wrong, off with the worker
				break
			}

			fmt.Println("Saving as", imgdata.filename)

			func() { // to not hold all the files open when there is no need
				output, err := os.Create(IMGDIR + string(os.PathSeparator) + imgdata.filename) //And now, THE FILE!
				if err != err {
					panic(err)
				}
				defer output.Close() //Not forgetting to deal with it after completing download

				response, err := http.Get(imgdata.url)
				if err != nil {
					fmt.Println("Error while downloading", imgdata.url, "-", err)
					panic(err)
					return
				}
				defer response.Body.Close() //Same, we shall not listen to the void when we finished getting image

				hash := sha512.New()
				
				io.Copy(io.MultiWriter(output, hash), response.Body) //	Writing things we got from Derpibooru into the 

				b := make([]byte, hash.Size())
				hash.Sum(b[:0])

				//	fmt.Println("\n", hex.EncodeToString(b), "\n", imgdata.hash )

				if hex.EncodeToString(b) == imgdata.hash {
					fmt.Println("Hash correct")
				}	else {
					fmt.Println("Hash wrong")
				}

				//fmt.Println("\n", hex.EncodeToString(hash.Sum(nil)), "\n", imgdata.hash )
			
			}()
			
		} else {
			done <- true //well, there is no images in channel, it means we got them all, so synchronization is kicking in and ending the process

		}

	}
}

func parseTag(imgchan chan<- Image, tag string, key string) {

	source := "http://derpiboo.ru/search.json?nofav=&nocomments="

	if key != "" {
		source = source + "&key=" + key
	}

	fmt.Println("Searching as", source+"&q="+tag)
	var i int = 1

	for {
		fmt.Println("Searching page", i)
		resp, err := http.Get(source + "&q=" + tag + "&page=" + strconv.Itoa(i)) //Getting our nice http response. Needs checking for 404 and other responses that are... less expected
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close() //and not forgetting to close it when it's done

		var dats []map[string]interface{} //Because we got array incoming instead of single object, we using an slive of maps!

		//fmt.Println(resp)

		body, err := ioutil.ReadAll(resp.Body) //stolen from official documentation
		if err != nil {
			panic(err)
		}

		//fmt.Println(body)

		if err := json.Unmarshal(body, &dats); //transforming json into native slice of maps

		err != nil {
			panic(err)

		}

		if len(dats) == 0 {
			fmt.Println("Pages are over")
			break
		} //exit due to finishing all pages

		for _, dat := range dats {
			InfoToChannel(dat, imgchan)
		}
		
		i++
	}
	close(imgchan)
}

func InfoToChannel(dat map[string]interface{}, imgchan chan<- Image){
	
	var imgdata Image
	
	imgdata.url = "http:" + dat["image"].(string)
	imgdata.hash = dat["sha512_hash"].(string)
	imgdata.filename = strconv.FormatFloat(dat["id_number"].(float64), 'f', -1, 64) + "." + dat["file_name"].(string) + "." + dat["original_format"].(string)

	//	for troubleshooting later
	//	fmt.Println(dat)
	//	fmt.Println(imgdata.url)
	//	fmt.Println(imgdata.hash)
	//	fmt.Println(imgdata.filename)

	imgchan <- imgdata
}
