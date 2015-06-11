package settings

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vaughan0/go-ini" //We need some simple way to parse ini files, here it is, externally.
)

//Settings contain configuration used in ponydownloader
type Settings struct {
	QDepth int
	ImgDir string
	Key    string
}

//SetLog sets up logfile as I want it to: Copy to event.log, copy to commandline
func SetLog() (retlog *log.Logger, logfile *os.File) {

	logfile, err := os.OpenFile("event.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644) //file for putting errors into

	if err != nil {
		panic(err)
	}
	//elog - recoverable errors. log - just things that happen
	retlog = log.New(io.MultiWriter(logfile, os.Stderr), "Errors at ", log.LstdFlags) //Setting stuff for our logging: both errors and events.

	log.SetPrefix("Happens at ")
	log.SetFlags(log.LstdFlags)
	log.SetOutput(io.MultiWriter(logfile, os.Stdout)) //we write in file and stdout
	log.Println("Program start")

	return
}

//WriteConfig writes default, presumably sensible configuration into file.
func (WSet Settings) WriteConfig(elog log.Logger) {
	config, err := os.OpenFile("config.ini", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer func() {
		err = config.Close()
		if err != nil {
			elog.Fatalln("Could  not close configuration file")
		}
	}()

	if err != nil {
		elog.Fatalln("Could  not create configuration file")
	}

	defset := []string{"key = " + WSet.Key, "queue_depth = " + strconv.Itoa(WSet.QDepth), "downdir = " + WSet.ImgDir}

	_, err = fmt.Fprintln(config, strings.Join(defset, "\n"))

	if err != nil {
		elog.Fatalln("Could  not write in configuration file")
	}
}

//GetConfig gets configuration from the file, presuming it exist
func (WSet *Settings) GetConfig(elog log.Logger) {

	config, err := ini.LoadFile("config.ini") // Loading default config file and checking for various errors.

	if os.IsNotExist(err) {

		log.Println("Config.ini does not exist, creating it") //We can not live without config. We could, in theory, but writing default config if none exist can wait
		WSet.WriteConfig(elog)
		return
	}

	if err != nil && !os.IsNotExist(err) {
		elog.Panicln(err) //Oh, something is broken beyond my understanding. Sorry.
	}

	//Getting stuff from config, overwriting hardwired defaults when needed

	Key, ok := config.Get("", "key")

	if !ok || Key == "" {
		log.Println("'key' variable missing from 'main' section. It is vital for server-side filtering") //Empty key or key does not exist. Derpibooru works with this, but default image filter filters too much. Use key to set your own!
	}

	WSet.Key = Key

	QTemp, _ := config.Get("", "queue_depth")

	if QTemp != "" {
		WSet.QDepth, err = strconv.Atoi(QTemp)

		if err != nil {
			elog.Fatalln("Wrong configuration: Depth of the buffer queue is not a number")

		}
	}

	IDTemp, _ := config.Get("", "downdir")

	if IDTemp != "" {
		WSet.ImgDir = IDTemp
	}
}
