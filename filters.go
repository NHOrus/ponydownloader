package main

import "log"

//FilterChannel cuts off unneeded images
func FilterChannel(in ImageCh) (out ImageCh) {
	if !opts.Filter {
		log.Println("Filter is off")
		return in
	}
	log.Println("Filter is on")
	out = make(ImageCh, 1)
	go func() {
		for imgdata := range in {

			if imgdata.Score >= opts.Score {
				out <- imgdata
				continue
			}
			log.Println("Filtering " + imgdata.Filename)
		}
		close(out)
	}()
	return out
}
