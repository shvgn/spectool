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

	"github.com/shvgn/spectrum"
)

func main() {

	// Modificating flags
	// X Units
	nm2EvFlag := flag.Bool("ev", false, "Keep X in electron-volts")
	ev2NmFlag := flag.Bool("nm", false, "Keep X in nanometers")

	// X cutting options
	fromFlag := flag.Float64("from", -1.0, "X to start from")
	toFlag := flag.Float64("to", -1.0, "X to end with")

	// Spectra arithmetic operations with numbers
	addNumFlag := flag.Float64("addn", 0.0, "Add a number ")
	subNumFlag := flag.Float64("subn", 0.0, "Subtract a number ")
	mulNumFlag := flag.Float64("muln", 1.0, "Multiply by a number ")
	divNumFlag := flag.Float64("divn", 1.0, "Divide by a number ")

	// Spectra operations with other spectra
	addFlag := flag.String("add", "", "Add spectrum")
	subFlag := flag.String("sub", "", "Subtract spectrum")
	mulFlag := flag.String("mul", "", "Multiply by spectrum")
	divFlag := flag.String("div", "", "Divide by spectrum")

	// Spectra metadata
	noiseFlag := flag.Bool("noise", false, "(Not implemented) Subtract noise")

	meanFlag := flag.Bool("mean", false, "(Not implemented) Mean spectrum from all the passed data")
	smoothFlag := flag.String("smooth", "",
		"[ws,po]\t(Not implemented) Smooth data with optionally specified both window size and polynome order")

	// Non-modificating flags
	statsFlag := flag.Bool("s", false, "(Not implemented) Collect statistics on the data")
	colXFlag := flag.Int("xcol", 1, "Set number of the X column")
	colYFlag := flag.Int("ycol", 2, "Set number of the Y column")
	// colsFlag := flag.String("cols", "1,2", "Set numbers of X and Y columns")
	// inFmtFlag := flag.String("if", "ascii", "ascii|tsv|csv\tFormat of the input file")
	// outFmtFlag := flag.String("of", "ascii", "ascii|tsv|csv\tFormat of the output file")
	// pleFlag := flag.String("ple", "", "This is set of wavelength or energy walues: -ple=287.5,288,288.5")

	flag.Parse()

	var modificationsRequired bool = *nm2EvFlag || *ev2NmFlag ||
		*addFlag != "" || *subFlag != "" || *mulFlag != "" || *divFlag != ""

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

	// Choosing units for the processing
	// Forbid using -ev and -nm together
	if *nm2EvFlag && *ev2NmFlag {
		log.Fatal("Cannot work on nanometers and electron-volts simultaneously. Sorry.")
	}
	modifyUnits := *nm2EvFlag || *ev2NmFlag
	var ensureUnitsFunc func(float64) float64
	var unitsPreSuffix string

	switch modifyUnits {
	case *nm2EvFlag:
		ensureUnitsFunc = ensureEv
		unitsPreSuffix = "ev"
	case *ev2NmFlag:
		ensureUnitsFunc = ensureNm
		unitsPreSuffix = "nm"
	}

	if *fromFlag >= 0 {
		*fromFlag = ensureUnitsFunc(*fromFlag)
	}
	if *toFlag >= 0 {
		*toFlag = ensureUnitsFunc(*toFlag)
	}

	// Arithmetics to be done
	var addSpectrum, subSpectrum, mulSpectrum, divSpectrum *spectrum.Spectrum
	var err error
	if *addFlag != "" {
		if addSpectrum, err = spectrum.SpectrumFromFile(*addFlag); err != nil {
			log.Fatal(err)
		}
	}
	if *subFlag != "" {
		if subSpectrum, err = spectrum.SpectrumFromFile(*subFlag); err != nil {
			log.Fatal(err)
		}
	}
	if *mulFlag != "" {
		if mulSpectrum, err = spectrum.SpectrumFromFile(*mulFlag); err != nil {
			log.Fatal(err)
		}
	}
	if *divFlag != "" {
		if divSpectrum, err = spectrum.SpectrumFromFile(*divFlag); err != nil {
			log.Fatal(err)
		}
	}

	// Processing
	for _, sw := range modified {
		// Addition and subtracting of spectra should be done before noise calculation
		if *addFlag != "" {
			sw.s.Add(addSpectrum)
		}
		if *subFlag != "" {
			sw.s.Subtract(subSpectrum)
		}

		// Subtract the noise from the full-length signal
		if *noiseFlag {
			n := sw.s.Noise()
			sw.s.ModifyY(func(y float64) float64 { return y - n })
		}
		// Process the X units
		if modifyUnits {
			sw.s.ModifyX(ensureUnitsFunc)
			sw.fname = addPreSuffix(sw.fname, unitsPreSuffix)

			// Real spectrum X is always assumed to be positive
			xl, xr := *fromFlag, *toFlag
			if xl > 0 && xr > 0 {
				sw.s.Cut(xl, xr)
			}
		}

		if *mulFlag != "" {
			sw.s.Multiply(mulSpectrum)
		}
		if *divFlag != "" {
			sw.s.Divide(divSpectrum)
		}

		// Arithmetics with numbers
		if *addNumFlag != 0.0 {
			sw.s.ModifyY(func(y float64) float64 { return y + *addNumFlag })
		}
		if *subNumFlag != 0.0 {
			sw.s.ModifyY(func(y float64) float64 { return y - *subNumFlag })
		}
		if *mulNumFlag != 1.0 {
			sw.s.ModifyY(func(y float64) float64 { return y * *mulNumFlag })
		}
		if *divNumFlag != 1.0 {
			sw.s.ModifyY(func(y float64) float64 { return y / *divNumFlag })
		}

		if *smoothFlag != "" {
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

	// Saving
	for _, sw := range modified {
		sw.Write()
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
