package main

import (
	"os"
	"os/signal"
)

var (
	interruptParse = make(chan os.Signal, 1)
	interruptDL    = make(chan os.Signal, 1)
)

func init() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		interruptParse <- os.Interrupt
		close(interruptDL)
		<-sig
		lDone("Program interrupted by user's command")
		os.Exit(0)
	}()
}

func (imgchan ImageCh) interrupt() (outch ImageCh) {
	outch = make(ImageCh)
	go func() {
		for {
			select {
			case <-interruptDL:
				close(outch)
				return
			case img, ok := <-imgchan:
				if !ok {
					close(outch)
					imgchan = nil
					return
				}
				outch <- img
			}
		}
	}()

	return outch
}

func isParseInterrupted() bool {
	select {
	case <-interruptParse:
		return true
	default:
		return false
	}
}
