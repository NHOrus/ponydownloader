ponydownloader
==============

Ponydownloader seeks to provide useful command-line tool to download images from [Derpibooru](http://derpiboo.ru) using provided API

Currently ponydownloader provides three bits of functionality: download by image id, download by tag and filter images you download by their score.

Usage
-----



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
