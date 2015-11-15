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

To filter images one need to explicitly declare desire to do so by setting up flag `-filter` and then declare parameter and it's value.
Currently only filter available is filter by score: `--score <minimal score to accept`

Optional flag `-k` defines API key to use and overrides said key from configuration file. Key in configuration file gets rewritten. Derpibooru provides significant capability to exclude undesirable images server-side and API key allows one to switch from default settings to currently selected personal rule. One of the way to get API key is to look at your [account settings](https://derpibooru.org/users/edit)

#### Simple usage example:
```bash
./ponydownloader 415147
```

#### Complex usage example:
```bash
./ponydownloader -f --score 15 -p 3 -n 7 -t princess+luna%2C+safe
```
At the moment of writing both samples were working, you would get output looking approximately like quote below and images in default directory `img`

```
Derpibooru.org Downloader version 0.5.3
Happens at 2015/06/12 21:26:23 Program start
Happens at 2015/06/12 21:26:23 Processing tags princess+luna%2C+safe
Happens at 2015/06/12 21:26:23 Starting worker
Happens at 2015/06/12 21:26:23 Searching as https://derpibooru.org/search.json?&q=princess+luna%2C+safe
Happens at 2015/06/12 21:26:23 Searching page 3
Happens at 2015/06/12 21:26:23 Worker started; reading channel
Happens at 2015/06/12 21:26:23 Saving as 912564.png
Happens at 2015/06/12 21:26:24 Downloaded 336821 bytes in 0.28s, speed 1.15 MiB/s
Happens at 2015/06/12 21:26:24 Saving as 912535.jpeg
Happens at 2015/06/12 21:26:24 Downloaded 133826 bytes in 0.41s, speed 321.54 KiB/s
Happens at 2015/06/12 21:26:24 Saving as 912526.png
Happens at 2015/06/12 21:26:24 Downloaded 311033 bytes in 0.43s, speed 701.77 KiB/s

...
```

## How to install ponydownloader

There are two way to install this program: definitely working/developer way and most likely working/path of less resistance.

##### Path of less resistance.

Simple way does not work any more - all the binaries should not be in git repo, removed completely.

##### Path of more compilation (*NIX, console)

```
 go get -d git@github.com:NHOrus/ponydownloader.git
 cd $GOPATH/src/github.com/NHOrus/ponydownloader
 go build
 cp config.ini.sample config.ini
 // edit in your key into config.ini
 // run ponydownloader as explained earlier
``` 

###### Explanation:

Make sure your $GOPATH is set and correct and you got git installed and working.

We trust into `go get -d` to get source code of ponydownloader and all dependencies. Also, it updates them. `-d` flag  prevents installation. Alternatively, one can run `go get` to compile all the code and install binary into $GOPATH/bin 

Then we move into directory where source code sleeps endless sleep, compile everything into one binary executable and set it up for running in that same directory with the source. Or we can move it everywhere we want and run locally there.

config.ini
----------

```config.ini
key =  //your derpibooru.org key
downdir = img // in this directory your images would be saved
queue_depth = 20 // depth of queue of images, waiting for download. 
``` 

I feel that optimal depth is around 10-20, else there would be slowdown when parser requests next search page from derpibooru and feeds it's content to worker. It depends upon downloading speed and server response time. Ponydownloader downloads search  pages and images simultaneously, one by one.
If any line is empty, program would use default, build-in parameters. Empty `key` would end up with no API key being used.
