// Simple command line utility for the manipulation of columned ASCII data files
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

func TestAddPrePreSuffix(t *testing.T) {
	names := []struct {
		fname string
		sfx   string
		want  string
	}{
		{"data1", "sfx", "data1.sfx"},
		{"data1.txt", "sfx", "data1.sfx.txt"},
		{"data1.ev.dat", "sfx", "data1.sfx.ev.dat"},
		{"data1.nm", "sfx", "data1.sfx.nm"},
		{"d.nm.a.ev.t.a.ev.dat", "sfx", "d.nm.a.ev.t.a.sfx.ev.dat"},
	}
	for _, s := range names {
		got := addPreSuffix(s.fname, s.sfx)
		if got != s.want {
			t.Errorf("addPreSuffix(%q, %q) == %q, want %q",
				s.fname, s.sfx, got, s.want)
		}

	}
}
