ponydownloader
==============

Attempt to learn go and exploit derpiboo.ru public api to batch download files concurrently.

Current version gets a number of image and saves it in current directory

There are ways to go, but it provides minimal functionality as of now

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

Secion [main]

key => your derpiboo.ru key

./bin directory
---------------

Attempt at cross-compilation. Should contain statically linked binaries for all major systems, just in case. Try it if you trust me and author of cross-compilation scripts.

