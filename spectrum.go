// Simple command line utility for the manipulation of columned ASCII data files
package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shvgn/xy"
)

// Type to handle data with its name. s is for the original data Spectrum, dir
// stores the original directory of the file and fname is a new filename.
type Spectrum struct {
	xy    *xy.XY
	dir   string
	fname string
}

// Get new SpecWrapper from a file containing data with optional column numbers
// for Y alone or X and Y.
func NewSpectrum(fpath string, cols ...int) (*Spectrum, error) {
	s, err := xy.FromFile(fpath, cols...)
	if err != nil {
		return nil, err
	}
	dir, fname := filepath.Split(fpath)
	sw := &Spectrum{xy: s, fname: fname, dir: dir}
	return sw, nil
}

// String representation
func (sw *Spectrum) String() string {
	var buf bytes.Buffer
	// if sw.dir == "" {
	// 	buf.WriteString(fmt.Sprintf("Directory: %s\n", "."))
	// } else {
	// 	buf.WriteString(fmt.Sprintf("Directory: %s\n", sw.dir))
	// }
	// buf.WriteString(fmt.Sprintf("Filename: %s\n", sw.fname))
	buf.WriteString(sw.xy.String())
	buf.WriteString("\n")
	return buf.String()
}

// Write the data into a new corresponding file
func (sw *Spectrum) WriteFile(path string, fmt string, perm os.FileMode) error {
	var strFunc func() string
	if fmt == "ascii" {
		strFunc = sw.String
	}
	// if fmt == "tsv"     { strFunc = sw.TSVString }
	// if fmt == "csv"     { strFunc = sw.CSVString }
	// if fmt == "matlab"  { strFunc = sw.MATLABString }
	// if fmt == "json"    { strFunc = sw.JSONString }
	return ioutil.WriteFile(path, []byte(strFunc()), perm)
}
