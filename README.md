ponydownloader
==============

[![forthebadge](http://forthebadge.com/images/badges/fuck-it-ship-it.svg)](http://forthebadge.com)

WARNING
-------

This app under limited support due to loss of interest into ponies. Sorry for those people who may have wanted more.

---

Ponydownloader seeks to provide useful command-line tool to download images from [Derpibooru](https://derpibooru.org) using provided API

Currently ponydownloader provides three bits of functionality: download by image id, download by tag and filter images you download by their score.

Usage
-----

Currently ponydownloader got two main modes of usage: download single image by id and batch download of images by tag.

After invocation, ponydownloader would read `config.ini` if it exist or write default one in current working directory. Then it would write completed actions in file `event.log` and write image in a directory specified in configuration file.

To download single image, simply invoke ponydownloader with image id as argument.

To download all images with desired flag , invoke `ponydownloader -t <flag>`

One can manipulate start and stop pages for limiting amount of downloaded images or skipping images already present: `-p <star page>` `-n <stop page>` . Due to concurrent design of ponydownloader and insufficient documentation of Derpibooru API, queue may contain more images that response page from site and images ponydownloader currently saves may be from significantly earlier that page program declares as one being processed.

Currently only filter available is filter by score: `--score ` then minimal score to accept.

Optional flag `-k` defines API key to use and overrides said key from configuration file. Key in configuration file gets rewritten. Derpibooru provides significant capability to exclude undesirable images server-side and API key allows one to switch from default settings to currently selected personal rule. One of the way to get API key is to look at your [account settings](https://derpibooru.org/users/edit)

Full list of flags returned by `--help` command

#### Simple usage example:
```bash
./ponydownloader 415147
```

#### Complex usage example:
```bash
./ponydownloader -f --score 50 -p 3 -n 7 -t princess+luna%2C+safe
```
At the moment of writing both samples were working, you would get output looking approximately like quote below and images in default directory `img`

```
Derpibooru.org Downloader version 0.6.0
Happened at 2015/11/17 16:08:28 Program start
Happened at 2015/11/17 16:08:28 Processing tags princess+luna%2C+safe
Happened at 2015/11/17 16:08:28 Starting worker
Happened at 2015/11/17 16:08:28 Filter is on
Happened at 2015/11/17 16:08:28 Worker started; reading channel
Happened at 2015/11/17 16:08:28 Searching as https://derpibooru.org/search.json?q=princess+luna%2C+safe
Happened at 2015/11/17 16:08:28 Searching page 3
Happened at 2015/11/17 16:08:28 Filtering 1020863.jpeg
Happened at 2015/11/17 16:08:28 Saving as 1020830.png
Happened at 2015/11/17 16:08:29 Downloaded 1634245 bytes in 0.88s, speed 1.76 MiB/s
Happened at 2015/11/17 16:08:29 Saving as 1020688.png
Happened at 2015/11/17 16:08:29 Downloaded 95298 bytes in 0.06s, speed 1.58 MiB/s
Happened at 2015/11/17 16:08:29 Saving as 1020684.png
Happened at 2015/11/17 16:08:29 Downloaded 117766 bytes in 0.06s, speed 1.92 MiB/s
Happened at 2015/11/17 16:08:29 Saving as 1020682.png
...
```

## How to install ponydownloader

There are two way to install this program: definitely working/developer way and most likely working/path of less resistance.

##### Path of less resistance.

If you trust me, Github Releases should contain latest release version, [here](https://github.com/NHOrus/ponydownloader/releases) Copy ponydownloader-*your_os*-*your_architecture* in a folder you want to run it from or somewhere in your path. Enjoy. You may need to pass correct values of target directory for downloads and your API key through CLI or config.ini

Binaries may be outdated. My cross-compilation system may work not as well as intended. Binaries may be malicious, knowingly or unknowingly.

##### Path of more compilation (any system with Go installed, console)

```
 go get -d git@github.com:NHOrus/ponydownloader.git
 cd $GOPATH/src/github.com/NHOrus/ponydownloader
 go build
 cp config.ini.sample config.ini
 // edit in your key into config.ini or put it once as command line argument
 // run ponydownloader as explained earlier
``` 

###### Explanation:

Make sure your $GOPATH is set and correct and you got git installed and working.

We trust into `go get -d` to get source code of ponydownloader and all dependencies. Also, it updates them. `-d` flag  prevents installation. Alternatively, one can run `go get` to compile all the code and install binary into $GOPATH/bin 

Then we move into directory where source code sleeps endless sleep, compile everything into one binary executable and set it up for running in that same directory with the source. Or we can move it everywhere we want and run locally there.

config.ini
----------

```config.ini
key 		=	  // your derpibooru.org key
downdir 	= img // in this directory your images would be saved
queue_depth = 20  // depth of queue of images, waiting for download. 
``` 

I feel that optimal depth is around 10-20, else there would be slowdown when parser requests next search page from derpibooru and feeds it's content to worker. It depends upon downloading speed and server response time. Ponydownloader downloads search  pages and images simultaneously, one by one.
If any line is empty, program would use default, build-in parameters. Empty `key` would end up with no API key being used.
