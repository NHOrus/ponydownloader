package derpiapi

import (
	"encoding/json"
	//	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//Image contains data we got from API that we are using to filter and fetch said image next
type Image struct {
	Imgid          int    `json:"id_number"`
	URL            string `json:"image"`
	Filename       string
	Score          int    `json:"score"`
	OriginalFormat string `json:"original_format"`
	// Hashval     string `json:"sha512_hash"`
}

//Search returns to us array of searched images...
type Search struct {
	Images []Image `json:"search"`
}

type ImageCh chan Image

//infotochannel gets unmarchalled JSON info and plugs it into channel so it would be processed in other places
func infotochannel(dat Image, imgchan ImageCh) {
	dat.Filename = strconv.Itoa(dat.Imgid) + "." + dat.OriginalFormat
	dat.URL = "https:" + dat.URL
	if dat.OriginalFormat == "svg" {
		i := strings.LastIndex(dat.URL, ".")
		if i != -1 {
			dat.URL = dat.URL[:i] + ".svg" //Was afraid to extract things I needed from the date field, so extracting them from URL.
		}
	}
	imgchan <- dat
}

//ParseImg gets image ID and, fetches information about this image from Derpibooru and puts it into the channel.
func (imgchan ImageCh) ParseImg(imgid int, KEY string, elog *log.Logger) {

	source := "https://derpiboo.ru/images/" + strconv.Itoa(imgid) + ".json"
	if KEY != "" {
		source = source + "?key=" + KEY
	}

	log.Println("Getting image info at:", source)

	resp, err := http.Get(source) //Getting our nice http response. Needs checking for 404 and other responses that are... less expected
	if err != nil {
		elog.Println(err)
		return
	}

	defer resp.Body.Close() //and not forgetting to close it when it's done

	var dat Image

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

	infotochannel(dat, imgchan)

	close(imgchan) //closing channel, we are done here

	return
}

//DlImg downloads image on disk, given image data
func (imgchan ImageCh) DlImg(done chan bool, elog *log.Logger, IMGDIR string) {

	log.Println("Worker started; reading channel") //nice notification that we are not forgotten

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

				start := time.Now()

				response, err := http.Get(imgdata.URL)
				if err != nil {
					elog.Println("Error when getting image", imgdata.Imgid)
					elog.Println(err)
					return
				}
				defer response.Body.Close() //Same, we shall not listen to the void when we finished getting image

				size, err := io.Copy(output, response.Body) //	Writing things we got from Derpibooru into the file and into hasher
				if err != nil {
					elog.Println("Unable to write image on disk, id ", imgdata.Imgid)
					elog.Println(err)
					return
				}
				timed := time.Since(start).Seconds()
				log.Printf("Downloaded %d bytes in %.2fs, speed %s/s\n", size, timed, fmtbytes(float64(size)/timed))
			}()
		}

		//fmt.Println("\n", hex.EncodeToString(hash.Sum(nil)), "\n", imgdata.hash )

	}
}

func (imgchan ImageCh) ParseTag(tag string, KEY string, STARTPAGE int, STOPPAGE int, elog *log.Logger) {

	source := "https://derpiboo.ru/search.json?" //yay hardwiring url strings!

	if KEY != "" {
		source = source + "key=" + KEY +"&"
	}

	log.Println("Searching as", source + "q=" + tag)
	var working = true
	i := STARTPAGE
	for working {
		func() { //I suspect that all those returns could be dealt with in some way. But lazy.
			log.Println("Searching page", i)
			resp, err := http.Get(source + "q=" + tag + "&page=" + strconv.Itoa(i)) //Getting our nice http response. Needs checking for 404 and other responses that are... less expected
			defer resp.Body.Close()                                                  //and not forgetting to close it when it's done. And before we panic and die horribly.
			if err != nil {
				elog.Println("Error while getting search page", i)
				elog.Println(err)
				return
			}

			var dats Search //Because we got array incoming instead of single object, we using an slive of maps!

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

			if len(dats.Images) == 0 {
				log.Println("Pages are over") //Does not mean that process is over.
				working = false
				return
			} //exit due to finishing all pages

			for _, dat := range dats.Images {
				infotochannel(dat, imgchan)
			}
			if STOPPAGE != 0 && i > STOPPAGE { //stop page is included, but if not set? Work to the end of tag
				working = false
				return
			}
			i++

		}()
	}

	close(imgchan)
}
