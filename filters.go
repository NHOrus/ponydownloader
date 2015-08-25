package main

//FilterChannel cuts off unneeded images
func FilterChannel(in ImageCh) (out ImageCh) {
	//elog.Println("filtering")
	if !opts.Filter {
		//elog.Println("Filter is off")
		return in
	}
	//elog.Println("Filter is on")
	out = make(ImageCh, 1)
	go func() {
		for imgdata := range in {

			if imgdata.Score >= opts.Score {
				out <- imgdata
				continue
			}
			//elog.Println("Filtering " + imgdata.Filename)
		}
		close(out)
	}()
	return out
}
