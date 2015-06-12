// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code in any way you wish.
//
// Evgeny Shevchenko
// shvgn@protonmail.ch
// 2015

package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	VALUE_BRACE_OPEN    string = "["
	VALUE_BRACE_CLOSE   string = "]"
	OPERATION_DELIMITER string = "."
)

// Now spectool knows only nanometers and electron-volts
func isAnyKnownSuffix(s string) bool {
	knownSuffixes := []string{".nm", ".ev"}
	for _, sfx := range knownSuffixes {
		if s == sfx {
			return true
		}
	}
	return false
}

// Adds pre-suffix or extension to filename or changes the existing one
// addPreSuffix("data1.txt", "ev") == data1.ev.txt
// addPreSuffix("data1.nm", "ev") == data1.ev
// addPreSuffix("data1.nm.txt", "ev") == data1.ev.txt
func addPreSuffix(fname, newSfx string) string {
	newSfx = OPERATION_DELIMITER + newSfx
	ext := filepath.Ext(fname) // Extension with dot
	if strings.IndexAny(ext, VALUE_BRACE_CLOSE) >= 0 {
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
		if strings.IndexAny(sfx, VALUE_BRACE_CLOSE) >= 0 {
			return subFname + newSfx + ext
		}
		return strings.TrimSuffix(subFname, sfx) + newSfx + ext
	}
	return subFname + newSfx + ext
}

// Adds per-pre-suffix to filename related to any operation made on the file
// addPrePreSuffix("data1.txt", "div(4.54)") == data1.div(4.54).txt
// addPrePreSuffix("data1.nm", "div(4.54)") == data1.div(4.54).nm
// addPrePreSuffix("data1.add(3).nm.txt", "div(4.54)") == data1.add(3).div(4.54).nm.txt
func AddPrePreSuffix(fname, newSfx string) string {
	newSfx = "." + newSfx
	ext := filepath.Ext(fname)
	if strings.IndexAny(ext, VALUE_BRACE_CLOSE) >= 0 {
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

// Add info about operation involving a spectrum in a file name
func (sw *SpectrumWrapper) AddSpOpSuffix(op, fname string) {
	sfx := op + VALUE_BRACE_OPEN + fname + VALUE_BRACE_CLOSE
	sw.fname = AddPrePreSuffix(sw.fname, sfx)
}

// Add info about operation involving a number in a file name
func (sw *SpectrumWrapper) AddNumOpSuffix(op string, num float64) {
	sfx := fmt.Sprintf("%s%s%v%s", op, VALUE_BRACE_OPEN, num, VALUE_BRACE_CLOSE)
	sw.fname = AddPrePreSuffix(sw.fname, sfx)
}

// Add info about operation involving a number in a file name
func (sw *SpectrumWrapper) AddOpSuffix(op, s string) {
	sfx := fmt.Sprintf("%s%s%s%s", op, VALUE_BRACE_OPEN, s, VALUE_BRACE_CLOSE)
	sw.fname = AddPrePreSuffix(sw.fname, sfx)
}
