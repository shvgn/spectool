// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code in any way you wish.
//
// Evgeny Shevchenko
// shvgn@protonmail.ch
// 2015

package main

const (
	LIGHT_SPEED     float64 = 299792458                            // meters per second
	PLANCK_CONSTANT float64 = 4.135667516e-15                      // electronvolts * second, h
	EVNM            float64 = PLANCK_CONSTANT * 1e+9 * LIGHT_SPEED // factor of nanometers and electron-volts
	MAX_ENERGY      float64 = 10.0                                 // electron-volts
	MIN_WAVELENGTH  float64 = EVNM / MAX_ENERGY                    // nanometers
)

// Converter from nanometers to electron-volts and in reverse
func ConvEvNm(x float64) float64 {
	return EVNM / x
}

// Function to ensure the argument is in electron-volts
func ensureEv(x float64) float64 {
	if x < MAX_ENERGY {
		return x
	}
	return ConvEvNm(x)
}

// Function to ensure the argument is in nanometers
func ensureNm(x float64) float64 {
	if x > MIN_WAVELENGTH {
		return x
	}
	return ConvEvNm(x)
}
