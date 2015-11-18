package main

import (
	"fmt"
	"io"
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

//FlagOpts are runtime boolean flags
type FlagOpts struct {
	NoHTTPS     bool `long:"nohttps" description:"Disable HTTPS and try to download insecurely"`
	UnsafeHTTPS bool `long:"unsafe-https" description:"Disable HTTPS security verification"`
}

//FiltOpts are filtration parameters
type FiltOpts struct {
	Filter bool `no-flag:" " short:"f" long:"filter" description:"If set, enables client-side filtering of downloaded images"`
	Score  int  `long:"score" description:"Filter option, minimal score of image for it to be downloaded"`
}

//TagOpts are options relevant to searching by tags
type TagOpts struct {
	Tag       string `short:"t" long:"tag" description:"Tag to download"`
	StartPage int    `short:"p" long:"startpage" description:"Starting page for search" default:"1"`
	StopPage  int    `short:"n" long:"stoppage" description:"Stopping page for search, default - parse all search pages"`
}

//Options provide program-wide options. At maximum, we got one persistent global and one short-living copy for writing in config file
type Options struct {
	*Settings
	*FlagOpts
	*FiltOpts
	*TagOpts
	Args struct {
		IDs []int `description:"Image IDs to download" optional:"yes"`
	} `positional-args:"yes"`
}

func getOptions() (opts *Options, args []string) {
	opts = new(Options)
	args, inisets := opts.configSetup()
	opts.Settings.checkedWriteIni(inisets)
	return
}

//checkedWriteIni writes default, presumably sensible configuration into file without overwriting old one if unchanged
func (sets *Settings) checkedWriteIni(oldsets *Settings) {
	if sets.isEqual(oldsets) { //If nothing to write, no double-writing files
		return
	}

	inifile, err := os.OpenFile("config.ini", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	if err != nil {
		lFatal("Could  not create configuration file")
	}

	defer func() {
		err = inifile.Close()
		if err != nil {
			lFatal("Could  not close configuration file")
		}
	}()

	err = sets.prettyWriteIni(inifile)

	if err != nil {
		lFatal("Could  not write in configuration file")
	}
}

//prettyWriteIni Uses tabwriter to make pretty ini file with
func (sets Settings) prettyWriteIni(inifile io.Writer) error {
	tb := tabwriter.NewWriter(inifile, 10, 8, 0, ' ', 0) //Tabs! Elastic! Pretty!

	fmt.Fprintf(tb, "key \t= %s\n", sets.Key)
	fmt.Fprintf(tb, "queue_depth \t= %s\n", strconv.Itoa(sets.QDepth))
	fmt.Fprintf(tb, "downdir \t= %s\n", sets.ImageDir)

	return tb.Flush()
}

//compareStatic compares only options I want to preserve across launches.
func (sets *Settings) isEqual(b *Settings) bool {
	if b == nil {
		return false
	}
	if sets.ImageDir == b.ImageDir &&
		sets.QDepth == b.QDepth &&
		sets.Key == b.Key {
		return true
	}
	return false
}

//configSetup reads static config from file and runtime options from commandline
//It also preserves static config for later comparsion with runtime to prevent
//rewriting it when no changes are made
func (opts *Options) configSetup() ([]string, *Settings) {
	err := flag.IniParse("config.ini", opts)
	if err != nil {
		switch err.(type) {
		default:
			lFatal(err)
		case *os.PathError:
			lWarn("config.ini not found, using defaults")
		}
	}
	inisets := *opts.Settings //copy value instead of reference - or we will get no results later

	args, err := flag.Parse(opts)
	checkFlagError(err) //Here we scream if something goes wrong and provide help if something goes meh.

	for _, arg := range os.Args {
		if strings.Contains(arg, "--score") {
			opts.Filter = true
		}
	}
	return args, &inisets
}

func checkFlagError(err error) {
	if err == nil {
		return
	}

	flagError := err.(*flag.Error)

	switch flagError.Type {
	case flag.ErrHelp:
		os.Exit(0) //Why fall through when asked for help? Just exit with suggestion
	case flag.ErrUnknownFlag:
		fmt.Println("Use --help to view all available options")
		os.Exit(0)
	default:
		lFatal("Can't parse flags: %s\n", err)
	}
}
