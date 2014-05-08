package derpiapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Image struct {
	Imgid    int
	Url      string
	Filename string
	Score    int
	//	Hashval     string
}

func InfoToChannel(dat map[string]interface{}, imgchan chan<- Image) {

	var imgdata Image

	imgdata.Url = "http:" + dat["image"].(string)
	//	imgdata.Hashval = dat["sha512_hash"].(string)
	if (dat["original_format"].(string) == "svg") {
		imgdata.Url = "https://derpicdn.net/img/download/" + strings.Join(strings.Split("/",dat["image"].(string)[6:8]),  "/") + "/" + strconv.FormatFloat(dat["id_number"].(float64), 'f', -1, 64) + ".svg"
	}
	imgdata.Filename = (strconv.FormatFloat(dat["id_number"].(float64), 'f', -1, 64) + "." + dat["file_name"].(string) + "." + dat["original_format"].(string))
	imgdata.Imgid = int(dat["id_number"].(float64))
	imgdata.Score = int(dat["score"].(float64))

	//	for troubleshooting - possibly debug flag?
	//	fmt.Println(dat)
	fmt.Println(imgdata.Url)
	//	fmt.Println(imgdata.Hashval)
	//	fmt.Println(imgdata.Filename)

	imgchan <- imgdata
}

func ParseImg(imgchan chan<- Image, imgid string, KEY string, elog *log.Logger) {

	source := "http://derpiboo.ru/images/" + imgid + ".json?nofav=&nocomments="
	if KEY != "" {
		source = source + "&key=" + KEY
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

func DlImg(imgchan <-chan Image, done chan bool, elog *log.Logger, IMGDIR string) {

	fmt.Println("Worker started; reading channel") //nice notification that we are not forgotten

	for {

		imgdata, more := <-imgchan

		if !more { //checking that there is an image in channel
			done <- true //well, there is no images in channel, it means we got them all, so synchronization is kicking in and ending the process
			break        //Just in case, so it would not stupidly die when program finishes - it will die smartly
		}

		if imgdata.Filename == "" {
			elog.Println("Empty filename. Oops?") //something somewhere had gone wrong, not a cause to die, going to the next image
		} else {

			log.Println("Saving as", imgdata.Filename)

			func() { // To not hold all the files open when there is no need. All pointers to files are in the scope of this function.

				output, err := os.Create(IMGDIR + string(os.PathSeparator) + imgdata.Filename) //And now, THE FILE!
				if err != err {
					elog.Println("Error when creating file for image" + strconv.Itoa(imgdata.Imgid))
					elog.Println(err) //Either we got no permisson or no space, end of line
					return
				}
				defer output.Close() //Not forgetting to deal with it after completing download

				response, err := http.Get(imgdata.Url)
				if err != nil {
					elog.Println("Error when getting image", imgdata.Imgid)
					elog.Println(err)
					return
				}
				defer response.Body.Close() //Same, we shall not listen to the void when we finished getting image

				io.Copy(output, response.Body) //	Writing things we got from Derpibooru into the file and into hasher

			}()
		}

		//fmt.Println("\n", hex.EncodeToString(hash.Sum(nil)), "\n", imgdata.hash )

	}
}

func ParseTag(imgchan chan<- Image, tag string, KEY string, STARTPAGE int, STOPPAGE int, elog *log.Logger) {

	source := "http://derpiboo.ru/search.json?nofav=&nocomments=" //yay hardwiring url strings!

	if KEY != "" {
		source = source + "&key=" + KEY
	}

	fmt.Println("Searching as", source+"&q="+tag)
	var working = true
	i := STARTPAGE
	for working {
		func() { //I suspect that all those returns could be dealt with in some way. But lazy.
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

			err = json.Unmarshal(body, &dats) //transforming json into native slice of maps

			if err != nil {
				elog.Println("Error while parsing search page", i)
				elog.Println(err)
				return

			}

			if len(dats) == 0 {
				fmt.Println("Pages are over") //Does not mean that process is over.
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
