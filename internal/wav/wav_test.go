package wav

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFourCCString(t *testing.T) {
	fourcc := FourCC{'R', 'I', 'F', 'F'}
	if fourcc.String() != "RIFF" {
		t.Errorf("Expected 'RIFF', got '%s'", fourcc.String())
	}
}

func TestFourCCEquals(t *testing.T) {
	fourcc := FourCC{'W', 'A', 'V', 'E'}
	if !fourcc.Equals("WAVE") {
		t.Errorf("Expected fourcc to equal 'WAVE'")
	}
	if fourcc.Equals("RIFF") {
		t.Errorf("Expected fourcc to not equal 'RIFF'")
	}
}

func createTestWavFile(t *testing.T, path string) {
	header := []byte{
		'R', 'I', 'F', 'F',
		36, 0, 0, 0,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		16, 0, 0, 0,
		1, 0,
		1, 0,
		0x80, 0xBB, 0, 0,
		0, 0, 0, 0,
		0, 0,
		16, 0,
		'd', 'a', 't', 'a',
		0, 0, 0, 0,
	}

	bytesPerSample := 16 / 8
	byteRate := 48000 * 1 * bytesPerSample
	header[28] = byte(byteRate & 0xFF)
	header[29] = byte((byteRate >> 8) & 0xFF)
	header[30] = byte((byteRate >> 16) & 0xFF)
	header[31] = byte((byteRate >> 24) & 0xFF)

	blockAlign := 1 * bytesPerSample
	header[32] = byte(blockAlign & 0xFF)
	header[33] = byte((blockAlign >> 8) & 0xFF)

	err := os.WriteFile(path, header, 0644)
	if err != nil {
		t.Fatalf("Failed to create test WAV file: %v", err)
	}
}

func TestReadWavFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test.wav")
	createTestWavFile(t, testFilePath)

	wavFile, err := ReadWavFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read WAV file: %v", err)
	}

	if !wavFile.Header.ChunkID.Equals("RIFF") {
		t.Errorf("Expected ChunkID 'RIFF', got '%s'", wavFile.Header.ChunkID.String())
	}
	if !wavFile.Header.Format.Equals("WAVE") {
		t.Errorf("Expected Format 'WAVE', got '%s'", wavFile.Header.Format.String())
	}

	if !wavFile.FmtChunk.SubChunkID.Equals("fmt ") {
		t.Errorf("Expected SubChunkID 'fmt ', got '%s'", wavFile.FmtChunk.SubChunkID.String())
	}
	if wavFile.FmtChunk.AudioFormat != 1 {
		t.Errorf("Expected AudioFormat 1 (PCM), got %d", wavFile.FmtChunk.AudioFormat)
	}
	if wavFile.FmtChunk.NumChannels != 1 {
		t.Errorf("Expected NumChannels 1 (mono), got %d", wavFile.FmtChunk.NumChannels)
	}
	if wavFile.FmtChunk.SampleRate != 48000 {
		t.Errorf("Expected SampleRate 48000, got %d", wavFile.FmtChunk.SampleRate)
	}
	if wavFile.FmtChunk.BitsPerSample != 16 {
		t.Errorf("Expected BitsPerSample 16, got %d", wavFile.FmtChunk.BitsPerSample)
	}

	if !wavFile.DataChunk.SubChunkID.Equals("data") {
		t.Errorf("Expected SubChunkID 'data', got '%s'", wavFile.DataChunk.SubChunkID.String())
	}
}

func TestConvertSamples(t *testing.T) {
	sampleData := []byte{
		0x00, 0x00,
		0xFF, 0x3F,
		0x00, 0x00,
		0x01, 0xC0,
	}

	reader := NewWavReader()
	samples := reader.ConvertToSamples(sampleData)

	expectedSamples := []float64{0, 0.5, 0, -0.5}

	if len(samples) != len(expectedSamples) {
		t.Fatalf("Expected %d samples, got %d", len(expectedSamples), len(samples))
	}

	const epsilon = 0.01
	for i, expected := range expectedSamples {
		if samples[i] < expected-epsilon || samples[i] > expected+epsilon {
			t.Errorf("Sample %d: expected %.2f, got %.2f", i, expected, samples[i])
		}
	}
}

func TestValidateWavFormat(t *testing.T) {
	validWav := &WavFile{
		FmtChunk: FmtSubChunk{
			AudioFormat:   1,
			NumChannels:   1,
			SampleRate:    48000,
			BitsPerSample: 16,
		},
	}

	err := ValidateWavFormat(validWav)
	if err != nil {
		t.Errorf("Expected valid WAV format, got error: %v", err)
	}

	invalidFormat := *validWav
	invalidFormat.FmtChunk.AudioFormat = 2
	err = ValidateWavFormat(&invalidFormat)
	if err != ErrUnsupportedAudioFormat {
		t.Errorf("Expected ErrUnsupportedAudioFormat, got %v", err)
	}

	invalidChannels := *validWav
	invalidChannels.FmtChunk.NumChannels = 2
	err = ValidateWavFormat(&invalidChannels)
	if err != ErrInvalidNumChannels {
		t.Errorf("Expected ErrInvalidNumChannels, got %v", err)
	}

	invalidRate := *validWav
	invalidRate.FmtChunk.SampleRate = 44100
	err = ValidateWavFormat(&invalidRate)
	if err != ErrInvalidSampleRate {
		t.Errorf("Expected ErrInvalidSampleRate, got %v", err)
	}

	invalidBits := *validWav
	invalidBits.FmtChunk.BitsPerSample = 8
	err = ValidateWavFormat(&invalidBits)
	if err != ErrInvalidBitsPerSample {
		t.Errorf("Expected ErrInvalidBitsPerSample, got %v", err)
	}
}

func TestRoundTripConversion(t *testing.T) {
	originalSamples := []float64{0, 0.5, 0, -0.5, 0.25, -0.75, 1.0, -1.0}

	writer := NewWavWriter()
	bytes := writer.ConvertFromSamples(originalSamples)

	reader := NewWavReader()
	resultSamples := reader.ConvertToSamples(bytes)

	if len(resultSamples) != len(originalSamples) {
		t.Fatalf("Expected %d samples after round trip, got %d",
			len(originalSamples), len(resultSamples))
	}

	const epsilon = 0.001
	for i, expected := range originalSamples {
		if resultSamples[i] < expected-epsilon || resultSamples[i] > expected+epsilon {
			t.Errorf("Sample %d after round trip: expected %.3f, got %.3f",
				i, expected, resultSamples[i])
		}
	}
}

func TestWriteWavFile(t *testing.T) {
	samples := []float64{0, 0.5, 0, -0.5, 0, 0.5, 0, -0.5}

	wavFile := &WavFile{
		Header: WavHeader{
			ChunkID:   FourCC{'R', 'I', 'F', 'F'},
			ChunkSize: 36 + uint32(len(samples)*2),
			Format:    FourCC{'W', 'A', 'V', 'E'},
		},
		FmtChunk: FmtSubChunk{
			SubChunkID:    FourCC{'f', 'm', 't', ' '},
			SubChunkSize:  16,
			AudioFormat:   1,
			NumChannels:   1,
			SampleRate:    48000,
			ByteRate:      48000 * 1 * 2,
			BlockAlign:    2,
			BitsPerSample: 16,
		},
		DataChunk: DataSubChunk{
			SubChunkID:   FourCC{'d', 'a', 't', 'a'},
			SubChunkSize: uint32(len(samples) * 2),
		},
		Samples: samples,
	}

	tmpDir := t.TempDir()
	outFilePath := filepath.Join(tmpDir, "output.wav")

	err := WriteWavFile(outFilePath, wavFile)
	if err != nil {
		t.Fatalf("Failed to write WAV file: %v", err)
	}

	readWav, err := ReadWavFile(outFilePath)
	if err != nil {
		t.Fatalf("Failed to read written WAV file: %v", err)
	}

	if len(readWav.Samples) != len(samples) {
		t.Fatalf("Expected %d samples, got %d", len(samples), len(readWav.Samples))
	}

	const epsilon = 0.001
	for i, expected := range samples {
		if readWav.Samples[i] < expected-epsilon || readWav.Samples[i] > expected+epsilon {
			t.Errorf("Sample %d: expected %.3f, got %.3f",
				i, expected, readWav.Samples[i])
		}
	}
}
