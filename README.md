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

How-to compile
--------------

Only thing this program needs for compilation  is a working Go compiler.

Correct first compilation after cloning repository / getting source code some other way is:

Move inside ponydownloader directory.

Make sure your $GOPATH is set and correct.

>go get

It downloads external dependency from github and prepares it for usage

>go build

you get binary with same name as directory with source code, by default - ponydownloader

Run it.

Alternatively, put ponydownloader folder inside $GOPATH/src and do wherever you want

>go build ponydownloader

then put config file in same place. In the future it would write default config.ini on first run, but not yet.

config.ini
----------

[main]

key ==> your derpiboo.ru key

downdir ==> in this directory your images would be saved

queue_depth ==> depth of download buffer. Leave it no less that default: else there would be slowdown when parser requests next search page from derpibooru and feeds it's content to worker

If any line is not changed from default, program would use default, build-in parameters.
Currenly only key is ought to be in configuration file. It may be empty, then program would ignore it's existence. 

./bin directory
---------------

Cross-compiled versions of the program. Should work. At least for me they are working on Windows.
