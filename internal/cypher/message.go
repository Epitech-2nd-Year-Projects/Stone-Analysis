package cypher

import (
	"fmt"
	"math"
	"stone-analysis/internal/dft"
	"strings"
)

func validateMessage(message string) error {
	for _, char := range strings.ToUpper(message) {
		if _, exists := charToFreqMap[char]; !exists {
			return fmt.Errorf("unsupported character in message: '%c' (ASCII: %d)", char, int(char))
		}
	}
	return nil
}

func embedCharacterInSpectrum(components []dft.FrequencyComponent, char rune, sampleRate float64, numSamples int) error {
	upperChar := rune(strings.ToUpper(string(char))[0])
	freqs, exists := charToFreqMap[upperChar]
	if !exists {
		return fmt.Errorf("character '%c' not found in frequency mapping", char)
	}

	const modificationFactor = 1.01

	for _, freq := range freqs {
		bin := getFrequencyBin(float64(freq), sampleRate, numSamples)

		if bin >= 0 && bin < len(components) {
			originalMag := components[bin].Magnitude

			components[bin].Magnitude *= modificationFactor

			phase := components[bin].Phase
			components[bin].Real = components[bin].Magnitude * math.Cos(phase)
			components[bin].Imag = components[bin].Magnitude * math.Sin(phase)

			_ = originalMag
		}
	}

	return nil
}

func addMessageMarkers(components []dft.FrequencyComponent, messageLength int, sampleRate float64, numSamples int) {
	startMarkerFreq := 15000.0
	lengthMarkerBase := 16000.0

	startBin := getFrequencyBin(startMarkerFreq, sampleRate, numSamples)
	if startBin >= 0 && startBin < len(components) {
		components[startBin].Magnitude *= 1.02
		phase := components[startBin].Phase
		components[startBin].Real = components[startBin].Magnitude * math.Cos(phase)
		components[startBin].Imag = components[startBin].Magnitude * math.Sin(phase)
	}

	for i := 0; i < 8; i++ {
		if (messageLength & (1 << i)) != 0 {
			lengthFreq := lengthMarkerBase + float64(i*100)
			lengthBin := getFrequencyBin(lengthFreq, sampleRate, numSamples)

			if lengthBin >= 0 && lengthBin < len(components) {
				components[lengthBin].Magnitude *= 1.015
				phase := components[lengthBin].Phase
				components[lengthBin].Real = components[lengthBin].Magnitude * math.Cos(phase)
				components[lengthBin].Imag = components[lengthBin].Magnitude * math.Sin(phase)
			}
		}
	}
}
