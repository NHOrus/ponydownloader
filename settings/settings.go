package settings

import (
	"fmt"
	"io"
	"log"
	"os"
	//"github.com/vaughan0/go-ini"
)

//Setting up logfile as I want it to: Copy to event.log, copy to commandline
func SetLog() (retlog *log.Logger, logfile *os.File) {

	logfile, err := os.OpenFile("event.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644) //file for putting errors into

	if err != nil {
		panic(err)
	}

	retlog = log.New(io.MultiWriter(logfile, os.Stderr), "Errors at ", log.LstdFlags) //Setting stuff for our logging: both errors and events.

	log.SetPrefix("Happens at ")
	log.SetFlags(log.LstdFlags)
	log.SetOutput(io.MultiWriter(logfile, os.Stdout)) //we write in file and stdout
	log.Println("Program start")

	return
}

//Writing default configuration data into file.
func WriteConfigDefault() {
	config, err := os.OpenFile("config.ini", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer config.Close()
	
	if err != nil {
		panic(err)
	}
	
	_, err = fmt.Fprintln(config,
	"[main]",
	"key =",
	"queue_depth =",
	"downdir =")

	if err != nil {
		panic(err)
	}
}
