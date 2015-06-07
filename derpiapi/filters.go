package derpiapi

import (
	"log"
)

//FilterSet describes parameters upon which we need to cut off unneeded images
type FilterSet struct {
	Filterflag bool //Enable or disalbe filter
	Scrfilter  int //minimal score upon which to filter things
}

//FilterChannel cuts off unneeded images
func FilterChannel(inchan <-chan Image, outchan chan<- Image, fset FilterSet) {

	for {

		imgdata, more := <-inchan

		if !more {
			close(outchan)
			return //Why make a bunch of layers of ifs if one can just end it all?
		}

		if filterByScore(imgdata, fset) {
			outchan <- imgdata
		} else {
			log.Println("Filtering " + imgdata.Filename)
		}
	}
}

func filterByScore(img Image, fset FilterSet) bool {
	return !fset.Filterflag || (fset.Filterflag && img.Score >= fset.Scrfilter)
}
