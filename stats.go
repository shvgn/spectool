// Simple command line utility for the manipulation of columned ASCII data files
package main

import (
	"fmt"

	"github.com/shvgn/xy"
)

// A spectrum stats
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

func stats(s *xy.XY) *Stats {
	st := &Stats{}
	st.area = s.Area()
	st.maxpos, st.maxheight = s.MaxY()
	st.fwhm = s.FWHM(st.maxpos)
	// notImplemented() // FIXME
	return st
}
