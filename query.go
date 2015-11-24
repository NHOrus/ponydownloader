package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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
		lErr("Incorrect request to server, error ", chk.Status)
		lErr("Possible API changes")
		return false
	default:
		lWarn("Got something weird from server: ", chk.Status)
		lWarn("Continuing anyway")
		return true
	}
}

func makeHTTPSUnsafe() {
	lWarn("Disabling HTTPS trust check is unsafe and may lead to your data being monitored, stolen or falsified by third party")
	lWarn("Are you really want to continue? [Yes/No]")
	if !checkUserResponse() {
		lInfo("Continuing in safe mode")
		return
	}
	lWarn("Continuing in unsafe mode")
	http.DefaultClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
}

func checkUserResponse() bool {
	var response string
	for i := 0; i < 10; i++ {
		time.Sleep(3 * time.Second)
		n, _ := fmt.Scanln(&response)
		r := strings.ToLower(response)
		if n != 1 {
			continue
		}
		if r == "y" || r == "yes" {
			return true
		}
		if r == "n" || r == "no" {
			return false
		}

	}
	return false

}
func (imgdata Image) saveImage(opts *Config) (size int64) { // To not hold all the files open when there is no need. All pointers to files are in the scope of this function.

	filepath := opts.ImageDir + string(os.PathSeparator) + imgdata.Filename

	fsize := getFileSize(filepath)

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

	sizestring, prs := response.Header["Content-Length"]
	if !prs {
		lErr("Filesize not provided")
	}
	expsize, err := strconv.ParseInt(sizestring[0], 10, 64)
	if err != nil {
		lErr("Unable to get expected filesize")
	}

	if expsize == fsize {
		lInfo("Skipping: no-clobber")
		return 0
	}

	if err != nil {
		lErr("Error when getting image: ", strconv.Itoa(imgdata.Imgid))
		lErr(err)
		return
	}

	output, err := os.Create(filepath) //And now, THE FILE! New, truncated, ready to write
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

	size, err = io.Copy(output, response.Body) //Preventing creation of temporary buffer in memory
	if err != nil {
		lErr("Unable to write image on disk, id: ", strconv.Itoa(imgdata.Imgid))
		lErr(err)
		return
	}
	timed := time.Since(start).Seconds()

	lInfof("Downloaded %d bytes in %.2fs, speed %s/s\n", size, timed, fmtbytes(float64(size)/timed))

	if expsize != size {
		lErr("Unable to download full image")
		return 0
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
