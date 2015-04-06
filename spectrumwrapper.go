// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code in any way you wish.
//
// Evgeny Shevchenko
// shvgn@protonmail.ch
// 2015

package main

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/shvgn/spectrum"
)

// Type to handle data with its name. s is for the original data Spectrum, dir
// stores the original directory of the file and fname is a new filename.
type SpectrumWrapper struct {
	s     *spectrum.Spectrum
	dir   string
	fname string
}

// Get new SpecWrapper from a file containing data with optional column numbers
// for Y alone or X and Y.
func NewSpecWrapper(fpath string, cols ...int) (*SpectrumWrapper, error) {
	s, err := spectrum.SpectrumFromFile(fpath, cols...)
	if err != nil {
		return nil, err
	}
	dir, fname := filepath.Split(fpath)
	sw := &SpectrumWrapper{s: s, fname: fname, dir: dir}
	return sw, nil
}

// Write the data into a new corresponding file
func (sw *SpectrumWrapper) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Directory: %s\n", sw.dir))
	buf.WriteString(fmt.Sprintf("Filename: %s\n", sw.fname))
	buf.WriteString(sw.s.String())
	buf.WriteString("\n")
	return buf.String()
}

// Write the data into a new corresponding file
func (sw *SpectrumWrapper) Write() {
	notImplemented()
	// writes the data sw.s to the file path sw.dir+sw.fname
}
