package main

import (
	"fmt"
	"testing"
)

type TestCase struct {
	fname string
	sfx   string
	want  string
}

func TestAddPreSuffix(t *testing.T) {

	cases := []TestCase{
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
		{"d.nm.a.ev.t.a.div(3.14).ev.dat", "ev", "d.nm.a.ev.t.a.div(3.14).ev.dat"},
		{"d.nm.a.ev.t.a.div(3.14).ev.dat", "nm", "d.nm.a.ev.t.a.div(3.14).nm.dat"},
		{"d.nm.a.ev.t.a.div(3.14).ev", "ev", "d.nm.a.ev.t.a.div(3.14).ev"},
		{"d.nm.a.ev.t.a.div(3.14).ev", "nm", "d.nm.a.ev.t.a.div(3.14).nm"},
	}

	// Fill with brackets
	for _, sfx := range knownSuffixes() {
		// Without .dat
		c1 := TestCase{
			fname: fmt.Sprintf("d.nm.a.ev.t.a.div%s3.14%s.ev",
				valueBraceOpen, valueBraceClose),
			sfx: sfx,
			want: fmt.Sprintf("d.nm.a.ev.t.a.div%s3.14%s.%s",
				valueBraceOpen, valueBraceClose, sfx),
		}
		cases = append(cases, c1)

		// With .dat
		c2 := TestCase{
			fname: fmt.Sprintf("d.nm.a.ev.t.a.div%s3.14%s.ev.dat",
				valueBraceOpen, valueBraceClose),
			sfx: sfx,
			want: fmt.Sprintf("d.nm.a.ev.t.a.div%s3.14%s.%s.dat",
				valueBraceOpen, valueBraceClose, sfx),
		}
		cases = append(cases, c2)
	}

	for _, s := range cases {
		got := addPreSuffix(s.fname, s.sfx)
		if got != s.want {
			t.Errorf("addPreSuffix(%q, %q) -> %q, want %q",
				s.fname, s.sfx, got, s.want)
		}
	}
}

func TestAddPrePreSuffix(t *testing.T) {
	cases := []TestCase{
		{"data1", "sfx", "data1.sfx"},
		{"data1.txt", "sfx", "data1.sfx.txt"},
		{"data1.ev.dat", "sfx", "data1.sfx.ev.dat"},
		{"data1.nm", "sfx", "data1.sfx.nm"},
		{"d.nm.a.ev.t.a.ev.dat", "sfx", "d.nm.a.ev.t.a.sfx.ev.dat"},
	}

	// Fill with brackets
	sfx := "sfx"

	// Without .dat
	c1 := TestCase{
		fname: fmt.Sprintf("d.nm.a.ev.t.a.div%s3.14%s.ev",
			valueBraceOpen, valueBraceClose),
		sfx: sfx,
		want: fmt.Sprintf("d.nm.a.ev.t.a.div%s3.14%s.%s.ev",
			valueBraceOpen, valueBraceClose, sfx),
	}
	cases = append(cases, c1)

	// With .dat
	c2 := TestCase{
		fname: fmt.Sprintf("d.nm.a.ev.t.a.div%s3.14%s.ev.dat",
			valueBraceOpen, valueBraceClose),
		sfx: sfx,
		want: fmt.Sprintf("d.nm.a.ev.t.a.div%s3.14%s.%s.ev.dat",
			valueBraceOpen, valueBraceClose, sfx),
	}
	cases = append(cases, c2)

	for _, s := range cases {
		got := addPrePreSuffix(s.fname, s.sfx)
		if got != s.want {
			t.Errorf("addPrePreSuffix(%q, %q) -> %q, want %q",
				s.fname, s.sfx, got, s.want)
		}
	}
}
