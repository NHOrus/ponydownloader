package main

//FilterSet describes parameters upon which we need to cut off unneeded images
type FilterSet struct {
	MinScore int //minimal score upon which to filter things
}

//FilterChannel cuts off unneeded images
func FilterChannel(in <-chan Image, out chan<- Image, fset FilterSet) {

	for imgdata := range in {

		if imgdata.Score >= fset.MinScore {
			out <- imgdata
			continue
		}
		elog.Println("Filtering " + imgdata.Filename)
	}
	close(out)
}
