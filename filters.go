package main

type filtrator func(ImageCh) ImageCh

var filters []filtrator

//If filter isn't on, skip. If any of filter parameters is given, filtration is on
func filterInit(opts *FiltOpts, enableLog bool) {
	if !opts.Filter {
		lCondInfo(enableLog, "Filter is off")
		filters = append(filters, nopFilter) //First class function and array of them
		return
	}
	lCondInfo(enableLog, "Filter is on")
	if opts.ScoreF {
		filters = append(filters, scoreFilterGenerator(opts.Score, enableLog))
	}
	if opts.FavesF {
		filters = append(filters, favesFilterGenerator(opts.Faves, enableLog))
	}
}

//Do nothing
func nopFilter(in ImageCh) ImageCh {
	return in
}

func scoreFilterGenerator(score int, enableLog bool) filtrator {
	return func(in ImageCh) ImageCh {
		out := make(ImageCh)
		go func() {
			for imgdata := range in {

				if imgdata.Score >= score { //Capturing score inside lambda, to prevent passing it around each invocation
					out <- imgdata
					continue
				}
				lCondInfo(enableLog, "Filtering "+imgdata.Filename)
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

func favesFilterGenerator(faves int, enableLog bool) filtrator {
	return func(in ImageCh) ImageCh {
		out := make(ImageCh)
		go func() {
			for imgdata := range in {

				if imgdata.Faves >= faves { //Capturing score inside lambda, to prevent passing it around each invocation
					out <- imgdata
					continue
				}
				lCondInfo(enableLog, "Filtering "+imgdata.Filename)
			}
			close(out)
		}()
		return out
	}
}
