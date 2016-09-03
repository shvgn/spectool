// Simple command line utility for the manipulation of columned ASCII data files
package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shvgn/xy"
)

// Spectrum handles data with its name. xy is for the original numeric XY data, dir
// stores the path of the original directory of the file and fname is a new filename.
type Spectrum struct {
	xy    *xy.XY
	dir   string
	fname string
}

// NewSpectrum creates new Spectrum from a file containing data with
// optional column numbers for Y alone or X and Y.
func NewSpectrum(fpath string, cols ...int) (*Spectrum, error) {
	data, err := xy.FromFile(fpath, cols...)
	if err != nil {
		return nil, err
	}
	dir, fname := filepath.Split(fpath)
	spec := &Spectrum{xy: data, fname: fname, dir: dir}
	return spec, nil
}

// String representation
func (s *Spectrum) String() string {
	var buf bytes.Buffer
	// if sw.dir == "" {
	// 	buf.WriteString(fmt.Sprintf("Directory: %s\n", "."))
	// } else {
	// 	buf.WriteString(fmt.Sprintf("Directory: %s\n", sw.dir))
	// }
	// buf.WriteString(fmt.Sprintf("Filename: %s\n", sw.fname))
	buf.WriteString(s.xy.String())
	buf.WriteString("\n")
	return buf.String()
}

// WriteFile writes the data into a new corresponding file
func (s *Spectrum) WriteFile(path string, fmt string, perm os.FileMode) error {
	var strFunc func() string
	if fmt == "ascii" {
		strFunc = s.String
	}
	// if fmt == "tsv"     { strFunc = s.TSVString }
	// if fmt == "csv"     { strFunc = s.CSVString }
	// if fmt == "matlab"  { strFunc = s.MATLABString }
	// if fmt == "json"    { strFunc = s.JSONString }
	return ioutil.WriteFile(path, []byte(strFunc()), perm)
}
