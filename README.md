spectool
========

A command line tool for processing an ASCII two-column data files (X and Y) containing spectroscopic data (mostly photoluminescence, transmission and photocurrent)


Data
----

The ASCII data file can be just a two-column text with numeric data separated by tabs or spaces. In common the data is expected to look like this

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

`spectool -help` for help and usage message
`spectool -nm2ev file1 file2 file3`for converting nanometers to electron-volts

```spectool -ple   file1 file2 file3 ...``` for extraction of PL excitation spectra

```spectool -mean  file1 file2 file3 ...``` merge spectra into a mean one

```spectool -stats file1 file2 file3 ...``` generate statistics file (where)


Multiple tasks could look like this 
```
spectool -nm2ev -ydiv=excitationCurveFile -ysub=darkNoiseFile \
         -xstart=230 -xend=320 \
         -stats -ple=245,245.5,245.75,238 \
         file1 file2 file3 ...
```



Cutting
-------

This makes cut of a data. It cuts all points whosw _x_'s exceed the range of [xstart,xend] and throw away points whose _y_'s exceed the range of [ystart,yend]
```
spectool -xstart=float -xend=float -ystart=float -yend=float file1 file2 file3 ...
``` 



Arithmetics
-----------
This must be useful for both float numbers and reference data files. 
E.g. subtract known noise level or known spectrum of a scattered light (dark spectrum),
divide by used filters curves or scale factors, and so on. Addition and subtraction are 
of higher priority than multiplication and division when they are used simultaneously.
```
spectool -yadd=float file1 file2 file3 ...
spectool -yadd=fileS file1 file2 file3 ...
spectool -ymul=... file1 file2 file3 ...
spectool -ysub=... file1 file2 file3 ...
spectool -ydiv=... file1 file2 file3 ...
``` 


Tasks
-----
- [ ] naming and placing of files
- [ ] arithmetic operations including interpolations
- [ ] calculation of 
  - [ ] noise level
  - [ ] statistics data including 
    - [ ] area
    - [ ] maximum position
    - [ ] full width at half-maximum (FWHM) for main peaks
- [ ] analisys of X units for useful predictions on X conversion (-xconv key instead of both -nm2ev and -ev2nm)

