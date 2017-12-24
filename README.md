ponydownloader
==============

[![forthebadge](http://forthebadge.com/images/badges/fuck-it-ship-it.svg)](http://forthebadge.com) [![forthebadge](http://forthebadge.com/images/badges/oooo-kill-em.svg)](http://forthebadge.com) [![forthebadge](http://forthebadge.com/images/badges/uses-badges.svg)](http://forthebadge.com)

WARNING
-------

This app is under limited support due to loss of interest in ponies. Sorry for those people who may have wanted more.

---

Ponydownloader seeks to provide useful command line tool to download images from [Derpibooru](https://derpibooru.org) en-masse.

Ponydownloader can download by image ID, download by tag and filter out images by their score and/or number of favorites.

Usage
-----


#### Simple usage example:
```bash
./ponydownloader 415147
```

#### Complex usage example:
```bash
./ponydownloader --score 50 -p 3 -n 7 -t "princess luna, safe" --logfilter
```

Output: 

```
Derpibooru.org Downloader, version 0.11.0

Happened at 2017/12/20 03:44:27 Program start
Happened at 2017/12/20 03:44:27 Processing tags princess luna, safe
Happened at 2017/12/20 03:44:27 Starting worker
Happened at 2017/12/20 03:44:27 Worker started; reading channel
Happened at 2017/12/20 03:44:27 Searching as https://derpibooru.org/search.json?sbq=princess+luna%2C+safe
Happened at 2017/12/20 03:44:27 Searching page 3
Happened at 2017/12/20 03:44:27 Saving as 1605729.jpeg
Happened at 2017/12/20 03:44:27 Saving as 1605666.jpeg
Happened at 2017/12/20 03:44:27 Saving as 1605673.jpeg
Happened at 2017/12/20 03:44:27 Saving as 1606115.png
Happened at 2017/12/20 03:44:27 Searching page 4
Happened at 2017/12/20 03:44:28 Skipping: no-clobber
Happened at 2017/12/20 03:44:28 Saving as 1605591.png
Happened at 2017/12/20 03:44:28 Downloaded 1056255 bytes in 0.35s, speed 2.86 MiB/s
Happened at 2017/12/20 03:44:28 Saving as 1605586.png
Happened at 2017/12/20 03:44:28 Filtering  1605358.jpeg
Happened at 2017/12/20 03:44:28 Downloaded 381653 bytes in 0.54s, speed 689.22 KiB/s
Happened at 2017/12/20 03:44:28 Saving as 1605449.jpeg
Happened at 2017/12/20 03:44:28 Downloaded 267790 bytes in 0.56s, speed 464.06 KiB/s
...
```
#### Usage:

 - `-t,	--tag`		Tags to search and download images with. Same rules as Depribooru search: Downloaded images must have all tags passed to this flag.
 - `-k,	--key`		API key to use for Derpibooru access under your account. Can be found in your [account settings](https://derpibooru.org/users/edit). Pass once, better yet put in configuration file. Once passed, gets saved in configuration file.
 - `--dir`			Target directory to save images. Default directory - `img` under current directory. To explicitely save into current directory, pass `--dir=""`
 - `-q	--queue`	Queue Depth, how many images should wait to be downloaded. Default - 50, one page of search. Best leave default.  

#### Limiting amount of downloaded images:
 - `-p, --startpage`	Start downloading from p-th page of search, skipping 50*p images.
 - `-n, --stoppage`	Stop downloading on n-th page, downloading 50*n images at most.

Ponydownloader ignores `n` less than `p` and downloads exactly 50 images when `p` is equal `n`.

#### Filtering options

 - `--score` 		Minimal score image must possess to be downloaded
 - `--faves`		Minimal amount of favorites image must possess to be downloaded
 - `--logfilter`	Note that images were filtered out from download queue

Those options exists to skip low-quality images. If both present, images must possess both score and number of favorites to be downloaded. `logfilter`, by default set to true, makes a note in `events.log` of all discarded images.

#### Notes

Ability to download by tags is not exclusive with bare image IDs: given both, all images with tags and all images with passed IDs would be downloaded.  
Ponydownloader writes a log into `events.log`, containing errors, ID of downloaded and filtered out images, search pages processed and some other helpful information. This file is automatically rotated, with a hardcoded limit of 1Mb per file and 10 files total. That allow to keep log for about ~15k images downloaded.

At start, ponydownloader reads `config.ini`, command line, then writes all set static parameters - `key`, `dir`, `queue` and `logfilter` into it, creating new one if config.ini didn't exist previously.  
Derpibooru provides significant capability to filter out images server-side, for example spoilers or explicit ones. Passing key allows one to enable them and fine-tune some additional settings, instead of passing tags with each request.

## How to install ponydownloader

There are two way to install this program: definitely working/developer way and most likely working/path of less resistance.

##### Path of less resistance.

If you trust me, Github Releases should contain latest release version, [here](https://github.com/NHOrus/ponydownloader/releases). Copy ponydownloader-*your_os*-*your_architecture* in a folder you want to run it from or somewhere in your path. Enjoy. You may need to pass correct values of target directory for downloads and your API key through CLI or config.ini

Ponydownloader needs to run from terminal and does not run interactively.

Binaries may be outdated. My cross-compilation system may work not as well as intended. Binaries may be malicious, knowingly or unknowingly.

##### Path of more compilation (any system with Go installed, console)

```
 go get -d github.com/NHOrus/ponydownloader.git
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

Sample config.ini
----------

```config.ini
key		=	// your derpibooru.org key
downdir		= img	// in this directory your images would be saved
queue_depth	= 50	// depth of queue of images, waiting for download. Default value - one search page
logfilter	= false	// should app write ID discarded by filters images in log
```
