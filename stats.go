// Simple command line utility for the manipulation of columned ASCII data files
package main

import (
	"fmt"

	"github.com/shvgn/xy"
)

// Stats represents spectrum properties: area under the curve,
// the position of the maximum point and its height, and full with at half-maximum
// regardless of how many peaks are in the spectrum
type Stats struct {
	area      float64
	maxpos    float64
	maxheight float64
	fwhm      float64
}

func (st *Stats) String() string {
	return fmt.Sprintf("area: %f\tmaxpos: %f\tmaxheight: %f\tfwhm: %f",
		st.area, st.maxpos, st.maxheight, st.fwhm)
}

func stats(data *xy.XY) *Stats {
	st := &Stats{}
	st.area = data.Area()
	st.maxpos, st.maxheight = data.MaxY()
	st.fwhm = data.FWHM(st.maxpos)
	// notImplemented() // FIXME
	return st
}
