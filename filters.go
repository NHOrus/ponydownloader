package main

type filtrator func(ImageCh) ImageCh

var filters []filtrator

func filterInit(opts *Options) {
	if !opts.Filter {
		lInfo("Filter is off")
		filters = append(filters, nopFilter)
		return
	}
	lInfo("Filter is on")
	filters = append(filters, scoreFilterGenerator(opts.Score))
}

func nopFilter(in ImageCh) ImageCh {
	return in
}

func scoreFilterGenerator(score int) filtrator {
	return func(in ImageCh) ImageCh {
		out := make(ImageCh, 1)
		go func() {
			for imgdata := range in {

				if imgdata.Score >= score {
					out <- imgdata
					continue
				}
				lInfo("Filtering " + imgdata.Filename)
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
