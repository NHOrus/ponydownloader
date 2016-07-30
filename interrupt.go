package main

import (
	"os"
	"os/signal"
)

var (
	interrupter = make(chan os.Signal, 1)
)

func init() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		close(interrupter)
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
			case img, ok := <-imgchan:
				if !ok {
					close(outch)
					imgchan = nil
					return
				}
				outch <- img
			default:
				if isInterrupted() {
					close(outch)
					return
				}
			}
		}
	}()

	return outch
}

func isInterrupted() bool {
	select {
	case <-interrupter:
		return true
	default:
		return false
	}
}
