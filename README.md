ponydownloader
==============

Ponydownloader seeks to provide useful command-line tool to download images from [Derpibooru](http://derpiboo.ru) using provided API

Currently ponydownloader provides three bits of functionality: download by image id, download by tag and filter images you download by their score.

Usage
-----

Currently ponydownloader got two main modes of usage: download single image by id and batch download of images by tag.

After invocation, ponydownloader would read `config.ini` if it exist or write default one in current working directory. Then it would write completed actions in file `event.log` and write image in a directory specified in configuration file.

To download single image, simply invoke ponydownloader with image id as argument.

To download all images with desired flag , invoke `ponydownloader -t <flag>`

One can manipulate start and stop pages for limiting amount of downloaded images or skipping images already present: `-p <star page>` `-np <stop page>` . Due to concurrent design of ponydownloader and insufficient documentation of Derpibooru API, queue may contain more images that response page from site and images ponydownloader currently saves may be from significantly earlier that page program declares as one being processed.

To filter images one need to explicitly declare desire to do so by setting up flag `-filter` and then declare parameter and it's value.
Currently only filter available is filter by score: `-scr <minimal score to accept`

Optional flag `-k` defines API key to use and overrides said key from configuration file. Derpibooru provides significant capability to exclude undesirable images server-side and API key allows one to switch from default settings to currently selected personal rule. One of the way to get API key is to look at your [account settings](https://derpiboo.ru/users/edit)

#### Simple usage example:
```bash
./ponydownloader 415147
```

#### Complex usage example:
```bash
./ponydownloader -filter -scr 15 -p 3 -np 7 -t princess+luna%2C+safe
```
At the moment of writing both samples were working, you would get output looking approximately like quote below and images in default directory `img`

```
Derpiboo.ru Downloader version 0.2.0
Happens at 2014/05/07 16:17:54 Program start
Happens at 2014/05/07 16:17:54 Processing tags princess+luna%2C+safe
Searching as http://derpiboo.ru/search.json?nofav=&nocomments=&q=princess+luna%2C+safe
Searching page 3
Happens at 2014/05/07 16:17:54 Starting worker
Worker started; reading channel
Happens at 2014/05/07 16:17:58 Saving as 617803.sun_and_moon_by_wildberry_poptart-d7h5e3b.png.png
Happens at 2014/05/07 16:17:59 Saving as 617794.00_10_30_33_file.png
...
Happens at 2014/05/07 16:18:11 Filtering 617461.Untitled-11.jpg.jpg
...
```

## How to install ponydownloader

There are two way to install this program: definitely working/developer way and most likely working/path of less resistance.

##### Path of less resistance.

If you trust me, this git repository contains binaries in `bin` directory. This especially true for releases. They should just work. Copy ponydownloader-*your_os*-*your_architecture* in a folder you want to run it from or somewhere in your path. Enjoy.

Binaries may be outdated. My cross-compilation system may work not as well as intended. Binaries may be malicious, knowingly or unknowingly.

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
[main]
key =  //your derpiboo.ru key
downdir = img // in this directory your images would be saved
queue_depth = 20 // depth of queue of images, waiting for download. 
``` 
I feel that optimal depth is around 10-20, else there would be slowdown when parser requests next search page from derpibooru and feeds it's content to worker. It depends upon downloading speed and server response time. Ponydownloader downloads search  pages and images simultaneously, one by one.
If any line is empty, program would use default, build-in parameters. Empty `key` would end up with no API key being used.
