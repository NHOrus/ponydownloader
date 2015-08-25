package main

//FilterChannel cuts off unneeded images
func FilterChannel(in ImageCh) (out ImageCh) {

	if opts.Filter {
		return in
	}
	out = make(ImageCh)
	go func() {
		for imgdata := range in {

			if imgdata.Score >= opts.Score {
				out <- imgdata
				continue
			}
			elog.Println("Filtering " + imgdata.Filename)
		}
		close(out)
	}()
	return out
}
