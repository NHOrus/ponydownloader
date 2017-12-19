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

const version = "0.11.0"

//Config is concrete and stored in configuration file
type Config struct {
	ImageDir   string `long:"dir" description:"Target Directory" default:"img" ini-name:"downdir"`
	QDepth     int    `short:"q" long:"queue" description:"Length of the queue buffer" default:"50" ini-name:"queue_depth"`
	Key        string `short:"k" long:"key" description:"Derpibooru API key" ini-name:"key"`
	LogFilters Bool   `long:"logfilter" optional:" " optional-value:"true" description:"Enable logging of filtered images" ini-name:"logfilter"`
}

//FlagOpts are runtime boolean flags
type FlagOpts struct {
	UnsafeHTTPS bool `long:"unsafe-https" description:"Disable HTTPS security verification"`
}

//FiltOpts are filtration parameters
type FiltOpts struct {
	Score  int  `long:"score" description:"Filter option, minimal score of image for it to be downloaded"`
	Faves  int  `long:"faves" description:"Filter option, minimal amount of people who favored image for it to be downloaded"`
	ScoreF bool `no-flag:" "`
	FavesF bool `no-flag:" "`
}

//TagOpts are options relevant to searching by tags
type TagOpts struct {
	Tag       string `short:"t" long:"tag" description:"Tag to download"`
	StartPage int    `short:"p" long:"startpage" description:"Starting page for search" default:"1"`
	StopPage  int    `short:"n" long:"stoppage" description:"Stopping page for search, default - parse all search pages"`
}

//Options provide program-wide options. At maximum, we got one persistent global and one short-living copy for writing in config file
type Options struct {
	*Config
	*FlagOpts
	*FiltOpts
	*TagOpts
	Args struct {
		IDs []int `description:"Image IDs to download" optional:"yes"`
	} `positional-args:"yes"`
}

func getOptions() (opts *Options, args []string) {
	opts = new(Options)
	err := flag.IniParse("config.ini", opts)
	if err != nil {
		switch err.(type) {
		default:
			lFatal(err)
		case *os.PathError:
			lWarn("config.ini not found, using defaults")
		}
	}
	inisets := *opts.Config //copy value instead of reference - or we will get no results later

	args, err = flag.Parse(opts)
	if err != nil {

		switch err.(*flag.Error).Type {
		case flag.ErrHelp:
			os.Exit(0) //Why fall through when asked for help? Just exit with suggestion
		case flag.ErrUnknownFlag:
			fmt.Println("Use --help to view all available options")
			os.Exit(0)
		default:
			lFatal("Can't parse flags: ", err)
		}
	}
	dirF := false
	for _, arg := range os.Args {
		if strings.Contains(arg, "--score") {
			opts.ScoreF = true
		}
		if strings.Contains(arg, "--faves") {
			opts.FavesF = true
		}

		if strings.Contains(arg, "--dir") {
			dirF = true
		}
	}

	if !dirF {
		opts.Config.ImageDir = inisets.ImageDir
	}

	if opts.Config.isEqual(&inisets) { //If nothing to write, no double-writing files
		return
	}

	inifile, err := os.OpenFile("config.ini", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)

	if err != nil {
		lFatal("Could  not create configuration file")
	}

	defer func() {
		err = inifile.Close()
		if err != nil {
			lFatal("Could  not close configuration file")
		}
	}()

	err = opts.Config.prettyWriteIni(inifile)

	if err != nil {
		lFatal("Could  not write in configuration file")
	}
	return
}

//prettyWriteIni Uses tabwriter to make pretty ini file with
func (sets *Config) prettyWriteIni(inifile io.Writer) error {
	tb := tabwriter.NewWriter(inifile, 10, 8, 0, ' ', 0) //Tabs! Elastic! Pretty!

	fmt.Fprintf(tb, "key \t= %s\n", sets.Key)
	fmt.Fprintf(tb, "queue_depth \t= %s\n", strconv.Itoa(sets.QDepth))
	fmt.Fprintf(tb, "downdir \t= %s\n", sets.ImageDir)
	fmt.Fprintf(tb, "logfilter \t= %t\n", sets.LogFilters)

	return tb.Flush() //Returns and passes error upstairs
}

//isEqual compares only options I want to preserve across launches.
func (sets *Config) isEqual(b *Config) bool {
	if b == nil {
		return false
	}
	if sets.ImageDir == b.ImageDir &&
		sets.QDepth == b.QDepth &&
		sets.Key == b.Key &&
		sets.LogFilters == b.LogFilters {
		return true
	}
	return false
}
