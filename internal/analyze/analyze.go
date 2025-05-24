package analyze

import (
	"fmt"
	"stone-analysis/internal/dft"
	"stone-analysis/internal/wav"
)

func Analyze(inFile string, n int) error {
	wavFile, err := wav.ReadWavFile(inFile)
	if err != nil {
		return fmt.Errorf("wav.ReadWavFile(%s): %w", inFile, err)
	}

	dftResult, err := dft.AnalyzeFrequencies(wavFile.Samples, float64(wavFile.FmtChunk.SampleRate), true)
	if err != nil {
		return fmt.Errorf("dft.AnalyzeFrequencies(): %w", err)
	}

	topFrequencies := dftResult.GetTopFrequencies(n)

	fmt.Printf("Top %d frequencies:\n", n)
	for _, freq := range topFrequencies {
		fmt.Printf("%.1f Hz\n", freq.Frequency)
	}

	return nil
}
