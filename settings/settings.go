package settings

import (
	"io"
	"log"
	"os"
	//"github.com/vaughan0/go-ini"
)

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

func WriteConfig() {

}
