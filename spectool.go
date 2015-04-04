// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code in any way you wish.
//
// Evgeny Shevchenko
// shvgn@protonmail.ch
// 2015

package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {

	nm2EvFlag := flag.Bool("nm2ev", false, "Convert X from nanometers to electron-volts")
	ev2NmFlag := flag.Bool("ev2nm", false, "Convert X from electron-volts to nanometers")
	// inFmtFlag := flag.String("if", "ascii", "ascii|tsv|csv\tFormat of the input file")
	// outFmtFlag := flag.String("of", "ascii", "ascii|tsv|csv\tFormat of the output file")
	// pleFlag := flag.String("ple", "", "This is set of wavelength or energy walues: -ple=287.5,288,288.5")
	flag.Parse()

	// var parseSpecFunc func(io.Reader, int, int) (*spectrum.Spectrum, error)
	// switch *inFmtFlag {
	// case "tsv":
	// 	parseSpecFunc = spectrum.ReadFromTSV
	// case "csv":
	// 	parseSpecFunc = nil
	// case "ascii":
	// 	parseSpecFunc = nil
	// }
	//
	// var saveSpecFunc func(io.Writer) error
	// switch *outFmtFlag {
	// case "tsv":
	// 	saveSpecFunc = nil
	// case "csv":
	// 	saveSpecFunc = nil
	// case "ascii":
	// 	saveSpecFunc = nil
	// }

	// Filling the data
	spectra := make([]*SpectrumWrapper, len(flag.Args()))
	for _, fname := range flag.Args() {
		sw, err := NewSpecWrapper(fname)
		if err != nil {
			log.Println("Cannot read spectrum from file", fname+":", err.Error(), "- Skipping.")
			continue
		}
		fmt.Println(sw)
		spectra = append(spectra, sw)
	}

	// Processing
	for _, sw := range spectra {
		if *nm2EvFlag {
			sw.s.ModifyX(ensureEv)
			sw.fname = addPreSuffix("ev")
		} else if *ev2NmFlag {
			sw.s.ModifyX(ensureNm)
			sw.fname = addPreSuffix("nm")
		}

	}

	// --------------------------------------------------------------------------
	// Take PLE values into account

	// pleVals := make([]float64, 0)

	// if *pleSet != "" {
	// 	if *pleSet == "all" {
	// 		// Make the all!
	// 	}

	// 	// If it contains colons ':' than interpret is as matlab/julia range
	// 	// e.g. 345:0.2:330 is the same as 330:0.2:345
	// 	// 345:330 equals 330:1:345
	// 	// aaaand we must also take into account electron-volts 4.5:0.05:5

	// 	pleValStr := strings.Split(*pleSet, ",")
	// 	fmt.Println("pleValStr: ", pleValStr)

	// 	for _, v := range pleValStr {
	// 		pleval, err := strconv.ParseFloat(v, 64)
	// 		if err != nil {
	// 			fmt.Println("Warning! Could not parse value for PLE:", v, "- Skipping.")
	// 			continue
	// 		}
	// 		pleVals = append(pleVals, pleval)
	// 	}
	// }

	// --------------------------------------------------------------------------

	// fmt.Println("Number of values to average: ", *averPtr)
	// fmt.Println("PLE detection values passed: ", *pleSet)
	fmt.Println("nm to eV: ", *nm2EvFlag)
	fmt.Println("eV to nm: ", *ev2NmFlag)
	// fmt.Println("PLE detection values parsed: ", pleVals)
	fmt.Println("Other cmd arguments: ", flag.Args())
}

/************************************************************************

Interface

Tasks:

A specrum file is a two-column ASCII file with numbers, the columns being
separated by space characters such as multiple whitespaces or tabs (TSV file).
Headers are allowed. If a header has colon ':', the colon is considered to be
the delimeter, otherwise it will be first space character met after the first
word. Comments must start with #.


HeaderName: Header Value with a bunch of whitespaces
Header Name 2: And this must work, too
Header Value is now from the second word and towards the end
1.000000 4.3123353
1.010000    12434,53432

...


1. specify operation with the special key?

	spectool -op=nm2ev file1 file2 file3 ...

2. seaprate key for each operation

	spectool -nm2ev file1 file2 file3 ...
	spectool -ple file1 file2 file3 ...
	spectool -mean file1 file2 file3 ...
	spectool -stats file1 file2 file3 ...

	spectool -stats -ple 360 360.5 362 [...] -nm2ev file1 file2 file3 ...
		or spectool -stats -ple=all -nm2ev file1 file2 file3 ...

	spectool -nm2ev file1 fil32 file3 -noise -stats -from=235 -to=310 -xshift=1.23 -yshift=40 -mul -div -add -sub


*************************************************************************/
