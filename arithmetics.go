package main

const (
	lightSpeed     float64 = 299792458                          // Meters per second
	planckConstant float64 = 4.135667516e-15                    // Electron-volt * second, h
	eVtoNm         float64 = planckConstant * 1e+9 * lightSpeed // Factor of nanometer and electron-volt
	maxEnergy      float64 = 10.0                               // Electron-volts
	minWavelength  float64 = eVtoNm / maxEnergy                 // Nanometers
)

// ConvEvNm converts a value from nanometers to electron-volts and in reverse
func ConvEvNm(x float64) float64 {
	return eVtoNm / x
}

// ensureEv ensures the argument is in electron-volts using maxEnergy constant
func ensureEv(x float64) float64 {
	if x < maxEnergy {
		return x
	}
	return ConvEvNm(x)
}

// ensureNm ensures the argument is in nanometers using minWavelength constant
func ensureNm(x float64) float64 {
	if x > minWavelength {
		return x
	}
	return ConvEvNm(x)
}
