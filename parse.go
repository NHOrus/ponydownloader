package main

import (
	"encoding/json"
	"path"
	"strconv"
	"sync"
	//	"github.com/davecgh/go-spew/spew"
)

//RawImage contains data we got from API that needs to be modified before further usage
type RawImage struct {
	Imgid          int    `json:"id_number"`
	URL            string `json:"image"`
	Score          int    `json:"score"`
	OriginalFormat string `json:"original_format"`
	Faves          int    `json:"faves"`
}

//Image contains data needed to filter fetch and save image
type Image struct {
	Imgid    int
	URL      string
	Filename string
	Score    int
	Faves    int
}

//Search returns to us array of searched images...
type Search struct {
	Images []RawImage `json:"search"`
}

//ImageCh is a channel of image data. You can put images into channel by parsing
//Derpibooru API by id(s) or  by tags and you can download images that are already
//in channel
type ImageCh chan Image

//Push gets unmarchalled JSON info, massages it and plugs it into channel so it
//would be processed in other places
func trim(dat RawImage) Image {

	tfn := strconv.Itoa(dat.Imgid) + "." + dat.OriginalFormat
	return Image{
		Imgid:    dat.Imgid,
		Filename: tfn,
		URL:      prefix + "/" + path.Dir(dat.URL) + "/" + tfn,
		Score:    dat.Score,
		Faves:    dat.Faves,
	}
}

//ParseImg gets image IDs, fetches information about those images from Derpibooru and pushes them into the channel.
func (imgchan ImageCh) ParseImg(ids []int, key string) {

	for _, imgid := range ids {

		if isParseInterrupted() {
			break
		}

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
		var dat RawImage
		if err := json.Unmarshal(body, &dat); //transforming json into native map

		err != nil {
			lErr(err)
			continue
		}

		imgchan <- trim(dat)
	}

	close(imgchan) //closing channel, we are done here

	return
}

//DlImg reads image data from channel and downloads specified images to disc
func (imgchan ImageCh) downloadImages(opts *Config) {

	lInfo("Worker started; reading channel") //nice notification that we are not forgotten
	var n int
	var size int64
	var l sync.Mutex
	var wg sync.WaitGroup
	for k := 0; k < 4; k++ {
		wg.Add(1)
		go func() {
			for imgdata := range imgchan {

				lInfo("Saving as", imgdata.Filename)

				tsize, ok := imgdata.saveImage(opts)
				l.Lock()
				size += tsize
				if ok {
					n++
				}
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	lInfof("Downloaded %d images, for a total of %s", n, fmtbytes(float64(size)))
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

		if isParseInterrupted() {
			break
		}

		lInfo("Searching page", i)

		body, err := getRemoteJSON(source + "&page=" + strconv.Itoa(i))
		if err != nil {
			lErr("Error while getting json from page ", i)
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
			imgchan <- trim(dat)
		}

	}

	close(imgchan)
}

func isParseInterrupted() bool {
	select {
	case <-interruptParse:
		return true
	default:
		return false
	}
}
