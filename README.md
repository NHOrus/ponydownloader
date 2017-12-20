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

With default settings, images are downloaded into directory `img` under current directory

After invocation, ponydownloader would attempt to read `config.ini` in current directory and create default one if it does not exist. Then it would write completed actions in file `event.log` and write image in a directory specified in configuration file. To protect innocent and prevent excessive accumulation of logs, rotation is implemented, currenlty caps at 1 Mb per file and 10 logfiles total.

To download single image, simply invoke ponydownloader with image ID as argument. For multiple images, separate their IDs by space.

To download all images with desired tags, invoke `ponydownloader -t "tag A, tag B,.."`

One can manipulate start and stop pages for tag download for limiting amount of downloaded images or skipping images already present: `-p <star page>` `-n <stop page>`. If stop page is less that start page, only images from start page are downloaded. Derpibooru provides results in fifty per page.

Currently available filters are:
-  filter by score: `--score ` then minimal score to accept.
-  filter by favorites: `--faves ` then minimal number of people who favored the to accept.

Optional flag `-k` defines API key to use and overwrites said key from configuration file. You may want to back up your old key in this case. One of the way to get API key is to look at your [account settings](https://derpibooru.org/users/edit)
Derpibooru provides significant capability to exclude undesirable images server-side and API key allows one to switch from default settings to currently selected personal rule.

Optional flag `--logfilter` turns on detailed log of each image discarded before download. It's saved into config file. To turn off, pass `--logfilter=false`

Full list of flags returned by `--help` command.

## How to install ponydownloader

There are two way to install this program: definitely working/developer way and most likely working/path of less resistance.

##### Path of less resistance.

If you trust me, Github Releases should contain latest release version, [here](https://github.com/NHOrus/ponydownloader/releases). Copy ponydownloader-*your_os*-*your_architecture* in a folder you want to run it from or somewhere in your path. Enjoy. You may need to pass correct values of target directory for downloads and your API key through CLI or config.ini

Ponydownloader needs to run from terminal and does not run interactively.

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
key			=		// your derpibooru.org key
downdir		= img	// in this directory your images would be saved
queue_depth	= 50	// depth of queue of images, waiting for download. Default value - one search page
logfilter	= false	// should app write ID discarded by filters images in log
```

If any line is empty, program would use default, build-in parameters and writes them in config, when appropriate. Empty `key` would end up with no API key being used.
