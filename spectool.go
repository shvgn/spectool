// Simple command line utility for the manipulation of columned ASCII data files
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
)

var (
	verboseFlag bool // Verbosity of the output

	keepEvFlag bool // Keep units in electron-volts
	keepNmFlag bool // Keep units in nanometers

	xFrom float64 // Cut X from value
	xTo   float64 // Cut X up to value

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

	// outFmtFlag string // Ouput format
	outDirFlag string // Ouput directory for resulting spectra
)

// Message on an arithmetic operation
func opMessage(op, value string) {
	if verboseFlag {
		fmt.Printf("  %s   %v\n", op, value)
	}
}

func init() {
	/*
	 Modificating flags
	*/
	// X Units
	flag.BoolVar(&keepEvFlag, "2ev", false, "Keep X in electron-volts")
	flag.BoolVar(&keepNmFlag, "2nm", false, "Keep X in nanometers")

	// X cutting options
	flag.Float64Var(&xFrom, "xfrom", math.Inf(-1), "X to start from")
	flag.Float64Var(&xTo, "xto", math.Inf(1), "X to end with")

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

	/*
	 Non-modificating flags
	*/
	// flag.BoolVar(&statsFlag, "s", false, "(Not implemented) Collect statistics on the data")
	flag.IntVar(&colXFlag, "colx", 1, "Set number of the X column in passed data files")
	flag.IntVar(&colYFlag, "coly", 2, "Set number of the Y column in passed data ASCII files")

	// flag.StringVar(&outFmtFlag, "of", "ascii", "Format of the output file (not implemented)")
	flag.StringVar(&outDirFlag, "od", "", "Directory for output files")

	flag.BoolVar(&verboseFlag, "v", false, "Verbose the actions")

	flag.Parse()
}

func main() {

	var (
		spectra   []*Spectrum
		spectrum  *Spectrum
		filePaths []string
		err       error
	)

	// Parse filenames. Those are considered to be paths or the parsing falls to globs
	// in order to work in both Windows cmd and Unix shells
	for _, cmdArg := range flag.Args() {

		// Parse from a file path
		spectrum, err = NewSpectrum(cmdArg, colXFlag, colYFlag)
		if err == nil {
			spectra = append(spectra, spectrum)
			continue
		}

		// Try a glob
		filePaths, err = filepath.Glob(cmdArg)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, filePath := range filePaths {
			spectrum, err = NewSpectrum(filePath, colXFlag, colYFlag)
			if err != nil {
				// Not strict about files parsing problems
				log.Printf("%s: Parse error: %q. Skip.", filePath, err)
				continue
			}
			spectra = append(spectra, spectrum)
		}
	}

	// Choose units for the processing and forbid using -2ev and -2nm together
	if keepEvFlag && keepNmFlag {
		log.Fatal("Cannot work on multiple X units simultaneously. Sorry.")
	}
	modifyUnits := keepEvFlag || keepNmFlag
	// Assign a killer that will be redefined fot the converion of X units
	ensureUnitsFunc := func(x float64) float64 {
		if modifyUnits { // Something went wrong
			log.Fatal("Unexpected units conversion")
			return math.NaN()
		}
		return x
	}

	// The string that will be added in the filename right before file extension
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
	if xFrom > xTo { // The X cutting flags are invalid
		log.Fatal("Invalid order of X cutting flags")
	}

	var (
		cutLeft  bool
		cutRight bool
		checker1 float64 // We want not to mess with the order of cutting values
		checker2 float64
	)

	if !math.IsInf(xFrom, -1) {
		cutLeft = true
		checker1 = xFrom
		xFrom = ensureUnitsFunc(xFrom)
	}
	if !math.IsInf(xTo, 1) {
		cutRight = true
		checker2 = xTo
		xTo = ensureUnitsFunc(xTo)
	}

	if modifyUnits && (cutLeft || cutRight) {
		// Check the order of cut boundaries
		if cutLeft && cutRight {
			if xFrom > xTo {
				// X order was possibly reversed by the X conversion
				xFrom, xTo = xTo, xFrom
			}
		} else {
			// Artificial checker is required because in the case of one-sided cut
			// we dont't have the second boundary to compare. Moreover, units conversion
			// is not strict on the present units of X. Therefore, we don't know
			// how to correspond passed cutting boundaries to boundaries of X
			if cutLeft {
				// Make checker2 > checker1
				checker2 = checker1 + (math.Abs(checker1)+1)*rand.Float64()
			}
			if cutRight {
				// Make checker1 < checker2
				checker1 = checker2 - (math.Abs(checker2)+1)*rand.Float64()
			}
			newChecker1, newChecker2 := ensureUnitsFunc(checker1), ensureUnitsFunc(checker2)
			if newChecker1 > newChecker2 {
				// Modified X has reversed order
				xFrom, xTo = xTo, xFrom
				cutLeft, cutRight = cutRight, cutLeft
			}
		}

	}

	// Arithmetics operands
	var addSpectrum, subSpectrum, mulSpectrum, divSpectrum *Spectrum
	if addFlag != "" {
		if addSpectrum, err = NewSpectrum(addFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			addSpectrum.xy.ModifyX(ensureUnitsFunc)
		}
	}
	if subFlag != "" {
		if subSpectrum, err = NewSpectrum(subFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			subSpectrum.xy.ModifyX(ensureUnitsFunc)
		}
	}
	if mulFlag != "" {
		if mulSpectrum, err = NewSpectrum(mulFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			mulSpectrum.xy.ModifyX(ensureUnitsFunc)
		}
	}
	if divFlag != "" {
		if divSpectrum, err = NewSpectrum(divFlag); err != nil {
			log.Fatal(err)
		}
		if modifyUnits {
			divSpectrum.xy.ModifyX(ensureUnitsFunc)
		}
	}

	// Processing
	l := len(spectra)
	for i, spectrum := range spectra {
		if verboseFlag && l > 1 {
			fmt.Println(fmt.Sprintf("%d/%d  ", i+1, l) + spectrum.dir + spectrum.fname)
		}

		// Subtract the noise from the full-length signal
		if noiseFlag {
			n := spectrum.xy.Noise()
			opMessage("-", fmt.Sprintf("%s (noise)", humanize.Ftoa(n)))
			spectrum.xy.ModifyY(func(y float64) float64 { return y - n })
			spectrum.AddOpSuffix("noise", humanize.Ftoa(n))
		}

		// Process the X units
		if modifyUnits {
			spectrum.xy.ModifyX(ensureUnitsFunc)
			spectrum.fname = addPreSuffix(spectrum.fname, unitsPreSuffix)
		}

		/*
			Cutting is done within boundaries and with spectra after making
			sure all X units are the same.
		*/
		cutFmt := "%s,%s"
		if cutLeft && cutRight {
			opMessage(">", humanize.Ftoa(xFrom))
			opMessage("<", humanize.Ftoa(xTo))
			spectrum.xy.Cut(xFrom, xTo)
			spectrum.AddOpSuffix("x", fmt.Sprintf(cutFmt, humanize.Ftoa(xFrom), humanize.Ftoa(xTo)))
			// sw.AddOpSuffix("from", humanize.Ftoa(fromFlag))
			// sw.AddOpSuffix("to", humanize.Ftoa(toFlag))
		} else {
			if cutLeft {
				opMessage(">", fmt.Sprintf("%v", humanize.Ftoa(xFrom)))
				xLast, _ := spectrum.xy.LastPoint()
				spectrum.xy.Cut(xFrom, xLast)
				spectrum.AddOpSuffix("x", fmt.Sprintf(cutFmt, humanize.Ftoa(xFrom), ""))
				// sw.AddOpSuffix("from", humanize.Ftoa(fromFlag))
			}
			if cutRight {
				opMessage("<", fmt.Sprintf("%v", humanize.Ftoa(xTo)))
				xFirst, _ := spectrum.xy.FirstPoint()
				spectrum.xy.Cut(xFirst, xTo)
				spectrum.AddOpSuffix("x", fmt.Sprintf(cutFmt, "", humanize.Ftoa(xTo)))
				// sw.AddOpSuffix("to", humanize.Ftoa(toFlag))
			}
		}

		/*
			Arithmetics with spectra
			Addition and subtracting of spectra should be done before noise calculation?
		*/
		if addFlag != "" {
			spectrum.xy.Add(addSpectrum.xy)
			opMessage("+", addSpectrum.fname)
			spectrum.AddSpOpSuffix("add", addSpectrum.fname)
		}
		if subFlag != "" {
			spectrum.xy.Subtract(subSpectrum.xy)
			opMessage("-", subSpectrum.fname)
			spectrum.AddSpOpSuffix("sub", subSpectrum.fname)
		}

		if mulFlag != "" {
			opMessage("*", mulSpectrum.fname)
			spectrum.xy.Multiply(mulSpectrum.xy)
			spectrum.AddSpOpSuffix("mul", mulSpectrum.fname)
		}
		if divFlag != "" {
			spectrum.xy.Divide(divSpectrum.xy)
			opMessage("/", divSpectrum.fname)
			spectrum.AddSpOpSuffix("div", divSpectrum.fname)
		}

		/*
			Arithmetics with numbers
		*/
		if addNumFlag != 0.0 {
			opMessage("+", fmt.Sprintf("%v", addNumFlag))
			spectrum.xy.ModifyY(func(y float64) float64 { return y + addNumFlag })
			spectrum.AddNumOpSuffix("add", addNumFlag)
		}
		if subNumFlag != 0.0 {
			opMessage("-", fmt.Sprintf("%v", subNumFlag))
			spectrum.xy.ModifyY(func(y float64) float64 { return y - subNumFlag })
			spectrum.AddNumOpSuffix("sub", subNumFlag)
		}
		if mulNumFlag != 1.0 {
			opMessage("*", fmt.Sprintf("%v", mulNumFlag))
			spectrum.xy.ModifyY(func(y float64) float64 { return y * mulNumFlag })
			spectrum.AddNumOpSuffix("mul", mulNumFlag)
		}
		if divNumFlag != 1.0 {
			opMessage("/", fmt.Sprintf("%v", divNumFlag))
			spectrum.xy.ModifyY(func(y float64) float64 { return y / divNumFlag })
			spectrum.AddNumOpSuffix("div", divNumFlag)
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
	// if outFmtFlag == "" {
	// 	outFmtFlag = "ascii"
	// }

	for _, spectrum := range spectra {
		var path string
		var perm os.FileMode = 0644 // FIXME Why use something else?

		if outDirFlag != "" {
			path = filepath.Join(outDirFlag, spectrum.fname)
		} else {
			path = filepath.Join(spectrum.dir, spectrum.fname)
		}

		// err := sw.WriteFile(path, outFmtFlag, perm)
		err := spectrum.WriteFile(path, "ascii", perm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
