package dft

import "math"

func ApplyHammingWindow(samples []float64) []float64 {
	n := len(samples)
	if n == 0 {
		return samples
	}

	windowed := make([]float64, n)

	for i := 0; i < n; i++ {
		window := 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(n-1))
		windowed[i] = samples[i] * window
	}

	return windowed
}

func ApplyHannWindow(samples []float64) []float64 {
	n := len(samples)
	if n == 0 {
		return samples
	}

	windowed := make([]float64, n)

	for i := 0; i < n; i++ {
		window := 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(n-1)))
		windowed[i] = samples[i] * window
	}

	return windowed
}

func ApplyBlackmanWindow(samples []float64) []float64 {
	n := len(samples)
	if n == 0 {
		return samples
	}

	windowed := make([]float64, n)

	a0 := 0.42
	a1 := 0.5
	a2 := 0.08

	for i := 0; i < n; i++ {
		arg := 2 * math.Pi * float64(i) / float64(n-1)
		window := a0 - a1*math.Cos(arg) + a2*math.Cos(2*arg)
		windowed[i] = samples[i] * window
	}

	return windowed
}
