package main

import (
	"encoding/json"
	"os"
	"path"
	"strconv"
	"time"
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

func (imgdata *Image) isDeepEqual(b Image) bool {
	if b.Imgid == imgdata.Imgid &&
		b.URL == imgdata.URL &&
		b.Filename == imgdata.Filename &&
		b.Score == imgdata.Score &&
		b.Faves == imgdata.Faves {
		return true
	}
	return false
}

//ImageCh is a channel of image data. You can put images into channel by parsing
//Derpibooru API by id(s) or  by tags and you can download images that are already
//in channel
type ImageCh chan Image

//Push gets unmarchalled JSON info, massages it and plugs it into channel so it
//would be processed in other places
func (imgchan ImageCh) push(dat RawImage) {
	
	tfn := strconv.Itoa(dat.Imgid) + "." + dat.OriginalFormat
	imgchan <- Image{
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

		imgchan.push(dat)
	}

	close(imgchan) //closing channel, we are done here

	return
}

//DlImg reads image data from channel and downloads specified images to disc
func (imgchan ImageCh) downloadImages(opts *Config) {

	lInfo("Worker started; reading channel") //nice notification that we are not forgotten
	var n, size int
	for imgdata := range imgchan {

		if imgdata.Filename == "" {
			lErr("Empty filename. Oops?") //something somewhere had gone wrong, not a cause to die, going to the next image
			continue
		}

		lInfo("Saving as", imgdata.Filename)

		size += imgdata.saveImage(opts)
		n++
	}
	lInfof("Downloaded %d images, for a total of %s", n, fmtbytes(float64(size)))
}

func (imgdata Image) saveImage(opts *Config) (size int) { // To not hold all the files open when there is no need. All pointers to files are in the scope of this function.

	output, err := os.Create(opts.ImageDir + string(os.PathSeparator) + imgdata.Filename) //And now, THE FILE!
	if err != nil {
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

	body, header, err := getImage(imgdata.URL)

	if err != nil {
		lErr("Error when getting image: ", strconv.Itoa(imgdata.Imgid))
		lErr(err)
		return
	}

	size, err = output.Write(body) //
	if err != nil {
		lErr("Unable to write image on disk, id: ", strconv.Itoa(imgdata.Imgid))
		lErr(err)
		return
	}
	timed := time.Since(start).Seconds()

	lInfof("Downloaded %d bytes in %.2fs, speed %s/s\n", size, timed, fmtbytes(float64(size)/timed))

	sizestring, prs := header["Content-Length"]
	if !prs {
		lErr("Filesize not provided")
		return
	}

	expsize, err := strconv.Atoi(sizestring[0])
	if err != nil {
		lErr("Unable to get expected filesize")
		return
	}
	if expsize != size {
		lErr("Unable to download full image")
		return 0
	}
	return
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
