package main

import (
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"hash"
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
	Hashval        string `json:"sha512_hash"`
}

//Search returns to us array of searched images...
type Search struct {
	Images []Image `json:"search"`
}

//ImageCh is a channel of image data. You can put images into channel by parsing
//Derpibooru API by id(s) or  by tags and you can download images that are already
//in channel
type ImageCh chan Image

//Push gets unmarchalled JSON info, massages it and plugs it into channel so it
//would be processed in other places
func (imgchan ImageCh) push(dat Image) {
	dat.Filename = strconv.Itoa(dat.Imgid) + "." + dat.OriginalFormat
	dat.URL = prefix + dat.URL
	if dat.OriginalFormat == "svg" {
		i := strings.LastIndex(dat.URL, ".")
		if i != -1 {
			dat.URL = dat.URL[:i] + ".svg" //Was afraid to extract things I needed from the date field, so extracting them from URL.
		}
	}
	imgchan <- dat
}

//ParseImg gets image IDs, fetches information about those images from Derpibooru and pushes them into the channel.
func (imgchan ImageCh) ParseImg() {

	for _, imgid := range opts.Args.IDs {
		source := prefix + "//derpibooru.org/images/" + strconv.Itoa(imgid) + ".json"
		if opts.Key != "" {
			source = source + "?key=" + opts.Key
		}

		log.Println("Getting image info at:", source)

		response, err := http.Get(source) //Getting our nice http response. Needs checking for 404 and other responses that are... less expected
		if err != nil {
			if _, ok := err.(x509.UnknownAuthorityError); ok && opts.Unsafe {
				log.Println("Warning: ", err)
			} else {
				elog.Println(err)
				continue
			}
		}

		defer func() {
			err = response.Body.Close() //and not forgetting to close it when it's done
			if err != nil {
				elog.Fatalln("Could  not close server response")
			}
		}()
		var dat Image

		if !okHTTPStatus(response) {
			continue
		}

		body, err := ioutil.ReadAll(response.Body) //stolen from official documentation
		if err != nil {
			elog.Println(err)
			continue
		}

		if err := json.Unmarshal(body, &dat); //transforming json into native map

		err != nil {
			elog.Println(err)
			continue
		}

		imgchan.push(dat)
	}

	close(imgchan) //closing channel, we are done here

	return
}

//DlImg reads image data from channel and downloads specified images to disc
func (imgchan ImageCh) DlImg() {

	log.Println("Worker started; reading channel") //nice notification that we are not forgotten

	hasher := sha512.New() //checksums will be done in this

	for imgdata := range imgchan {

		if imgdata.Filename == "" {
			elog.Println("Empty filename. Oops?") //something somewhere had gone wrong, not a cause to die, going to the next image
			continue
		}

		log.Println("Saving as", imgdata.Filename)

		imgdata.saveImage(hasher)

	}
	done <- true
}

func (imgdata Image) saveImage(hasher hash.Hash) { // To not hold all the files open when there is no need. All pointers to files are in the scope of this function.

	output, err := os.Create(opts.ImageDir + string(os.PathSeparator) + imgdata.Filename) //And now, THE FILE!
	if err != err {
		elog.Println("Error when creating file for image" + strconv.Itoa(imgdata.Imgid))
		elog.Println(err) //Either we got no permisson or no space, end of line
		return
	}
	defer func() {
		err = output.Close() //Not forgetting to deal with it after completing download
		if err != nil {
			elog.Fatalln("Could  not close downloaded file")
		}
	}()

	start := time.Now() //timing download time. We can't begin it sooner, not sure if we can begin it later

	response, err := http.Get(imgdata.URL)

	if err != nil {
		if _, ok := err.(x509.UnknownAuthorityError); ok && opts.Unsafe { //With flag to sidestep outdated root certificates
			log.Println("Warning: ", err)
		} else {
			elog.Println("Error when getting image", imgdata.Imgid)
			elog.Println(err)
			return
		}
	}
	defer func() {
		err = response.Body.Close() //Same, we shall not listen to the void when we finished getting image
		if err != nil {
			elog.Fatalln("Could  not close server response")
		}
	}()
	if !okHTTPStatus(response) {
		return
	}
	size, err := io.Copy(output, io.TeeReader(response.Body, hasher)) //	Writing things we got from Derpibooru into the file and into hasher
	if err != nil {
		elog.Println("Unable to write image on disk, id ", imgdata.Imgid)
		elog.Println(err)
		return
	}
	timed := time.Since(start).Seconds()

	hash := hex.EncodeToString(hasher.Sum(nil)) //Ah, hexadecimals. It may be easier to convert hex to []byte, but whatever

	if hash != imgdata.Hashval {
		elog.Println("Hash mismatch, got ", hash, " instead of ", imgdata.Hashval)
	}

	hasher.Reset() //cleanup, so we don't allocate each loop

	log.Printf("Downloaded %d bytes in %.2fs, speed %s/s\n", size, timed, fmtbytes(float64(size)/timed))
}

//ParseTag gets image tags, fetches information about all images it could from Derpibooru and pushes them into the channel.
func (imgchan ImageCh) ParseTag() {

	//Unlike main, I don't see how I could separate bits out to decrease complexity
	source := prefix + "//derpibooru.org/search.json?q=" + opts.Tag //yay hardwiring url strings!

	if opts.Key != "" {
		source = source + "&key=" + opts.Key
	}

	log.Println("Searching as", source)

	for i := opts.StartPage; opts.StopPage == 0 || i <= opts.StopPage; i++ {
		log.Println("Searching page", i)

		response, err := http.Get(source + "&page=" + strconv.Itoa(i))
		//Getting our nice http response.

		//This error check may be given it's own function, later. Not sure of best way to do it.
		if err != nil {
			if _, ok := err.(x509.UnknownAuthorityError); ok && opts.Unsafe { //With flag to sidestep outdated root certificates
				log.Println("Warning: ", err)
			} else {
				elog.Println("Error while getting search page", i)
				elog.Println(err)
				continue
			}
		}

		defer func() {
			err = response.Body.Close() //and not forgetting to close it when it's done. And before we panic and die horribly.
			if err != nil {
				elog.Fatalln("Could  not close server response")
			}
		}()

		if !okHTTPStatus(response) { //Checking that we weren't given crap instead of candy
			continue
		}

		body, err := ioutil.ReadAll(response.Body) //stolen from official documentation
		if err != nil {
			elog.Println("Error while reading search page", i)
			elog.Println(err)
			continue
		}

		var dats Search                   //Because we got array incoming instead of single object, we using a slive of maps!
		err = json.Unmarshal(body, &dats) //transforming json into native view

		if err != nil {
			elog.Println("Error while parsing search page", i)
			elog.Println(err)
			if serr, ok := err.(*json.SyntaxError); ok { //In case crap was still given, we are looking at it.
				log.Println("Occurred at offset:", serr.Offset)
			}
			continue

		}

		if len(dats.Images) == 0 {
			log.Println("Pages are over") //Does not mean that process is over.
			break
		} //exit due to finishing all pages

		for _, dat := range dats.Images {
			imgchan.push(dat)
		}

	}

	close(imgchan)
}

func okHTTPStatus(chk *http.Response) bool {
	switch chk.StatusCode {
	case http.StatusOK, http.StatusNotModified:
		return true
	case http.StatusGatewayTimeout,
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusHTTPVersionNotSupported:
		elog.Println("Server error: ", chk.Status)
		return false
	case http.StatusBadRequest,
		http.StatusTeapot,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusRequestURITooLong,
		http.StatusExpectationFailed:
		elog.Println("Incorrect request to server, error ", chk.Status)
		elog.Println("Possible API changes")
		return false
	default:
		elog.Println("Got something weird from server: ", chk.Status)
		elog.Println("Continuing anyway")
		return true
	}
}
