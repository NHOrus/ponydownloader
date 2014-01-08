package main

import (
	"fmt"
	"os"
	"encoding/json"
	"github.com/vaughan0/go-ini"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"flag"
	
//	"errors"
//	"log"
//	"net"
//	"crypto/sha512"
//	"encoding/hex"

)

//	defaults:
var (
	WORKERS	int64	= 10	//Number of workers
	IMGDIR	string	= "img"	//default download directory
	)

func main() {

	fmt.Println("Derpiboo.ru Downloader version 0.0.4 \nWorking")

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
		WORKERS, err = strconv.ParseInt( W_temp, 10, 0)
		if err != nil { fmt.Println("Wrong configuration: Amount of workers is not a number"); os.Exit(1) }
	}
	
	flag.Parse()
	
	length := flag.NArg()
	if length == 0 {
		fmt.Println("Nothing to download, bye!")
		os.Exit(0)
	}

	//	creating directory for downloads if not yet done
    if err := os.MkdirAll(IMGDIR, 0777); err != nil {
         panic(err)
    }

	imgid := flag.Arg(length-1) //Last argument given presumed to be number of image on site. Temporally, because later would do with flags.

	//	fmt.Println(key) //Just checking that I am not wrong

	_, err = strconv.ParseInt(imgid, 10, 0)
	if err != nil { fmt.Println("Wrong input: can not parse", imgid, "as a number"); os.Exit(1) }
	
	
	fmt.Println("Processing image No " + imgid)
	
	imgdat := make (chan Image, 1)
	
	parseImg (imgdat, imgid, key)

	dlimage(imgdat)

}

type Image struct {
	url      string
	filename string
	hash     string
}

func parseImg(imgchan chan<- Image, imgid string, key string) {

	source := "http://derpiboo.ru/" + imgid + ".json?nofav=&nocomments="
	if (  key != "" ) { source = source + "&key=" + key }
	
//	fmt.Println(source)

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
	
	return 
}

func dlimage(imgchan <-chan Image) {
//	fmt.Println("reading channel")
	
	
	imgdata := <-imgchan
	
	fmt.Println("Saving as ", imgdata.filename)
	PathSep, _ := strconv.Unquote(strconv.QuoteRune(os.PathSeparator))

	output, err := os.Create(IMGDIR + PathSep + imgdata.filename)
	if err !=err { panic(err)}
	defer output.Close()

	response, err := http.Get(imgdata.url)
	if err != nil {
		panic(err)	
		fmt.Println("Error while downloading", imgdata.url, "-", err)
		return
	}
	defer response.Body.Close()
	
//	hash := sha512.New()

	io.Copy(output, response.Body)

/*	io.Copy(hash, response.Body)
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

}

func parseTag ( imgchan chan<- Image, tag string, key string) {

}

