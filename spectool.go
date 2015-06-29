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
	"math"
	"math/rand"
	"os"
	"path/filepath"
)

var (
	verboseFlag bool // Verbosity of the output

	keepEvFlag bool // Keep units in electron-volts
	keepNmFlag bool // Keep units in nanometers

	fromFlag float64 // Cut X from value
	toFlag   float64 // Cut X up to value

	addNumFlag float64 // Add number
	subNumFlag float64 // Subtract number
	mulNumFlag float64 // Multiply by number
	divNumFlag float64 // Divide by number

	addFlag string // Add spectrum
	subFlag string // Subtract spectrum
	mulFlag string // Multiply by spectrum
	divFlag string // Divide by spectrum

	noiseFlag bool // Calculate and subtract noise
	// meanFlag   bool   // Calculate arithmetic mean of spectra
	// smoothFlag string // Smooth spectra
	// statsFlag  bool   // Calculate metainfo

	colXFlag int // Column in ASCII file to read X from
	colYFlag int // Column in ASCII file to read Y from

	outFmtFlag string // Ouput format
	outDirFlag string // Ouput directory for resulting spectra
)

// Message on an arithmetic operation
func opMessage(op, value string) {
	if verboseFlag {
		fmt.Printf("  %s   %v\n", op, value)
	}
}

func init() {
	// Modificating flags
	// X Units
	flag.BoolVar(&keepEvFlag, "2ev", false, "Keep X in electron-volts")
	flag.BoolVar(&keepNmFlag, "2nm", false, "Keep X in nanometers")

	// X cutting options
	flag.Float64Var(&fromFlag, "xfrom", math.Inf(-1), "X to start from")
	flag.Float64Var(&toFlag, "xto", math.Inf(1), "X to end with")

	// Spectra arithmetic operations with numbers
	flag.Float64Var(&addNumFlag, "nadd", 0.0, "Add a number ")
	flag.Float64Var(&subNumFlag, "nsub", 0.0, "Subtract a number ")
	flag.Float64Var(&mulNumFlag, "nmul", 1.0, "Multiply by a number ")
	flag.Float64Var(&divNumFlag, "ndiv", 1.0, "Divide by a number ")

	// Spectra operations with other spectra
	flag.StringVar(&addFlag, "spadd", "", "Add spectrum")
	flag.StringVar(&subFlag, "spsub", "", "Subtract spectrum")
	flag.StringVar(&mulFlag, "spmul", "", "Multiply by spectrum")
	flag.StringVar(&divFlag, "spdiv", "", "Divide by spectrum")

	// Spectra metadata
	flag.BoolVar(&noiseFlag, "n", false, "Subtract noise")

	// flag.BoolVar(&meanFlag, "mean", false, "(Not implemented) Mean spectrum from all the passed data")
	// flag.StringVar(&smoothFlag, "smooth", "",
	// 	"[ws,po]\t(Not implemented) Smooth data with optionally specified both window size and polynome order")

	// Non-modificating flags
	// flag.BoolVar(&statsFlag, "s", false, "(Not implemented) Collect statistics on the data")
	flag.IntVar(&colXFlag, "colx", 1, "Set number of the X column in passed data files")
	flag.IntVar(&colYFlag, "coly", 2, "Set number of the Y column in passed data ASCII files")

	flag.StringVar(&outFmtFlag, "of", "ascii", "[ascii|tsv|csv]   Format of the output file")
	flag.StringVar(&outDirFlag, "od", "", "Directory for output files")

	flag.BoolVar(&verboseFlag, "v", false, "Verbose the actions")

	flag.Parse()
}

func main() {

	// Parsing filenames from passed strings. Those are considered to be files
	// and globs in order to work in both Windows cmd and Unix shells
	var spData []*SpectrumWrapper
	var sw *SpectrumWrapper
	var err error
	var files []string

	for _, arg := range flag.Args() { // Remaining arguments are filepaths or globs to process
		if sw, err = NewSpecWrapper(arg, colXFlag, colYFlag); err != nil {
			if verboseFlag {
				fmt.Println("Cannot open file", arg, ":", err, "Trying with glob...")
			}
			if files, err = filepath.Glob(arg); err != nil {
				if verboseFlag {
					fmt.Println("Nor filename nor correct glob pattern. Skipping.")
				}
				fmt.Println(err)
				continue
			}
			for _, f := range files { // arg is a valid glob pattern
				sw, err = NewSpecWrapper(f, colXFlag, colYFlag)
				if err != nil {
					fmt.Println(f+": Parse error:", err, "Skipping.")
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
	if keepEvFlag && keepNmFlag {
		log.Fatal("Cannot work on multiple X units simultaneously. Sorry.")
	}
	modifyUnits := keepEvFlag || keepNmFlag
	ensureUnitsFunc := func(x float64) float64 {
		if modifyUnits { // Something went wrong
			log.Fatal("Unexpected units conversion")
			return math.NaN()
		}
		return x
	}

	var unitsPreSuffix string
	if keepEvFlag {
		ensureUnitsFunc = ensureEv
		unitsPreSuffix = "ev"
	}
	if keepNmFlag {
		ensureUnitsFunc = ensureNm
		unitsPreSuffix = "nm"
	}

	// Boundaries for cutting X from and to
	if fromFlag > toFlag { // The X cutting flags are invalid
		log.Fatal("Invalid order of X cutting flags")
	}

	var cutLeft, cutRight bool
	var checker1, checker2 float64 // We want not to mess with the order of cutting values

	if !math.IsInf(fromFlag, -1) {
		cutLeft = true
		checker1 = fromFlag
		fromFlag = ensureUnitsFunc(fromFlag)
	}
	if !math.IsInf(toFlag, 1) {
		cutRight = true
		checker2 = toFlag
		toFlag = ensureUnitsFunc(toFlag)
	}

	if modifyUnits && (cutLeft || cutRight) { // Check the order of cut boundaries
		if cutLeft && cutRight {
			if fromFlag > toFlag { // X order was reversed
				fromFlag, toFlag = toFlag, fromFlag
				cutLeft, cutRight = cutRight, cutLeft
			}
		} else { // Artificial checker required
			if cutLeft {
				checker2 = checker1 + (math.Abs(checker1)+1)*rand.Float64() // checker2 > checker1
			}
			if cutRight {
				checker1 = checker2 - (math.Abs(checker2)+1)*rand.Float64() // checker1 < checker2
			}
			newChecker1, newChecker2 := ensureUnitsFunc(checker1), ensureUnitsFunc(checker2)
			if newChecker1 > newChecker2 { // Modified X has reversed order
				fromFlag, toFlag = toFlag, fromFlag
				cutLeft, cutRight = cutRight, cutLeft
			}
		}

	}

	// Arithmetics operands
	var addSpectrum, subSpectrum, mulSpectrum, divSpectrum *SpectrumWrapper
	if addFlag != "" {
		if addSpectrum, err = NewSpecWrapper(addFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			addSpectrum.s.ModifyX(ensureUnitsFunc)
		}
	}
	if subFlag != "" {
		if subSpectrum, err = NewSpecWrapper(subFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			subSpectrum.s.ModifyX(ensureUnitsFunc)
		}
	}
	if mulFlag != "" {
		if mulSpectrum, err = NewSpecWrapper(mulFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			mulSpectrum.s.ModifyX(ensureUnitsFunc)
		}
	}
	if divFlag != "" {
		if divSpectrum, err = NewSpecWrapper(divFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			divSpectrum.s.ModifyX(ensureUnitsFunc)
		}
	}

	// Processing
	l := len(spData)
	for i, sw := range spData {
		if verboseFlag {
			fmt.Println(fmt.Sprintf("%d/%d  ", i, l) + sw.dir + sw.fname)
		}

		// Subtract the noise from the full-length signal
		if noiseFlag {
			n := sw.s.Noise()
			opMessage("-", fmt.Sprintf("%s (noise)", humanize.Ftoa(n)))
			sw.s.ModifyY(func(y float64) float64 { return y - n })
			sw.AddOpSuffix("noise", humanize.Ftoa(n))
		}

		// Process the X units
		if modifyUnits {
			sw.s.ModifyX(ensureUnitsFunc)
			sw.fname = addPreSuffix(sw.fname, unitsPreSuffix)
		}

		// Cutting is done within boundaries and with spectra after making
		// sure all X units are the same.
		// FIXME Make cut one-sided
		if cutLeft && cutRight {
			opMessage(">", humanize.Ftoa(fromFlag))
			opMessage("<", humanize.Ftoa(toFlag))
			sw.s.Cut(fromFlag, toFlag)
			sw.AddOpSuffix("x", fmt.Sprintf("%s-%s", humanize.Ftoa(fromFlag), humanize.Ftoa(toFlag)))
			// sw.AddOpSuffix("from", humanize.Ftoa(fromFlag))
			// sw.AddOpSuffix("to", humanize.Ftoa(toFlag))
		} else {
			if cutLeft {
				opMessage(">", fmt.Sprintf("%v", humanize.Ftoa(fromFlag)))
				xLast, _ := sw.s.LastPoint()
				sw.s.Cut(fromFlag, xLast)
				sw.AddOpSuffix("x", fmt.Sprintf("%s-", humanize.Ftoa(fromFlag)))
				// sw.AddOpSuffix("from", humanize.Ftoa(fromFlag))
			}
			if cutRight {
				opMessage("<", fmt.Sprintf("%v", humanize.Ftoa(toFlag)))
				xFirst, _ := sw.s.FirstPoint()
				sw.s.Cut(xFirst, toFlag)
				sw.AddOpSuffix("x", fmt.Sprintf("-%s", humanize.Ftoa(toFlag)))
				// sw.AddOpSuffix("to", humanize.Ftoa(toFlag))
			}
		}

		// Addition and subtracting of spectra should be done before noise calculation?
		if addFlag != "" {
			sw.s.Add(addSpectrum.s)
			opMessage("+", addSpectrum.fname)
			sw.AddSpOpSuffix("add", addSpectrum.fname)
		}
		if subFlag != "" {
			sw.s.Subtract(subSpectrum.s)
			opMessage("-", subSpectrum.fname)
			sw.AddSpOpSuffix("sub", subSpectrum.fname)
		}

		if mulFlag != "" {
			opMessage("*", mulSpectrum.fname)
			sw.s.Multiply(mulSpectrum.s)
			sw.AddSpOpSuffix("mul", mulSpectrum.fname)
		}
		if divFlag != "" {
			sw.s.Divide(divSpectrum.s)
			opMessage("/", divSpectrum.fname)
			sw.AddSpOpSuffix("div", divSpectrum.fname)
		}

		// Arithmetics with numbers
		if addNumFlag != 0.0 {
			opMessage("+", fmt.Sprintf("%v", addNumFlag))
			sw.s.ModifyY(func(y float64) float64 { return y + addNumFlag })
			sw.AddNumOpSuffix("add", addNumFlag)
		}
		if subNumFlag != 0.0 {
			opMessage("-", fmt.Sprintf("%v", subNumFlag))
			sw.s.ModifyY(func(y float64) float64 { return y - subNumFlag })
			sw.AddNumOpSuffix("sub", subNumFlag)
		}
		if mulNumFlag != 1.0 {
			opMessage("*", fmt.Sprintf("%v", mulNumFlag))
			sw.s.ModifyY(func(y float64) float64 { return y * mulNumFlag })
			sw.AddNumOpSuffix("mul", mulNumFlag)
		}
		if divNumFlag != 1.0 {
			opMessage("/", fmt.Sprintf("%v", divNumFlag))
			sw.s.ModifyY(func(y float64) float64 { return y / divNumFlag })
			sw.AddNumOpSuffix("div", divNumFlag)
		}

		// if smoothFlag != "" {
		// 	// SMOOTH THEM ALL!!!1
		// }
		// if meanFlag {
		// 	// MEAN THEM ALL
		// 	sw.AddNumOpSuffix("mean", float64(len(spData)))
		// }
		// if statsFlag {
		// 	// Calculate stats
		// 	// fmt.Println(stats(sw.s))
		// }
	}

	// Saving
	// Directory to save in
	if outDirFlag != "" {
		var perm os.FileMode = 0755 // FIXME Why use something else?
		err := os.MkdirAll(outDirFlag, perm)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Output format
	if outFmtFlag == "" {
		outFmtFlag = "ascii"
	}

	for _, sw := range spData {
		var path string
		var perm os.FileMode = 0644 // FIXME Why use something else?

		if outDirFlag != "" {
			path = filepath.Join(outDirFlag, sw.fname)
		} else {
			path = filepath.Join(sw.dir, sw.fname)
		}

		err := sw.WriteFile(path, outFmtFlag, perm)
		if err != nil {
			log.Fatal(err)
		}
	}

	// --------------------------------------------------------------------------
	// TODO Take PLE values into account

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

}
