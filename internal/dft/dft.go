package dft

import "math"

func DFT(samples []float64, sampleRate float64) (*DFTResult, error) {
	if len(samples) == 0 {
		return &DFTResult{}, nil
	}

	originalLength := len(samples)

	if !isPowerOf2(originalLength) {
		spectrum := make([]Complex, originalLength)
		for k := 0; k < originalLength; k++ {
			var sum Complex
			for n := 0; n < originalLength; n++ {
				angle := -2 * math.Pi * float64(k) * float64(n) / float64(originalLength)
				sum.Real += samples[n] * math.Cos(angle)
				sum.Imag += samples[n] * -math.Sin(angle)
			}
			spectrum[k] = sum
		}

		numComponents := originalLength / 2
		components := make([]FrequencyComponent, numComponents)
		freqResolution := sampleRate / float64(originalLength)

		for i := 0; i < numComponents; i++ {
			mag := spectrum[i].Magnitude()
			var magnitude float64
			if i == 0 {
				magnitude = mag / float64(originalLength)
			} else {
				magnitude = 2 * mag / float64(originalLength)
			}
			components[i] = FrequencyComponent{
				Frequency: float64(i) * freqResolution,
				Magnitude: magnitude,
				Phase:     spectrum[i].Phase(),
				Real:      spectrum[i].Real,
				Imag:      spectrum[i].Imag,
			}
		}

		var nyquist Complex
		if originalLength%2 == 0 {
			nyquist = spectrum[originalLength/2]
		}

		return &DFTResult{
			Components:     components,
			Nyquist:        nyquist,
			SampleRate:     sampleRate,
			SampleCount:    originalLength,
			FreqResolution: freqResolution,
		}, nil
	}

	paddedLength := nextPowerOf2(originalLength)
	complexInput := make([]Complex, paddedLength)
	for i := 0; i < originalLength; i++ {
		complexInput[i] = Complex{Real: samples[i], Imag: 0}
	}

	fftResult, err := FFT(complexInput)
	if err != nil {
		return nil, err
	}

	numComponents := paddedLength / 2
	components := make([]FrequencyComponent, numComponents)
	freqResolution := sampleRate / float64(paddedLength)

	for i := 0; i < numComponents; i++ {
		mag := fftResult[i].Magnitude()
		var magnitude float64
		if i == 0 {
			magnitude = mag / float64(paddedLength)
		} else {
			magnitude = 2 * mag / float64(paddedLength)
		}
		components[i] = FrequencyComponent{
			Frequency: float64(i) * freqResolution,
			Magnitude: magnitude,
			Phase:     fftResult[i].Phase(),
			Real:      fftResult[i].Real,
			Imag:      fftResult[i].Imag,
		}
	}

	var nyquist Complex
	if paddedLength%2 == 0 {
		nyquist = fftResult[paddedLength/2]
	}

	return &DFTResult{
		Components:     components,
		Nyquist:        nyquist,
		SampleRate:     sampleRate,
		SampleCount:    originalLength,
		FreqResolution: freqResolution,
	}, nil
}

func IDFT(result *DFTResult) ([]float64, error) {
	if len(result.Components) == 0 {
		return []float64{}, nil
	}

	compLen := len(result.Components)
	paddedLength := compLen * 2
	complexSpectrum := make([]Complex, paddedLength)

	for i, comp := range result.Components {
		complexSpectrum[i] = Complex{Real: comp.Real, Imag: comp.Imag}
	}

	complexSpectrum[compLen] = result.Nyquist

	for i := 1; i < compLen; i++ {
		idx := paddedLength - i
		complexSpectrum[idx] = Complex{
			Real: result.Components[i].Real,
			Imag: -result.Components[i].Imag,
		}
	}

	timeDomain, err := IFFT(complexSpectrum)
	if err != nil {
		return nil, err
	}

	samples := make([]float64, result.SampleCount)
	for i := 0; i < result.SampleCount && i < len(timeDomain); i++ {
		samples[i] = timeDomain[i].Real
	}

	return samples, nil
}

func AnalyzeFrequencies(samples []float64, sampleRate float64, useWindowing bool) (*DFTResult, error) {
	processedSamples := samples
	if useWindowing && len(samples) > 1 {
		processedSamples = ApplyHammingWindow(samples)
	}
	return DFT(processedSamples, sampleRate)
}
