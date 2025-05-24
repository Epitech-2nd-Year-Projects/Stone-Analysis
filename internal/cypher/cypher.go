package cypher

import (
	"fmt"
	"stone-analysis/internal/dft"
	"stone-analysis/internal/wav"
	"strings"
)

func Cypher(inFile, outFile, message string) error {
	if err := validateMessage(message); err != nil {
		return fmt.Errorf("message validation failed: %w", err)
	}

	if len(message) == 0 {
		return fmt.Errorf("message cannot be empty")
	}

	if len(message) > 255 {
		return fmt.Errorf("message too long (max 255 characters, got %d)", len(message))
	}

	wavFile, err := wav.ReadWavFile(inFile)
	if err != nil {
		return fmt.Errorf("failed to read input WAV file '%s': %w", inFile, err)
	}

	dftResult, err := dft.DFT(wavFile.Samples, float64(wavFile.FmtChunk.SampleRate))
	if err != nil {
		return fmt.Errorf("DFT transformation failed: %w", err)
	}

	addMessageMarkers(dftResult.Components, len(message),
		float64(wavFile.FmtChunk.SampleRate), len(wavFile.Samples))

	upperMessage := strings.ToUpper(message)
	for i, char := range upperMessage {
		err := embedCharacterInSpectrum(dftResult.Components, char,
			float64(wavFile.FmtChunk.SampleRate), len(wavFile.Samples))
		if err != nil {
			return fmt.Errorf("failed to embed character '%c' at position %d: %w", char, i, err)
		}
	}

	modifiedSamples, err := dft.IDFT(dftResult)
	if err != nil {
		return fmt.Errorf("IDFT transformation failed: %w", err)
	}

	outputWav := &wav.WavFile{
		Header:    wavFile.Header,
		FmtChunk:  wavFile.FmtChunk,
		DataChunk: wavFile.DataChunk,
		Samples:   modifiedSamples,
	}

	dataSize := uint32(len(modifiedSamples) * 2)
	outputWav.Header.ChunkSize = 36 + dataSize
	outputWav.DataChunk.SubChunkSize = dataSize

	err = wav.WriteWavFile(outFile, outputWav)
	if err != nil {
		return fmt.Errorf("failed to write output WAV file '%s': %w", outFile, err)
	}

	return nil
}
