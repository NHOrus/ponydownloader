package settings

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/vaughan0/go-ini" //We need some simple way to parse ini files, here it is, externally.
)

type Settings struct {
	QDepth int
	ImgDir string
	Key    string
}

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
func (WSet Settings) WriteConfig(elog *log.Logger) {
	config, err := os.OpenFile("config.ini", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer config.Close()

	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintln(config,
		"[main]", "",
		"key = "+WSet.Key,
		"queue_depth = "+string(WSet.QDepth),
		"downdir = "+WSet.ImgDir)

	if err != nil {
		panic(err)
	}
}

func (DSet *Settings) GetConfig(elog *log.Logger) {

	config, err := ini.LoadFile("config.ini") // Loading default config file and checking for various errors.

	if os.IsNotExist(err) {
		
		log.Println("Config.ini does not exist, creating it") //We can not live without config. We could, in theory, but writing default config if none exist can wait
		DSet.WriteConfig(elog)
		return
	}

	if err != nil && !os.IsNotExist(err) {
		elog.Panicln(err) //Oh, something is broken beyond my understanding. Sorry.
	}

	//Getting stuff from config, overwriting hardwired defaults when needed

	Key, ok := config.Get("main", "key")

	if !ok || Key == "" {
		elog.Println("'key' variable missing from 'main' section. It is vital for server-side filtering") //Empty key or key does not exist. Derpibooru works with this, but default image filter filters too much. Use key to set your own!
	}

	DSet.Key = Key

	Q_temp, _ := config.Get("main", "queue_depth")

	if Q_temp != "" {
		DSet.QDepth, err = strconv.Atoi(Q_temp)

		if err != nil {
			elog.Fatalln("Wrong configuration: Depth of the buffer queue is not a number")

		}
	}

	ID_temp, _ := config.Get("main", "downdir")

	if ID_temp != "" {
		DSet.ImgDir = ID_temp
	}
}
