package main

type filtrator func(<-chan Image) <-chan Image

var filters []filtrator

//If filter isn't on, skip. If any of filter parameters is given, filtration is on
func filterInit(opts *FiltOpts, enableLog bool) {

	if opts.ScoreF {
		filters = append(filters, filterGenerator(func(i Image) bool { return i.Score >= opts.Score }, enableLog))
	}
	if opts.FavesF {
		filters = append(filters, filterGenerator(func(i Image) bool { return i.Faves >= opts.Faves }, enableLog))
	}
}

func filterGenerator(filt func(Image) bool, enableLog bool) filtrator {
	return func(in <-chan Image) <-chan Image {
		out := make(chan Image)
		go func() {
			for imgdata := range in {

				if filt(imgdata) { //Capturing score inside lambda, to prevent passing it around each invocation
					out <- imgdata
					continue
				}
				lCondInfo(enableLog, "Filtering ", imgdata.Filename)
			}
			close(out)
		}()
		return out
	}
}

//FilterChannel cuts off unneeded images
func FilterChannel(in <-chan Image) (out <-chan Image) {
	out = in
	for _, filter := range filters {
		out = filter(out)
	}
	return
}
