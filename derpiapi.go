package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
//	"github.com/davecgh/go-spew/spew"
)

//Image contains data we got from API that we are using to filter and fetch said image next
type Image struct {
	Imgid          int    `json:"id_number"`
	URL            string `json:"image"`
	Filename       string
	Score          int    `json:"score"`
	OriginalFormat string `json:"original_format"`
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
	dat.URL = prefix + "/" + path.Dir(dat.URL) + "/" + dat.Filename
	if dat.OriginalFormat == "svg" {
		i := strings.LastIndex(dat.URL, ".")
		if i != -1 {
			dat.URL = dat.URL[:i] + ".svg" //Was afraid to extract things I needed from the date field, so extracting them from URL.
		}
	}
	imgchan <- dat
}

//ParseImg gets image IDs, fetches information about those images from Derpibooru and pushes them into the channel.
func (imgchan ImageCh) ParseImg(ids []int, key string) {

	for _, imgid := range ids {
		source := prefix + "//derpibooru.org/images/" + strconv.Itoa(imgid) + ".json"
		if key != "" {
			source = source + "?key=" + key
		}

		lInfo("Getting image info at:", source)

		body, err := getRemoteJSON(source)
		if err != nil {
			lErr(err)
			continue
		}
		var dat Image
		if err := json.Unmarshal(body, &dat); //transforming json into native map

		err != nil {
			lErr(err)
			continue
		}

		imgchan.push(dat)
	}

	close(imgchan) //closing channel, we are done here

	return
}

//DlImg reads image data from channel and downloads specified images to disc
func (imgchan ImageCh) downloadImages(opts *Settings) {

	lInfo("Worker started; reading channel") //nice notification that we are not forgotten

	for imgdata := range imgchan {

		if imgdata.Filename == "" {
			lErr("Empty filename. Oops?") //something somewhere had gone wrong, not a cause to die, going to the next image
			continue
		}

		lInfo("Saving as", imgdata.Filename)

		imgdata.saveImage(opts)

	}
}

func (imgdata Image) saveImage(opts *Settings) { // To not hold all the files open when there is no need. All pointers to files are in the scope of this function.

	output, err := os.Create(opts.ImageDir + string(os.PathSeparator) + imgdata.Filename) //And now, THE FILE!
	if err != err {
		lErr("Error when creating file for image: ", strconv.Itoa(imgdata.Imgid))
		lErr(err) //Either we got no permisson or no space, end of line
		return
	}
	defer func() {
		err = output.Close() //Not forgetting to deal with it after completing download
		if err != nil {
			lFatal("Could  not close downloaded file")
		}
	}()

	start := time.Now() //timing download time. We can't begin it sooner, not sure if we can begin it later

	response, err := http.Get(imgdata.URL)

	if err != nil {
		lErr("Error when getting image: ", strconv.Itoa(imgdata.Imgid))
		lErr(err)
		return
	}
	defer func() {
		err = response.Body.Close() //Same, we shall not listen to the void when we finished getting image
		if err != nil {
			lFatal("Could  not close server response")
		}
	}()
	if !okHTTPStatus(response) {
		return
	}
	
	size, err := io.Copy(output, response.Body) //
	if err != nil {
		lErr("Unable to write image on disk, id: ", strconv.Itoa(imgdata.Imgid))
		lErr(err)
		return
	}
	timed := time.Since(start).Seconds()

	lInfof("Downloaded %d bytes in %.2fs, speed %s/s\n", size, timed, fmtbytes(float64(size)/timed))

	sizestring, prs := response.Header["Content-Length"]
	if !prs {
		lErr("Filesize not provided")
		return
	}

	expsize, err := strconv.ParseInt(sizestring[0], 10, 64)
	if err != nil {
		lErr("Unable to get expected filesize")
		return
	}
	if expsize != size {
		lErr("Unable to download full image")
	}
}

//ParseTag gets image tags, fetches information about all images it could from Derpibooru and pushes them into the channel.
func (imgchan ImageCh) ParseTag(opts *TagOpts, key string) {

	//Unlike main, I don't see how I could separate bits out to decrease complexity
	source := prefix + "//derpibooru.org/search.json?q=" + opts.Tag //yay hardwiring url strings!

	if key != "" {
		source = source + "&key=" + key
	}

	lInfo("Searching as", source)

	for i := opts.StartPage; opts.StopPage == 0 || i <= opts.StopPage; i++ {
		lInfo("Searching page", i)

		body, err := getRemoteJSON(source + "&page=" + strconv.Itoa(i))
		if err != nil {
			lErr("Error while json from page ", i)
			lErr(err)
			continue
		}

		var dats Search                   //Because we got array incoming instead of single object, we using a slive of maps!
		err = json.Unmarshal(body, &dats) //transforming json into native view

		if err != nil {
			lErr("Error while parsing search page", i)
			lErr(err)
			if serr, ok := err.(*json.SyntaxError); ok { //In case crap was still given, we are looking at it.
				lErr("Occurred at offset: ", serr.Offset)
			}
			continue

		}

		if len(dats.Images) == 0 {
			lInfo("Pages are all over") //Does not mean that process is over.
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
		lErr("Server error: ", chk.Status)
		return false
	case http.StatusBadRequest,
		http.StatusTeapot,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusRequestURITooLong,
		http.StatusExpectationFailed:
		lErr("Incorrect request to server, error ", chk.Status)
		lErr("Possible API changes")
		return false
	default:
		lWarn("Got something weird from server: ", chk.Status)
		lWarn("Continuing anyway")
		return true
	}
}

func getRemoteJSON(source string) (body []byte, err error) {
	response, err := http.Get(source)
	//Getting our nice http response.

	//This error check may be given it's own function, later. Not sure of best way to do it.
	if err != nil {
		return nil, err

	}

	defer func() {
		err = response.Body.Close() //and not forgetting to close it when it's done. And before we panic and die horribly.
		if err != nil {
			lFatal("Could  not close server response")
		}
	}()

	if !okHTTPStatus(response) { //Checking that we weren't given crap instead of candy
		return nil, errors.New("Incorrect server response")
	}

	body, err = ioutil.ReadAll(response.Body) //stolen from official documentation
	if err != nil {
		return nil, err
	}

	return body, nil
}
