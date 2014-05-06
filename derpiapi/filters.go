package derpiapi

import (
	"log"
)

type FilterSet struct {
	Filterflag bool
	Scrfilter  int
}

func FilterChannel(inchan <-chan Image, outchan chan<- Image, fset FilterSet) {

	for {

		imgdata, more := <-inchan

		if !more {
			close(outchan)
			return //Why make a bunch of layers of ifs if one can just end it all?
		}

		if !fset.Filterflag || (fset.Filterflag && imgdata.Score >= fset.Scrfilter) {
			outchan <- imgdata
		} else {
			log.Println("Filtering " + imgdata.Filename)
		}
	}
}
