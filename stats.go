// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code in any way you wish.
//
// Evgeny Shevchenko
// shvgn@protonmail.ch
// 2015

package main

import (
	"fmt"

	spectrum "github.com/shvgn/xy"
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

func stats(s *spectrum.Spectrum) *Stats {
	st := &Stats{}
	st.area = s.Area()
	st.maxpos, st.maxheight = s.MaxY()
	st.fwhm = s.FWHM(st.maxpos)
	// notImplemented() // FIXME
	return st
}
