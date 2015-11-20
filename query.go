package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
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

func getImage(source string) ([]byte, http.Header, error) {
	response, err := http.Get(source)
	if err != nil {
		return nil, nil, err

	}
	defer func() {
		err = response.Body.Close() //Same, we shall not listen to the void when we finished getting image
		if err != nil {
			lFatal("Could  not close server response")
		}
	}()
	if !okHTTPStatus(response) {
		return nil, nil, fmt.Errorf("Incorrect server response")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, response.Header, err
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
		n, err := fmt.Scanln(&response)
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
