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
	"os"
	"path/filepath"
)

const (
	Version  string = "1.0"
	Author   string = "Eugene Shevchenko"
	Email    string = "shvgn@protonmail.ch"
	URL      string = "https://github.com/shvgn/spectool"
	Liscence string = "MIT"
)

// Global for the verbosity control
var verboseFlag bool

// Message on an arithmetic operation
func opMessage(op, value string) {
	if verboseFlag {
		fmt.Printf("  %s   %v\n", op, value)
	}
}

func main() {

	// Modificating flags
	// X Units
	nm2EvFlag := flag.Bool("2ev", false, "Keep X in electron-volts")
	ev2NmFlag := flag.Bool("2nm", false, "Keep X in nanometers")

	// X cutting options
	fromFlag := flag.Float64("xfrom", -1.0, "X to start from")
	toFlag := flag.Float64("xto", -1.0, "X to end with")

	// Spectra arithmetic operations with numbers
	addNumFlag := flag.Float64("nadd", 0.0, "Add a number ")
	subNumFlag := flag.Float64("nsub", 0.0, "Subtract a number ")
	mulNumFlag := flag.Float64("nmul", 1.0, "Multiply by a number ")
	divNumFlag := flag.Float64("ndiv", 1.0, "Divide by a number ")

	// Spectra operations with other spectra
	addFlag := flag.String("spadd", "", "Add spectrum")
	subFlag := flag.String("spsub", "", "Subtract spectrum")
	mulFlag := flag.String("spmul", "", "Multiply by spectrum")
	divFlag := flag.String("spdiv", "", "Divide by spectrum")

	// Spectra metadata
	noiseFlag := flag.Bool("n", false, "Subtract noise")

	meanFlag := flag.Bool("mean", false, "(Not implemented) Mean spectrum from all the passed data")
	smoothFlag := flag.String("smooth", "",
		"[ws,po]\t(Not implemented) Smooth data with optionally specified both window size and polynome order")
	// pleFlag := flag.String("ple", "", "This is set of wavelength or energy walues: -ple=287.5,288,288.5")

	// Non-modificating flags
	statsFlag := flag.Bool("s", false, "(Not implemented) Collect statistics on the data")
	colXFlag := flag.Int("colx", 1, "Set number of the X column in passed data files")
	colYFlag := flag.Int("coly", 2, "Set number of the Y column in passed data ASCII files")
	// colsFlag := flag.String("cols", "1,2", "Set numbers of X and Y columns")
	// inFmtFlag := flag.String("if", "ascii", "ascii|tsv|csv\tFormat of the input file")
	outFmtFlag := flag.String("of", "ascii", "[ascii|tsv|csv]   Format of the output file")
	outDirFlag := flag.String("od", "", "Directory for output files. If not specified new files are placed near original ones")
	// verboseFlag := flag.Bool("v", false, "Verbose the actions")
	flag.BoolVar(&verboseFlag, "v", false, "Verbose the actions")

	flag.Parse()

	// Parsing filenames from passed strings. Those are considered to be files
	// and globs in order to work in both Windows cmd and Unix shells
	var spData []*SpectrumWrapper
	var sw *SpectrumWrapper
	var err error
	var files []string

	for _, arg := range flag.Args() {

		if sw, err = NewSpecWrapper(arg, *colXFlag, *colYFlag); err != nil {
			if verboseFlag {
				fmt.Println("Cannot open file", arg, ":",
					err, "Trying woth glob...")
			}
			if files, err = filepath.Glob(arg); err != nil {
				if verboseFlag {
					fmt.Println("Nor filename nor correct glob pattern.")
				}
				fmt.Println(err)
				continue
			}
			for _, f := range files { // arg is a valid glob pattern
				sw, err = NewSpecWrapper(f, *colXFlag, *colYFlag)
				if err != nil {
					fmt.Println(f+": skipped file because of parse error:", err)
					continue
				}
				spData = append(spData, sw)
			}
			continue // Appended files from the glob
		}
		spData = append(spData, sw) // arg is a valid filename
	}

	// Choosing units for the processing
	// Forbid using -2ev and -2nm together. FIXME Why?
	if *nm2EvFlag && *ev2NmFlag {
		log.Fatal("Cannot work on nanometers and electron-volts simultaneously. Sorry.")
	}
	modifyUnits := *nm2EvFlag || *ev2NmFlag
	var ensureUnitsFunc func(float64) float64
	var unitsPreSuffix string

	if *nm2EvFlag {
		ensureUnitsFunc = ensureEv
		unitsPreSuffix = "ev"
	} else if *ev2NmFlag {
		ensureUnitsFunc = ensureNm
		unitsPreSuffix = "nm"
	} else {
		ensureUnitsFunc = func(x float64) float64 {
			log.Fatal("Unexpected units conversion")
			return 0.0
		}
	}

	// X from and to
	if *fromFlag >= 0 {
		*fromFlag = ensureUnitsFunc(*fromFlag)
	}
	if *toFlag >= 0 {
		*toFlag = ensureUnitsFunc(*toFlag)
	}
	if *fromFlag > *toFlag {
		*fromFlag, *toFlag = *toFlag, *fromFlag
	}

	// Arithmetics operands
	var addSpectrum, subSpectrum, mulSpectrum, divSpectrum *SpectrumWrapper
	// var err error
	if *addFlag != "" {
		if addSpectrum, err = NewSpecWrapper(*addFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			addSpectrum.s.ModifyX(ensureUnitsFunc)
		}
	}
	if *subFlag != "" {
		if subSpectrum, err = NewSpecWrapper(*subFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			subSpectrum.s.ModifyX(ensureUnitsFunc)
		}
	}
	if *mulFlag != "" {
		if mulSpectrum, err = NewSpecWrapper(*mulFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			mulSpectrum.s.ModifyX(ensureUnitsFunc)
		}
	}
	if *divFlag != "" {
		if divSpectrum, err = NewSpecWrapper(*divFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			divSpectrum.s.ModifyX(ensureUnitsFunc)
		}
	}

	// Processing
	for _, sw := range spData {
		if verboseFlag {
			fmt.Println()
			fmt.Println(sw.dir + sw.fname)
		}

		// Subtract the noise from the full-length signal
		if *noiseFlag {
			n := sw.s.Noise()
			opMessage("-", fmt.Sprintf("%v (noise)", n))
			sw.s.ModifyY(func(y float64) float64 { return y - n })
			sw.AddNumOpSuffix("noise", n)
		}

		// Process the X units
		if modifyUnits {
			sw.s.ModifyX(ensureUnitsFunc)
			sw.fname = addPreSuffix(sw.fname, unitsPreSuffix)

			// Real spectrum X is always assumed to be positive
			// FIXME Make cut one-sided
			xl, xr := *fromFlag, *toFlag
			if xl > 0 && xr > 0 {
				opMessage(">", fmt.Sprintf("%v", xl))
				opMessage("<", fmt.Sprintf("%v", xr))
				sw.s.Cut(xl, xr)
			}
		}
		// Addition and subtracting of spectra should be done before noise calculation
		if *addFlag != "" {
			sw.s.Add(addSpectrum.s)
			opMessage("+", addSpectrum.fname)
			sw.AddSpOpSuffix("add", addSpectrum.fname)
		}
		if *subFlag != "" {
			sw.s.Subtract(subSpectrum.s)
			opMessage("-", subSpectrum.fname)
			sw.AddSpOpSuffix("sub", subSpectrum.fname)
		}

		if *mulFlag != "" {
			opMessage("*", mulSpectrum.fname)
			sw.s.Multiply(mulSpectrum.s)
			sw.AddSpOpSuffix("mul", mulSpectrum.fname)
		}
		if *divFlag != "" {
			sw.s.Divide(divSpectrum.s)
			opMessage("/", divSpectrum.fname)
			sw.AddSpOpSuffix("div", divSpectrum.fname)
		}

		// Arithmetics with numbers
		if *addNumFlag != 0.0 {
			opMessage("+", fmt.Sprintf("%v", *addNumFlag))
			sw.s.ModifyY(func(y float64) float64 { return y + *addNumFlag })
			sw.AddNumOpSuffix("add", *addNumFlag)
		}
		if *subNumFlag != 0.0 {
			opMessage("-", fmt.Sprintf("%v", *subNumFlag))
			sw.s.ModifyY(func(y float64) float64 { return y - *subNumFlag })
			sw.AddNumOpSuffix("sub", *subNumFlag)
		}
		if *mulNumFlag != 1.0 {
			opMessage("*", fmt.Sprintf("%v", *mulNumFlag))
			sw.s.ModifyY(func(y float64) float64 { return y * *mulNumFlag })
			sw.AddNumOpSuffix("mul", *mulNumFlag)
		}
		if *divNumFlag != 1.0 {
			opMessage("/", fmt.Sprintf("%v", *divNumFlag))
			sw.s.ModifyY(func(y float64) float64 { return y / *divNumFlag })
			sw.AddNumOpSuffix("div", *divNumFlag)
		}

		if *smoothFlag != "" {
			// SMOOTH THEM ALL!!!1
		}
		if *meanFlag {
			// MEAN THEM ALL
			sw.AddNumOpSuffix("mean", float64(len(spData)))
		}
		if *statsFlag {
			// Calculate stats
			// fmt.Println(stats(sw.s))
		}
	}

	// Saving
	// Directory to save in
	if *outDirFlag != "" {
		var perm os.FileMode = 0755 // FIXME Why use something else?
		err := os.MkdirAll(*outDirFlag, perm)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Output format
	if *outFmtFlag == "" {
		*outFmtFlag = "ascii"
	}

	for _, sw := range spData {
		var path string
		var perm os.FileMode = 0644 // FIXME Why use something else?

		if *outDirFlag != "" {
			path = filepath.Join(*outDirFlag, sw.fname)
		} else {
			path = filepath.Join(sw.dir, sw.fname)
		}

		err := sw.WriteFile(path, *outFmtFlag, perm)
		if err != nil {
			log.Fatal(err)
		}
	}

	// --------------------------------------------------------------------------
	// Take PLE values into account TODO

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
