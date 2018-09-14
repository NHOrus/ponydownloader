package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"os"
	"strconv"
	"time"
)

func getJSON(source string) (body []byte, err error) {
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
		return nil, fmt.Errorf("Incorrect server response")
	}

	body, err = ioutil.ReadAll(response.Body) //stolen from official documentation
	if err != nil {
		return nil, err
	}

	return body, nil
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
		lErr("Incorrect request to server, error:", chk.Status)
		lErr("Possible API changes")
		return false
	default:
		lWarn("Got something weird from server:", chk.Status)
		lWarn("Continuing anyway")
		return true
	}
}

func (imgdata Image) saveImage(opts *Config) (size int64, ok bool) { // To not hold all the files open when there is no need. All file descriptors are in the scope of this function.

	filepath := constructFilepath(imgdata.Filename, opts.ImageDir)

	fsize := getFileSize(filepath)

	start := time.Now() //Timing download time. We can't begin it sooner, not sure if we can begin it later

	response, err := http.Get(imgdata.URL.String())

	if err != nil {
		lErr("Error when getting image: ", imgdata.Imgid)
		lErr(err)
		return
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			lFatal("Could not close server response")
		}
	}()

	if !okHTTPStatus(response) {
		return
	}

	expsize := getRemoteSize(response.Header)

	if expsize == fsize {
		lInfo("Skipping: no-clobber")
		return
	}

	if err != nil {
		lErr("Error when getting image: ", imgdata.Imgid)
		lErr(err)
		return
	}

	output, err := os.Create(filepath) //And now, THE FILE! New, truncated, ready to write
	if err != nil {
		lErr("Error when creating file for image: ", imgdata.Imgid)
		lErr(err) //Either we got no permission or no space, end of line
		return
	}
	defer func() {
		err = output.Close() //Not forgetting to deal with it after completing download
		if err != nil {
			lFatal("Could  not close downloaded file")
		}
	}()

	size, err = io.Copy(output, response.Body) //Preventing creation of temporary buffer in memory
	if err != nil {
		lErr("Unable to write image on disk, id: ", imgdata.Imgid)
		lErr(err)
		return
	}
	timed := time.Since(start).Seconds()

	lInfof("Downloaded %d bytes in %.2fs, speed %s/s\n", size, timed, fmtbytes(float64(size)/timed))
	ok = true

	if expsize != size {
		lErr("Unable to download full image")
	}
	return
}

func getFileSize(path string) int64 {
	fstat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return fstat.Size()

}

func getRemoteSize(head http.Header) (expsize int64) {

	sizestring, ok := head["Content-Length"]
	if !ok {
		lErr("Filesize not provided")
	}

	expsize, err := strconv.ParseInt(sizestring[0], 10, 64)
	if err != nil {
		lErr("Unable to get expected filesize")
	}
	return
}

func constructFilepath(filename string, imagedir string) string {
	if imagedir == "" {
		return filename
	}
	return imagedir + string(os.PathSeparator) + filename
}
