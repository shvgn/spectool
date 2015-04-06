// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code in any way you wish.
//
// Evgeny Shevchenko
// shvgn@protonmail.ch
// 2015

package main

import "testing"

func TestAddPreSuffix(t *testing.T) {
	names := []struct {
		fname string
		sfx   string
		want  string
	}{
		{"data1.txt", "ev", "data1.ev.txt"},
		{"data1.dat", "nm", "data1.nm.dat"},
		{"data1.nm", "ev", "data1.ev"},
		{"data1.ev", "nm", "data1.nm"},
		{"data1.nm", "nm", "data1.nm"},
		{"data1.ev", "ev", "data1.ev"},
		{"d.nm.a.ev.t.a.ev.dat", "ev", "d.nm.a.ev.t.a.ev.dat"},
		{"d.nm.a.ev.t.a.ev.dat", "nm", "d.nm.a.ev.t.a.nm.dat"},
		{"d.nm.a.ev.t.a.ev", "ev", "d.nm.a.ev.t.a.ev"},
		{"d.nm.a.ev.t.a.ev", "nm", "d.nm.a.ev.t.a.nm"},
	}
	for _, s := range names {
		got := addPreSuffix(s.fname, s.sfx)
		if got != s.want {
			t.Errorf("addPreSuffix(%q, %q) == %q, want %q",
				s.fname, s.sfx, got, s.want)
		}

	}
}
