spectool
========
[![GoDoc](https://godoc.org/github.com/shvgn/spectool?status.svg)](https://godoc.org/github.com/shvgn/spectool)

Simple command line tool for processing ASCII files with columns of numeric data. Spectool features some capabilities for spectroscopic data (optical in particular) though it is useful for some basic arithmetic processing of arbitrary data. The goal is time saving before building plots in GUI software which usually require too much clicks. It is can be downloaded as binary files for Linux, OS X and Windows. The interface and code are not quite polished yet.


Data
----

The ASCII data must contail multiple-column text with numeric data separated by tabs or spaces. In general the data is expected to look like this

```
# Lines starting with '#' are ignored
# Either do empty lines

# Space characters are assumed to be delimeters by default

# Columns are counted from 1
#1   2   3   4
A1  A2  A3  A4
B1  B2  B3  B4
C1  C2  C3  C4
D1  D2  D3  D4
E1  E2  E3  E4
F1  F2  F3  F4
```

By default, X and Y are taken as columns 1 and 2 respectively. To take other columns in account one can use flags `-xcol` and `-ycol`


Interface
---------

### Arithmetics

For simplicity keys for arithmetic operations are distinguished for spectra and numbers. The cure.dat file must contain valid data as other files.
```
spectool -ndiv=1000 file1 file2 file3 ...
spectool -spdiv=/path/to/calibration/curve.dat file1 file2 file3 ...
``` 


### Spectrum capabilities

`spectool -2ev file1 file2 file3` for keeping/converting X units in electron-volts, fileneames being renamed from filename.ext or filename.nm.ext to filename.ev.ext

`spectool -2nm file1 file2 file3` for keeping/converting X units in nanometers, fileneames being renamed from filename.ext or filename.ev.ext to filename.nm.ext

Multiple tasks could look like this 
```
spectool -2ev -n -spdiv=ApparatusSpectra.dat -xfrom=230 -xto=320 -od=res -s -v spectrum*.txt
```
Which means calculate and subtract noise (`-n`), keep X values in electron-volts (`-2ev`), cut the data from 230 nm to 320 nm (these values will also be converted and used in electron-volts), show some analysis data (`-s`) and verbose output (`-v`) and create direcory 'res' (`-od`) to put the resulting ascii files there. These new files will be named as follows (e.g. we took file spectrum1.txt as an input for the command above):
```
spectrum1.noise[1.34].x[3.874506,5.390617].div[ApparatusSpectra.dat].ev.txt
```

### Options
```
Usage of spectool:
  -2ev
    	Keep X in electron-volts
  -2nm
    	Keep X in nanometers
  -colx int
    	Set number of the X column in passed data files (default 1)
  -coly int
    	Set number of the Y column in passed data ASCII files (default 2)
  -n	Subtract noise
  -nadd float
    	Add a number
  -ndiv float
    	Divide by a number  (default 1)
  -nmul float
    	Multiply by a number  (default 1)
  -nsub float
    	Subtract a number
  -od string
    	Directory for output files
  -of string
    	[ascii|tsv|csv]   Format of the output file (default "ascii")
  -spadd string
    	Add spectrum
  -spdiv string
    	Divide by spectrum
  -spmul string
    	Multiply by spectrum
  -spsub string
    	Subtract spectrum
  -v	Verbose the actions
  -xfrom float
    	X to start from (default -Inf)
  -xto float
    	X to end with (default +Inf)
```




TODO
-----
- [x] naming and placing of resulting files
- [x] arithmetic operations involving interpolations
- [ ] calculation of 
  - [x] noise level
  - [ ] mean spectra
  - [ ] metadata (`-s`)
    - [x] area under curve
    - [x] maximum position (x,y)
    - [ ] full width at half-maximum (FWHM) for main peaks
    - [ ] position of FWHM's for main peaks
- [x] One-way cutting of X
- [ ] Smoothing by Savitsky-Golay or Holoborodko filters ([must be in xy](https://github.com/shvgn/xy))
- [ ] Divide all resulting data by maximum hight/area of a reference spectrum, so they could be always compared in relative units to some reference data (e.g. `-div-by-max`, `-div-by-area=/path/to/referenceSpectrum.txt`)
- [ ] Separate peaks by minimum (x,y) in a region passed by extra keys (`-sepfrom`, `-septo`). Useful for precessing of separate peaks
- [ ] Input formats (e.g. `-if=csv`) ([must be in xy](https://github.com/shvgn/xy))
  - [x] ASCII/TSV
  - [ ] CSV
  - [ ] JSON
- [ ] Ouput formats (e.g. `-of=matlab`)([must be in xy](https://github.com/shvgn/xy))
  - [x] ASCII/TSV
  - [ ] CSV
  - [ ] Matlab 2-D array
  - [ ] JSON
  - [ ] Origin-friendly?
