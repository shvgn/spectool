// The code is provided "as is" without any warranty and shit.
// You are free do anything you want with it.
//
// Evgenii Shevchenko a.k.a @shvgn
// 2014

package main

import (
	"flag"
	"fmt"
	"github.com/shvgn/spectrum"
	// "strconv"
	// "strings"
)

const (
	EVNM           = 1239.84193        // Coefficient equal to nm*eV
	MAX_ENERGY     = 8.0               // electron-volts
	MIN_WAVELENGTH = EVNM / MAX_ENERGY // nanometers
)

// Converter from nanometers to electron-volts and in reverse
func conv_evnm(x float64) float64 {
	return EVNM / x
}

func main() {

	// --------------------------------------------------------------------------
	// We start with flags. All other arguments are supposed to be text files with spectra data

	pleSet := flag.String("ple", "", "This is set of wavelength of energy walues (in nm or ev) divided by commas for PLE extraction e.g. -ple=287.5,288,288.5")
	// averPtr := flag.Int("aver", 0, "Specifies number of neighbour values to take into account. The exact or neares value is taken if aver=0, if e.g. aver=2 than two more values are taken from both sides if possible resulting in averaging of 5 values.")

	nm2evPtr := flag.Bool("nm2ev", false, "Set this flag in order to convert X from nanometers to electron-volts")
	ev2nmPtr := flag.Bool("ev2nm", false, "Set this flag in order to convert X from electron-volts to nanometers")

	flag.Parse()

	// --------------------------------------------------------------------------
	// Filling the data
	spectra := make([]*spectrum.Spectrum, 1)

	for _, filePath := range flag.Args() {
		// specPtr := spectrum.NewSpectrum()
		// var specPtr *spectrum.Spectrum
		specPtr, err := spectrum.SpectrumFromFile(filePath)
		if err != nil {
			fmt.Println("Cannot read spectrum from file", filePath+":", err.Error(), "- Skipping.")
			continue
		}
		specPtr.ReadFromFile(filePath)
		fmt.Println(specPtr)
		spectra = append(spectra, specPtr)
	}

	// --------------------------------------------------------------------------
	// Processing
	// for _, sp := range spectra {
	// if flag1 do thing1
	// if flag2 do thing2
	// and so on
	// }

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
	fmt.Println("PLE detection values passed: ", *pleSet)
	fmt.Println("nm to eV: ", *nm2evPtr)
	fmt.Println("eV to nm: ", *ev2nmPtr)
	// fmt.Println("PLE detection values parsed: ", pleVals)
	fmt.Println("Other cmd arguments: ", flag.Args())
}

/************************************************************************

Interface

1. specify operation with the special key?

	spectool -op=nm2ev file1 file2 file3 ...

2. seaprate key for each operation

	spectool -nm2ev file1 file2 file3 ...
	spectool -ple file1 file2 file3 ...
	spectool -mean file1 file2 file3 ...
	spectool -stats file1 file2 file3 ...

	spectool -stats -ple -nm2ev file1 file2 file3 ...
		or spectool -stats -ple=all -nm2ev file1 file2 file3 ...



*************************************************************************/
