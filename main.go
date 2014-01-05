package main

import ("fmt"
//	"net"
	"os"
//	"errors"
	"log"
//	"io"
	"github.com/vaughan0/go-ini"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
	)	


func main(){
	fmt.Println("Check one")
	config, err := ini.LoadFile("config.ini")
	if os.IsNotExist(err) { panic("config.ini does not exist, create it")}
	if err != nil { log.Fatal(err) }

	key, ok := config.Get("main", "key")
	if !ok {
		panic("'key' variable missing from 'main' section")
		}
	length := len(os.Args)
	imgid := os.Args[length - 1] //temporally, because later would do with flags.
	
//	fmt.Println(key)
	
	fmt.Println("Processing image No " + imgid)
	
	imgdat := parse (imgid, key)

	go dlimage(imgdat)

}

type Image struct {
	url		string
	filename	string
	hash		string
	}

func parse(imgid string, key string) (imgdata Image) {

	source := "http://derpiboo.ru/" + imgid + ".json?nofav=&nocomments=?key=" + key
//	fmt.Println(source)
	
	resp, err := http.Get(source)
		if err != nil {
			panic(err)
		}
	
	defer resp.Body.Close()
	
	var dat map[string]interface{}
	
	body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
	
	if err := json.Unmarshal(body, &dat); err != nil {
        panic(err)
    }

	imgdata.url = "http:" + dat["image"].(string)
	imgdata.hash = dat["sha512_hash"].(string)
	imgdata.filename = strconv.FormatFloat(dat["id_number"].(float64), 'f', -1, 64) + "." + dat["file_name"].(string)

//	fmt.Println(dat)

	fmt.Println(imgdata.url)
	fmt.Println(imgdata.hash)
	fmt.Println(imgdata.filename)

	return
	}
	

func dlimage(imgdata Image) {}
