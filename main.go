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

	//	"errors"
	//	"log"
	//	"net"
	//	"crypto/sha512"
	//	"encoding/hex"
)

//	defaults:
var (
	WORKERS 	int64	= 10    //Number of workers
	IMGDIR  	string	= "img" //default download directory
	TAG     	string	= ""    //default string is empty, it can only ge extracted from command line
	STARTPAGE	int		= 1		//default start page, derpiboo.ru 1-indexed
	STOPPAGE	int		= 0		//default stop page, would stop parsing json when it ends
	)

func main() {

	config, err := ini.LoadFile("config.ini") // Loading default config file and checking for various errors.
	if os.IsNotExist(err) {
		panic("config.ini does not exist, create it")
	}

	if err != nil {
		panic(err)
	}

	//Getting stuff from config, overwriting defaults

	key, ok := config.Get("main", "key")
	if !ok {
		panic("'key' variable missing from 'main' section")
	}

	W_temp, _ := config.Get("main", "workers")
	if W_temp != "" {
		WORKERS, err = strconv.ParseInt(W_temp, 10, 0)
		if err != nil {
			fmt.Println("Wrong configuration: Amount of workers is not a number")
			os.Exit(1)
		}
	}

	ID_temp, _ := config.Get("main", "downdir")
	if ID_temp != "" {
		IMGDIR = ID_temp
	}

	//here shall be flag parser

	flag.StringVar(&TAG, "t", TAG, "Tag to download: Replace spaces with \"+\".")
	flag.Parse()
	
	fmt.Println("Derpiboo.ru Downloader version 0.0.6 \nWorking")


	length := flag.NArg()
	if length == 0 && TAG == "" {
		fmt.Println("Nothing to download, bye!")
		os.Exit(0)
	}

	//	creating directory for downloads if not yet done
	if err := os.MkdirAll(IMGDIR, 0777); err != nil {
		panic(err)
	}

	imgdat := make(chan Image, WORKERS)
	done := make(chan bool)
	
	if TAG == "" {

		imgid := flag.Arg(length - 1)

		//	fmt.Println(key) //Just checking that I am not wrong

		_, err = strconv.ParseInt(imgid, 10, 0)
		if err != nil {
			fmt.Println("Wrong input: can not parse", imgid, "as a number")
			os.Exit(1)
		}

		fmt.Println("Processing image No " + imgid)

		go parseImg(imgdat, imgid, key)

	}	else	{
	//	and here we parse tag and abuse stuff horribly
		
		fmt.Println("Trying to process tags")
		go parseTag(imgdat, TAG, key)
	}
	
	fmt.Println("Starting worker")
	go dlimage(imgdat, done)
	
	<-done

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

		fmt.Println(source)

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
	var imgdata Image
	imgdata.url = "http:" + dat["image"].(string)
	imgdata.hash = dat["sha512_hash"].(string) //for the future and checking that we got file right
	imgdata.filename = strconv.FormatFloat(dat["id_number"].(float64), 'f', -1, 64) + "." + dat["file_name"].(string) + "." + dat["original_format"].(string)

	//	fmt.Println(strconv.FormatFloat(dat["id_number"].(float64), 'f', -1, 64))

	//	fmt.Println(dat)

	//	for now and troubleshooting
	//	fmt.Println(imgdata.url)
	//	fmt.Println(imgdata.hash)
	//	fmt.Println(imgdata.filename)

	imgchan <- imgdata
	
	close(imgchan)

	return
}

func dlimage(imgchan <-chan Image, done chan bool) {
	//	fmt.Println("reading channel")
	
	for {

	imgdata, more := <-imgchan
	
	if	more {
	if imgdata.filename == "" { fmt.Println("Empty filename. Oops?"); break }
	
	fmt.Println("Saving as ", imgdata.filename)
	PathSep, _ := strconv.Unquote(strconv.QuoteRune(os.PathSeparator))

	output, err := os.Create(IMGDIR + PathSep + imgdata.filename)
	if err != err {
		panic(err)
	}
	defer output.Close()

	response, err := http.Get(imgdata.url)
	if err != nil {
		panic(err)
		fmt.Println("Error while downloading", imgdata.url, "-", err)
		return
	}
	defer response.Body.Close()

	io.Copy(output, response.Body)

	/*	hash := sha512.New()
		io.Copy(hash, response.Body)
		b := make([]byte, hash.Size())
		hash.Sum(b[:0])

		fmt.Println("\n", hex.EncodeToString(b), "\n", imgdata.hash )

		if hex.EncodeToString(b) == imgdata.hash {
			fmt.Println("Hash correct")
		}	else {
			fmt.Println("Hash wrong")
		}

		fmt.Println("\n", hex.EncodeToString(hash.Sum(nil)), "\n", imgdata.hash )
	*/
		} else	{
			done <- true
		
		}
				
		}
}

func parseTag(imgchan chan<- Image, tag string, key string) {
	
	if tag == "" { fmt.Println("Something has gone horribly wrong, no tag found?"); os.Exit(1) }
	source := "http://derpiboo.ru/search.json?nofav=&nocomments=&utf8=false"
	if key != "" { source = source + "&key=" + key }
	
	fmt.Println(source + "&q=" + tag)
	
	resp, err := http.Get(source + "&q=" + tag) //Getting our nice http response. Needs checking for 404 and other responses that are... less expected
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close() //and not forgetting to close it when it's done

	var dats []map[string]interface{}
	
	//fmt.Println(resp)
	
	body, err := ioutil.ReadAll(resp.Body) //stolen from official documentation
	if err != nil {
		panic(err)
	}
	
	//fmt.Println(body)

	if err := json.Unmarshal(body, &dats); //transforming json into native map

	err != nil {
		panic(err)
	
	}
	
	var imgdata Image
	
	for _, dat := range dats {
		
		imgdata.url = "http:" + dat["image"].(string)
		imgdata.hash = dat["sha512_hash"].(string) //for the future and checking that we got file right
		imgdata.filename = strconv.FormatFloat(dat["id_number"].(float64), 'f', -1, 64) + "." + dat["file_name"].(string) + "." + dat["original_format"].(string)
		
		imgchan <- imgdata
	}	
	
	close(imgchan)
}
