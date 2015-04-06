// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code in any way you wish.
//
// Evgeny Shevchenko
// shvgn@protonmail.ch
// 2015

package main

import (
	"path/filepath"
	"strings"
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
	newSfx = "." + newSfx
	ext := filepath.Ext(fname) // Extension with dot
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
		return strings.TrimSuffix(subFname, sfx) + newSfx + ext
	}
	return subFname + newSfx + ext
}
