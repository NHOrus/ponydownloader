package main

import (
	"encoding/json"
	"net/url"
	"path"
	"strconv"
	"sync"
	//	"github.com/davecgh/go-spew/spew"
)

var (
	derpiURL = url.URL{
		Scheme: "https",
		Host:   "derpibooru.org",
	}
	derpiquery url.Values
)

//RawImage contains data we got from API that needs to be modified before further usage
type RawImage struct {
	Imgid          string `json:"id"`
	URL            string `json:"image"`
	Score          int    `json:"score"`
	OriginalFormat string `json:"original_format"`
	Faves          int    `json:"faves"`
}

//Image contains data needed to filter fetch and save image
type Image struct {
	Imgid    string
	URL      *url.URL
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

	fn := dat.Imgid + "." + dat.OriginalFormat
	tu, _ := url.Parse(dat.URL)
	tu.Scheme = derpiURL.Scheme
	tu.Path = path.Dir(tu.Path) + "/" + fn
	return Image{

		Imgid:    dat.Imgid,
		Filename: fn,
		URL:      tu,
		Score:    dat.Score,
		Faves:    dat.Faves,
	}
}

//ParseImg gets image IDs, fetches information about those images from Derpibooru and pushes them into the channel.
func (imgchan ImageCh) ParseImg(ids []int, key string) {

	for _, imgid := range ids {

		if isInterrupted() {
			break
		}

		derpiURL.Path = strconv.Itoa(imgid) + ".json"
		derpiURL.RawQuery = derpiquery.Encode()

		lInfo("Getting image info at:", derpiURL)

		body, err := getJSON(derpiURL.String())
		if err != nil {
			lErr(err)
			break
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
	derpiURL.Path = "search.json"
	derpiquery.Add("sbq", opts.Tag)
	derpiURL.RawQuery = derpiquery.Encode()
	lInfo("Searching as", derpiURL.String())

	for page := opts.StartPage; opts.StopPage == 0 || page <= opts.StopPage; page++ {

		if isInterrupted() {
			break
		}

		lInfo("Searching page", page)
		derpiquery.Set("page", strconv.Itoa(page))
		derpiURL.RawQuery = derpiquery.Encode()

		body, err := getJSON(derpiURL.String())
		if err != nil {
			lErr("Error while getting json from page ", page)
			lErr(err)
			break
		}

		var dats Search
		err = json.Unmarshal(body, &dats)

		if err != nil {
			lErr("Error while parsing search page", page)
			lErr(err)
			if serr, ok := err.(*json.SyntaxError); ok { //In case crap was still given, we are looking at it.
				lErr("Occurred at offset: ", serr.Offset)
			}
			break

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
