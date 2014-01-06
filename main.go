package main

import (
"fmt"
//	"net"
	"os"
//	"errors"
//	"log"
	"io"
	"github.com/vaughan0/go-ini"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
	)	


func main(){

	fmt.Println("Derpiboo.ru Downloader version 0.0.2 \nWorking")
//	fmt.Println("Working")

	config, err := ini.LoadFile("config.ini") // Loading default config file and checking for various errors.

	if os.IsNotExist(err) { 
		panic("config.ini does not exist, create it")
		}

	if err != nil { panic(err) }

	key, ok := config.Get("main", "key") //Need to make things for key == nil

	if !ok {
		panic("'key' variable missing from 'main' section")
		}

	length := len(os.Args)
	if length == 1 {
		fmt.Println("Nothing to download, bye!")
		os.Exit(0)
	}
	
	imgid := os.Args[length - 1] //Last argument given presumed to be number of image on site. Temporally, because later would do with flags.
	
//	fmt.Println(key) //Just checking that I am not wrong
	
	fmt.Println("Processing image No " + imgid)
	
	imgdat := parseImg (imgid, key)

	dlimage(imgdat)

}

type Image struct {
	url		string
	filename	string
	hash		string
	}

func parseImg(imgid string, key string) (imgdata Image) {

	source := "http://derpiboo.ru/" + imgid + ".json?nofav=&nocomments=?key=" + key //correct way is to assemble all the different arguments and only then append them to source url. I can live with hardcoded source site. May be add check for derpiboo.ru or derpibooru.org ?
//	fmt.Println(source)
	
	resp, err := http.Get(source)	//Getting our nice http response. Needs checking for 404 and other responses that are... less expected
		if err != nil {
			panic(err)
		}
	
	defer resp.Body.Close()	//and not forgetting to close it when it's done
	
	var dat map[string]interface{}
	
	body, err := ioutil.ReadAll(resp.Body)	//stolen from official documentation
		if err != nil {
			panic(err)
		}
	
	if err := json.Unmarshal(body, &dat); //transforming json into native map
	
	err != nil {
	        panic(err)
		}

	imgdata.url = "http:" + dat["image"].(string)
	imgdata.hash = dat["sha512_hash"].(string)  //for the future and checking that we got file right
	imgdata.filename = strconv.FormatFloat(dat["id_number"].(float64), 'f', -1, 64) + "." + dat["file_name"].(string)
	
//	fmt.Println(strconv.FormatFloat(dat["id_number"].(float64), 'f', -1, 64))

//	fmt.Println(dat)

// for now and troubleshooting
//	fmt.Println(imgdata.url)
//	fmt.Println(imgdata.hash)
//	fmt.Println(imgdata.filename)

	return
	}
	

func dlimage(imgdata Image) {
	fmt.Println("Saving as ", imgdata.filename)
	output, err := os.Create(imgdata.filename)
	defer output.Close()
	
	response, err := http.Get(imgdata.url)
    if err != nil {
		fmt.Println("Error while downloading", imgdata.url, "-", err)
		return
    }
    defer response.Body.Close()
    
    io.Copy(output, response.Body)
	
	}

func parseTag ( tag string ) (imgdata chan-< Image) {}
