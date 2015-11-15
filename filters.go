package main

import "log"

type Filter func(ImageCh) ImageCh

var filters []Filter

func filterInit(opts Options) {
	if !opts.Filter {
		log.Println("Filter is off")
		filters = append(filters, nopFilter)
		return
	}
	log.Println("Filter is on")
	filters = append(filters, scoreFilterGenerator(opts))
}

func nopFilter(in ImageCh) ImageCh {
	return in
}

func scoreFilterGenerator(option Options) Filter {
	return func(in ImageCh) ImageCh {
		out := make(ImageCh, 1)
		go func() {
			for imgdata := range in {

				if imgdata.Score >= option.Score {
					out <- imgdata
					continue
				}
				log.Println("Filtering " + imgdata.Filename)
			}
			close(out)
		}()
		return out
	}
}

//FilterChannel cuts off unneeded images
func FilterChannel(in ImageCh) (out ImageCh) {
	out = in
	for _, filter := range filters {
		out = filter(out)
	}
	return
}
