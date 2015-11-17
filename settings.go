package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	flag "github.com/jessevdk/go-flags"
)

//Settings are concrete and stored in configuration file
type Settings struct {
	ImageDir string `long:"dir" description:"Target Directory" default:"img" ini-name:"downdir"`
	QDepth   int    `short:"q" long:"queue" description:"Length of the queue buffer" default:"20" ini-name:"queue_depth"`
	Key      string `short:"k" long:"key" description:"Derpibooru API key" ini-name:"key"`
}

//Flags are runtime boolean flags
type Flags struct {
	Filter  bool `no-flag:" " short:"f" long:"filter" description:"If set, enables client-side filtering of downloaded images"`
	Unsafe  bool `long:"unsafe" description:"If set, trusts in unknown authority"`
	NoHTTPS bool `long:"nohttps" description:"Disable HTTPS and try to download insecurely"`
}

//Options provide program-wide options. At maximum, we got one persistent global and one short-living copy for writing in config file
type Options struct {
	Settings
	Flags
	Tag       string `short:"t" long:"tag" description:"Tag to download"`
	StartPage int    `short:"p" long:"startpage" description:"Starting page for search" default:"1"`
	StopPage  int    `short:"n" long:"stoppage" description:"Stopping page for search, default - parse all search pages"`
	Score     int    `long:"score" description:"Filter option, minimal score of image for it to be downloaded"`
	Args      struct {
		IDs []int `description:"Image IDs to download" optional:"yes"`
	} `positional-args:"yes"`
}

//WriteConfig writes default, presumably sensible configuration into file.
func WriteConfig(opts, iniopts Settings) {

	if opts.compareStatic(&iniopts) { //If nothing to write, no double-writing files
		return
	}

	config, err := os.OpenFile("config.ini", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	if err != nil {
		lFatal("Could  not create configuration file")
		//Need to check if log file is created before config file. Also, work around screams
		//in case if nor log, nor config file could be created.
	}

	defer func() {
		err = config.Close()
		if err != nil {
			lFatal("Could  not close configuration file")
		}
	}()

	tb := tabwriter.NewWriter(config, 10, 8, 0, ' ', 0) //Tabs! Elastic! Pretty!
	fmt.Fprintf(tb, "key \t= %s\n", opts.Key)
	fmt.Fprintf(tb, "queue_depth \t= %s\n", strconv.Itoa(opts.QDepth))
	fmt.Fprintf(tb, "downdir \t= %s\n", opts.ImageDir)

	err = tb.Flush()

	if err != nil {
		lFatal("Could  not write in configuration file")
	}
}

//compareStatic compares only options I want to preserve across launches.
func (a *Settings) compareStatic(b *Settings) bool {
	if a.ImageDir == b.ImageDir &&
		a.QDepth == b.QDepth &&
		a.Key == b.Key {
		return true
	}
	return false
}

//configSetup reads static config from file and runtime options from commandline
//It also preserves static config for later comparsion with runtime to prevent
//rewriting it when no changes are made
func (opts *Options) configSetup() ([]string, Settings) {
	err := flag.IniParse("config.ini", opts)
	if err != nil {
		switch err.(type) {
		default:
			panic(err)
		case *os.PathError:
			lWarn("config.ini not found, using defaults")
		}
	}
	var iniopts = opts.Settings

	args, err := flag.Parse(opts)
	if err != nil {
		flagError := err.(*flag.Error)

		switch flagError.Type {
		case flag.ErrHelp:
			fallthrough
		case flag.ErrUnknownFlag:
			fmt.Println("Use --help to view all available options")
			os.Exit(0)
		default:
			lErr("Can't parse flags: %s\n", err)
			os.Exit(1)
		}
	}

	for _, arg := range os.Args {
		if strings.Contains(arg, "--score") {
			opts.Filter = true
		}
	}
	return args, iniopts
}

func doOptions() ( opts *Options, args []string){
	opts = new(Options)
	args, iniopts := opts.configSetup()
	WriteConfig(opts.Settings, iniopts)
	return
}
