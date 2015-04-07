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
)

func main() {

	// Modificating flags
	nm2EvFlag := flag.Bool("ev", false, "Convert X from nanometers to electron-volts")
	ev2NmFlag := flag.Bool("nm", false, "Convert X from electron-volts to nanometers")
	addFlag := flag.String("add", "", "Add a number or a spectrum")
	subFlag := flag.String("sub", "", "Subtract a number or a spectrum")
	mulFlag := flag.String("mul", "", "Multiply by a number or a spectrum")
	divFlag := flag.String("div", "", "Divide by a number or a spectrum")
	meanFlag := flag.Bool("mean", false, "Mean spectrum from all the passed data")
	smoothFlag := flag.String("smooth", "no",
		"[ws,po]\tSmooth data with optionally specified both window size and polynome order")

	// Not modificating flags
	statsFlag := flag.Bool("s", false, "Collect statistics on the data")
	colXFlag := flag.Int("xcol", 1, "Set number of the X column")
	colYFlag := flag.Int("ycol", 2, "Set number of the Y column")
	// colsFlag := flag.String("cols", "1,2", "Set numbers of X and Y columns")
	// inFmtFlag := flag.String("if", "ascii", "ascii|tsv|csv\tFormat of the input file")
	// outFmtFlag := flag.String("of", "ascii", "ascii|tsv|csv\tFormat of the output file")
	// pleFlag := flag.String("ple", "", "This is set of wavelength or energy walues: -ple=287.5,288,288.5")

	flag.Parse()

	var modificationsRequired = false
	if *nm2EvFlag || *ev2NmFlag ||
		*addFlag != "" || *subFlag != "" ||
		*mulFlag != "" || *divFlag != "" {
		modificationsRequired = true
	}
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
	originals := make([]*SpectrumWrapper, len(flag.Args()))
	var modified []*SpectrumWrapper
	for _, fname := range flag.Args() {
		sw, err := NewSpecWrapper(fname, *colXFlag, *colYFlag)
		if err != nil {
			fmt.Print("Error processing file " + fname + ": ")
			fmt.Println(err)
			continue
		}
		originals = append(originals, sw)
		if modificationsRequired {
			modified = append(modified, sw)
		}
	}

	// Processing
	for _, sw := range modified {
		if *nm2EvFlag {
			sw.s.ModifyX(ensureEv)
			sw.fname = addPreSuffix(sw.fname, "ev")
		} else if *ev2NmFlag {
			sw.s.ModifyX(ensureNm)
			sw.fname = addPreSuffix(sw.fname, "nm")
		}
		if *smoothFlag != "no" {
			// SMOOTH THEM ALL!!!1
		}
		if *meanFlag {
			// MEAN THEM ALL
		}
		if *statsFlag {
			// Calculate stats
			fmt.Println(stats(sw.s))
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
	// fmt.Println("nm to eV: ", *nm2EvFlag)
	// fmt.Println("eV to nm: ", *ev2NmFlag)
	// fmt.Println("PLE detection values parsed: ", pleVals)
	// fmt.Println("Other cmd arguments: ", flag.Args())
}

/************************************************************************

Interface

Tasks:

A specrum file is a two-column ASCII file with numbers, the columns being
separated by space characters such as multiple whitespaces or tabs (TSV file).
Headers are allowed. For an arbitrary ASCII file if a header has colon ':', the
colon is considered to be the delimeter, otherwise it will be first space
character met after the first word.


HeaderName: Header Value with a bunch of whitespaces
Header Name 2: And this must work, too
Header Value is now from the second word and towards the end
1.000000 4.3123353
1.010000    12434,53432

...

	spectool -nm2ev file1 file2 file3 ...
	spectool -ple file1 file2 file3 ...
	spectool -mean file1 file2 file3 ...
	spectool -stats file1 file2 file3 ...

	spectool -stats -ple 360 360.5 362 [...] -nm2ev file1 file2 file3 ...
		or spectool -stats -ple=all -nm2ev file1 file2 file3 ...

	spectool -nm2ev file1 fil32 file3 -noise -stats -from=235 -to=310 -xshift=1.23 -yshift=40 -mul -div -add -sub


*************************************************************************/
