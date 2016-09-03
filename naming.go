// Simple command line utility for the manipulation of columned ASCII data files
package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	valueBraceOpen     string = "["
	valueBraceClose    string = "]"
	operationDelimiter string = "."
)

// spectool knows only nanometers and electron-volts
func isAnyKnownSuffix(s string) bool {
	knownSuffixes := []string{".nm", ".ev"}
	for _, sfx := range knownSuffixes {
		if s == sfx {
			return true
		}
	}
	return false
}

// addPreSuffix adds pre-suffix or extension to filename or changes the existing one
// For example:
// addPreSuffix("data1.txt", "ev") == data1.ev.txt
// addPreSuffix("data1.nm", "ev") == data1.ev
// addPreSuffix("data1.nm.txt", "ev") == data1.ev.txt
func addPreSuffix(fname, newSfx string) string {
	newSfx = operationDelimiter + newSfx
	ext := filepath.Ext(fname) // Extension with dot
	if strings.IndexAny(ext, valueBraceClose) >= 0 {
		return fname + newSfx
	}
	subFname := strings.TrimSuffix(fname, ext)
	if isAnyKnownSuffix(ext) {
		if ext == newSfx {
			return fname
		}
		return subFname + newSfx
	}
	sfx := filepath.Ext(subFname)
	if isAnyKnownSuffix(sfx) {
		if sfx == newSfx {
			return fname
		}
		if strings.IndexAny(sfx, valueBraceClose) >= 0 {
			return subFname + newSfx + ext
		}
		return strings.TrimSuffix(subFname, sfx) + newSfx + ext
	}
	return subFname + newSfx + ext
}

// AddPrePreSuffix adds per-pre-suffix to filename related to any operation made on the file
// For example:
// addPrePreSuffix("data1.txt", "div(4.54)") == data1.div(4.54).txt
// addPrePreSuffix("data1.nm", "div(4.54)") == data1.div(4.54).nm
// addPrePreSuffix("data1.add(3).nm.txt", "div(4.54)") == data1.add(3).div(4.54).nm.txt
func AddPrePreSuffix(fname, newSfx string) string {
	newSfx = "." + newSfx
	ext := filepath.Ext(fname)
	if strings.IndexAny(ext, valueBraceClose) >= 0 {
		return fname + newSfx
	}
	subFname := strings.TrimSuffix(fname, ext)
	if isAnyKnownSuffix(ext) {
		return subFname + newSfx + ext
	}
	subSfx := filepath.Ext(subFname)
	if isAnyKnownSuffix(subSfx) {
		return strings.TrimSuffix(subFname, subSfx) + newSfx + subSfx + ext
	}
	return subFname + newSfx + ext
}

// AddSpOpSuffix adds info about operation involving a spectrum in a file name
func (s *Spectrum) AddSpOpSuffix(op, fname string) {
	sfx := op + valueBraceOpen + fname + valueBraceClose
	s.fname = AddPrePreSuffix(s.fname, sfx)
}

// AddNumOpSuffix adds info about operation involving a number in a file name
func (s *Spectrum) AddNumOpSuffix(op string, num float64) {
	sfx := fmt.Sprintf("%s%s%v%s", op, valueBraceOpen, num, valueBraceClose)
	s.fname = AddPrePreSuffix(s.fname, sfx)
}

// AddOpSuffix adds info about operation involving a number in a file name
func (s *Spectrum) AddOpSuffix(op, val string) {
	sfx := fmt.Sprintf("%s%s%s%s", op, valueBraceOpen, val, valueBraceClose)
	s.fname = AddPrePreSuffix(s.fname, sfx)
}
