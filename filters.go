package main

type filtrator func(ImageCh) ImageCh

var filters []filtrator

//If filter isn't on, skip. If any of filter parameters is given, filtration is on
func filterInit(opts *FiltOpts) {
	if !opts.Filter {
		lInfo("Filter is off")
		filters = append(filters, nopFilter) //First class function and array of them
		return
	}
	lInfo("Filter is on")
	filters = append(filters, scoreFilterGenerator(opts.Score))
}

//Do nothing
func nopFilter(in ImageCh) ImageCh {
	return in
}

func scoreFilterGenerator(score int) filtrator {
	return func(in ImageCh) ImageCh {
		out := make(ImageCh, 1) //minimal buffer, so nothing should grind to a halt, hopefully, when image is consumed
		go func() {
			for imgdata := range in {

				if imgdata.Score >= score { //Capturing score inside lambda, to prevent passing it around each invocation
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
