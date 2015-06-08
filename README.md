spectool
========

This is a simple command line tool for processing ASCII two-column data files (X and Y) containing spectroscopic data (mostly photoluminescence, transmission and photocurrent). The goal is combinig the functionality of shvgn/py_spectool scripts in one binary program that doesn't require any runtime setup anywhere (especially on Windows platform). The interface and code are not quite polished yet.


Data
----

The ASCII data file can be just a two-column text with numeric data separated by tabs or spaces. In general the data is expected to look like this

```
# Lines with the # symbol in the beginning are ignored
# Either do empty lines

# Space characters are assumed to be the delimeters.

# Columns are counted from 1
#1	2	3	4
A1	B1	C1	D1
A2	B2	C2	D2
A3	B3	C3	D3
A4	B4	C4	D4
A5	B5	C5	D5
A6	B6	C6	D6
A7	B7	C7	D7
A8	B8	C8	D8
```

X and Y are taken as columns 1 and 2 respectively. To take other columns in account one can use flags `xcol` and `ycol`


Interface
---------

`spectool -2ev file1 file2 file3` for keeping/converting X units in electron-volts, fileneames being renamed from filename.ext or filename.nm.ext to filename.ev.ext

`spectool -2nm file1 file2 file3` for keeping/converting X units in nanometers, fileneames being renamed from filename.ext or filename.ev.ext to filename.nm.ext

Multiple tasks could look like this 
```
spectool -2ev -n -spdiv=ApparatusSpectra.dat -xfrom=230 -xto=320 -od=res -s -v spectrum*.txt
```
Which means calculate and subtract noise (```-n```), keep X values in electron-volts (```-2ev```), cut the data from 230 nm to 320 nm (these values will also be converted and used in electron-volts), show some analysis data (```-s```) and verbose output (```-v```) and create direcory 'res' (```-od```) to put the resulting ascii files there. The new files will be named as follows (e.g. we took file spec1.txt as an input for the command above):
``` spec1.noise[1.34].div[ApparatusSpectra.dat].ev.txt ```

```
spectool -h
Usage of spectool:
  -2ev=false: Keep X in electron-volts
  -2nm=false: Keep X in nanometers
  -colx=1: Set number of the X column in passed data files
  -coly=2: Set number of the Y column in passed data ASCII files
  -mean=false: (Not implemented) Mean spectrum from all the passed data
  -n=false: Subtract noise
  -nadd=0: Add a number 
  -ndiv=1: Divide by a number 
  -nmul=1: Multiply by a number 
  -nsub=0: Subtract a number 
  -od="": Directory for output files. If not specified new files are placed near original ones
  -of="ascii": [ascii|tsv|csv]   Format of the output file
  -s=false: (Not implemented) Collect statistics on the data
  -smooth="": [ws,po]   (Not implemented) Smooth data with optionally specified both window size and polynome order
  -spadd="": Add spectrum
  -spdiv="": Divide by spectrum
  -spmul="": Multiply by spectrum
  -spsub="": Subtract spectrum
  -v=false: Verbose the actions
  -xfrom=-1: X to start from
  -xto=-1: X to end with
```


Arithmetics
-----------
For simplicity keys for arithmetic operations are distinguished for spectra and numbers. The cure.dat file must contain valid data as other files.
```
spectool -ndiv=1000 file1 file2 file3 ...
spectool -spdiv=/path/to/calibration/curve.dat file1 file2 file3 ...
``` 


TODO
-----
- [x] naming and placing of resulting files
- [x] arithmetic operations including interpolations
- [ ] calculation of 
  - [x] noise level
  - [ ] mean spectra
  - [ ] metadata (```-s```)
    - [ ] area under the curve
    - [ ] maximum position (x,y)
    - [ ] full width at half-maximum (FWHM) for main peaks
    - [ ] position of FWHM's for main peaks
- [ ] One-way cutting of X
- [ ] Smoothing by Savitsky-Golay or Holoborodko filters 
- [ ] Divide all resulting data by maximum hight/area of a reference spectrum, so there can be always compared in relative units to some reference (-div-by-max, -div-by-area=/path/to/referenceSpectrum.txt)
- [ ] Separate peaks by minimum (x,y) in a region passed by extra keys (-sepfrom, -septo). Useful for distinct precessing of peaks.
- [ ] Input formats (-if=csv)
  - [x] ASCII/TSV
  - [ ] CSV
- [ ] Ouput formats (-of=matlab)
  - [x] ASCII/TSV
  - [ ] CSV
  - [ ] Matlab 2-D array
  - [ ] JSON
  - [ ] Origin-friendly?
