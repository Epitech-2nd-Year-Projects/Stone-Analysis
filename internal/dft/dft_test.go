package dft

import (
	"math"
	"testing"
)

const epsilon = 1e-10

func TestComplexOperations(t *testing.T) {
	c1 := Complex{Real: 3, Imag: 4}
	c2 := Complex{Real: 1, Imag: 2}

	sum := c1.Add(c2)
	if sum.Real != 4 || sum.Imag != 6 {
		t.Errorf("Expected (4, 6), got (%.2f, %.2f)", sum.Real, sum.Imag)
	}

	diff := c1.Sub(c2)
	if diff.Real != 2 || diff.Imag != 2 {
		t.Errorf("Expected (2, 2), got (%.2f, %.2f)", diff.Real, diff.Imag)
	}

	prod := c1.Mul(c2)
	if prod.Real != -5 || prod.Imag != 10 {
		t.Errorf("Expected (-5, 10), got (%.2f, %.2f)", prod.Real, prod.Imag)
	}

	mag := c1.Magnitude()
	expected := 5.0
	if math.Abs(mag-expected) > epsilon {
		t.Errorf("Expected magnitude %.2f, got %.2f", expected, mag)
	}

	phase := c1.Phase()
	expectedPhase := math.Atan2(4, 3)
	if math.Abs(phase-expectedPhase) > epsilon {
		t.Errorf("Expected phase %.6f, got %.6f", expectedPhase, phase)
	}
}

func TestFFTErrors(t *testing.T) {
	_, err := FFT([]Complex{})
	if err != ErrEmptyInput {
		t.Errorf("Expected ErrEmptyInput, got %v", err)
	}

	input := []Complex{
		{Real: 1, Imag: 0},
		{Real: 2, Imag: 0},
		{Real: 3, Imag: 0},
	}
	_, err = FFT(input)
	if err != ErrInvalidFFTSize {
		t.Errorf("Expected ErrInvalidFFTSize, got %v", err)
	}

	validInput := []Complex{
		{Real: 1, Imag: 0},
		{Real: 0, Imag: 0},
		{Real: 0, Imag: 0},
		{Real: 0, Imag: 0},
	}
	result, err := FFT(validInput)
	if err != nil {
		t.Errorf("Expected no error for valid input, got %v", err)
	}
	if len(result) != 4 {
		t.Errorf("Expected 4 results, got %d", len(result))
	}
}

func TestIFFTErrors(t *testing.T) {
	_, err := IFFT([]Complex{})
	if err != ErrEmptyInput {
		t.Errorf("Expected ErrEmptyInput, got %v", err)
	}

	input := []Complex{
		{Real: 1, Imag: 0},
		{Real: 2, Imag: 0},
		{Real: 3, Imag: 0},
	}
	_, err = IFFT(input)
	if err != ErrInvalidFFTSize {
		t.Errorf("Expected ErrInvalidFFTSize, got %v", err)
	}
}

func TestFFTBasic(t *testing.T) {
	input := []Complex{
		{Real: 1, Imag: 0},
		{Real: 0, Imag: 0},
		{Real: 0, Imag: 0},
		{Real: 0, Imag: 0},
	}

	result, err := FFT(input)
	if err != nil {
		t.Fatalf("FFT failed: %v", err)
	}

	for i, c := range result {
		magnitude := c.Magnitude()
		if math.Abs(magnitude-1.0) > epsilon {
			t.Errorf("FFT impulse test bin %d: expected magnitude 1.0, got %.6f", i, magnitude)
		}
	}
}

func TestFFTSinusoid(t *testing.T) {
	sampleRate := 8.0
	freq := 1.0
	n := 8

	input := make([]Complex, n)
	for i := 0; i < n; i++ {
		t := float64(i) / sampleRate
		amplitude := math.Sin(2 * math.Pi * freq * t)
		input[i] = Complex{Real: amplitude, Imag: 0}
	}

	result, err := FFT(input)
	if err != nil {
		t.Fatalf("FFT failed: %v", err)
	}

	expectedBin := 1
	maxMagnitude := 0.0
	maxBin := 0

	for i, c := range result {
		magnitude := c.Magnitude()
		if magnitude > maxMagnitude {
			maxMagnitude = magnitude
			maxBin = i
		}
	}

	if maxBin != expectedBin {
		t.Errorf("Expected peak at bin %d, got peak at bin %d", expectedBin, maxBin)
	}
}

func TestFFTIFFTRoundTrip(t *testing.T) {
	original := []Complex{
		{Real: 1, Imag: 0},
		{Real: 2, Imag: 0},
		{Real: 3, Imag: 0},
		{Real: 4, Imag: 0},
	}

	forward, err := FFT(original)
	if err != nil {
		t.Fatalf("FFT failed: %v", err)
	}

	reconstructed, err := IFFT(forward)
	if err != nil {
		t.Fatalf("IFFT failed: %v", err)
	}

	for i, orig := range original {
		if math.Abs(reconstructed[i].Real-orig.Real) > epsilon ||
			math.Abs(reconstructed[i].Imag-orig.Imag) > epsilon {
			t.Errorf("Round trip failed at index %d: expected (%.6f, %.6f), got (%.6f, %.6f)",
				i, orig.Real, orig.Imag, reconstructed[i].Real, reconstructed[i].Imag)
		}
	}
}

func TestDFTBasic(t *testing.T) {
	samples := []float64{1, 0, 0, 0}
	sampleRate := 4.0

	result, err := DFT(samples, sampleRate)
	if err != nil {
		t.Fatalf("DFT failed: %v", err)
	}

	if result.SampleRate != sampleRate {
		t.Errorf("Expected sample rate %.1f, got %.1f", sampleRate, result.SampleRate)
	}

	if result.SampleCount != len(samples) {
		t.Errorf("Expected sample count %d, got %d", len(samples), result.SampleCount)
	}

	expectedComponents := 2
	if len(result.Components) != expectedComponents {
		t.Errorf("Expected %d components, got %d", expectedComponents, len(result.Components))
	}
}

func TestDFTSinusoid(t *testing.T) {
	sampleRate := 8.0
	freq := 1.0
	duration := 1.0
	numSamples := int(sampleRate * duration)

	samples := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		samples[i] = math.Sin(2 * math.Pi * freq * t)
	}

	result, err := DFT(samples, sampleRate)
	if err != nil {
		t.Fatalf("DFT failed: %v", err)
	}

	maxMagnitude := 0.0
	peakFreq := 0.0

	for _, comp := range result.Components {
		if comp.Magnitude > maxMagnitude {
			maxMagnitude = comp.Magnitude
			peakFreq = comp.Frequency
		}
	}

	if math.Abs(peakFreq-freq) > 0.1 {
		t.Errorf("Expected peak at %.1f Hz, got peak at %.1f Hz", freq, peakFreq)
	}
}

func TestGetTopFrequencies(t *testing.T) {
	components := []FrequencyComponent{
		{Frequency: 100, Magnitude: 0.5},
		{Frequency: 200, Magnitude: 1.0},
		{Frequency: 300, Magnitude: 0.3},
		{Frequency: 400, Magnitude: 0.8},
	}

	result := &DFTResult{
		Components: components,
	}

	top := result.GetTopFrequencies(2)

	if len(top) != 2 {
		t.Fatalf("Expected 2 top frequencies, got %d", len(top))
	}

	if top[0].Frequency != 200 || top[1].Frequency != 400 {
		t.Errorf("Expected frequencies [200, 400], got [%.0f, %.0f]", top[0].Frequency, top[1].Frequency)
	}
}

func TestDFTIDFTRoundTrip(t *testing.T) {
	original := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	sampleRate := 8.0

	dftResult, err := DFT(original, sampleRate)
	if err != nil {
		t.Fatalf("DFT failed: %v", err)
	}

	reconstructed, err := IDFT(dftResult)
	if err != nil {
		t.Fatalf("IDFT failed: %v", err)
	}

	if len(reconstructed) != len(original) {
		t.Fatalf("Expected %d samples, got %d", len(original), len(reconstructed))
	}

	for i, orig := range original {
		if math.Abs(reconstructed[i]-orig) > 1e-10 {
			t.Errorf("Round trip failed at index %d: expected %.6f, got %.6f",
				i, orig, reconstructed[i])
		}
	}
}

func TestWindowFunctions(t *testing.T) {
	samples := []float64{1, 1, 1, 1}

	hamming := ApplyHammingWindow(samples)
	if len(hamming) != len(samples) {
		t.Errorf("Hamming window changed length: expected %d, got %d", len(samples), len(hamming))
	}

	hann := ApplyHannWindow(samples)
	if len(hann) != len(samples) {
		t.Errorf("Hann window changed length: expected %d, got %d", len(samples), len(hann))
	}

	blackman := ApplyBlackmanWindow(samples)
	if len(blackman) != len(samples) {
		t.Errorf("Blackman window changed length: expected %d, got %d", len(samples), len(blackman))
	}

	for _, windowed := range [][]float64{hamming, hann, blackman} {
		if windowed[0] >= samples[0] || windowed[len(windowed)-1] >= samples[len(samples)-1] {
			t.Error("Window function should reduce edge values")
		}
	}

	empty := []float64{}
	emptyHamming := ApplyHammingWindow(empty)
	if len(emptyHamming) != 0 {
		t.Errorf("Expected empty result for empty input, got %d elements", len(emptyHamming))
	}
}

func TestAnalyzeFrequencies(t *testing.T) {
	sampleRate := 1000.0
	duration := 1.0
	numSamples := int(sampleRate * duration)

	samples := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		samples[i] = math.Sin(2*math.Pi*50*t) + 0.5*math.Sin(2*math.Pi*120*t)
	}

	result1, err := AnalyzeFrequencies(samples, sampleRate, false)
	if err != nil {
		t.Fatalf("AnalyzeFrequencies without windowing failed: %v", err)
	}

	result2, err := AnalyzeFrequencies(samples, sampleRate, true)
	if err != nil {
		t.Fatalf("AnalyzeFrequencies with windowing failed: %v", err)
	}

	top1 := result1.GetTopFrequencies(2)
	top2 := result2.GetTopFrequencies(2)

	if len(top1) < 2 || len(top2) < 2 {
		t.Fatal("Should detect at least 2 peak frequencies")
	}

	freqs1 := []float64{top1[0].Frequency, top1[1].Frequency}
	freqs2 := []float64{top2[0].Frequency, top2[1].Frequency}

	for _, freqs := range [][]float64{freqs1, freqs2} {
		found50 := false
		found120 := false

		for _, f := range freqs {
			if math.Abs(f-50) < 5 {
				found50 = true
			}
			if math.Abs(f-120) < 5 {
				found120 = true
			}
		}

		if !found50 || !found120 {
			t.Errorf("Expected to find peaks near 50 Hz and 120 Hz, got frequencies: %.1f, %.1f", freqs[0], freqs[1])
		}
	}
}

func TestEmptyInput(t *testing.T) {
	result, err := DFT([]float64{}, 48000)
	if err != nil {
		t.Fatalf("DFT with empty input failed: %v", err)
	}

	if len(result.Components) != 0 {
		t.Errorf("Expected 0 components for empty input, got %d", len(result.Components))
	}

	samples, err := IDFT(result)
	if err != nil {
		t.Fatalf("IDFT with empty input failed: %v", err)
	}

	if len(samples) != 0 {
		t.Errorf("Expected 0 samples for empty IDFT, got %d", len(samples))
	}
}

func TestDFTErrors(t *testing.T) {
	samples := []float64{1, 2, 3, 4}
	result, err := DFT(samples, 48000)

	if err != nil {
		t.Errorf("DFT should handle padding automatically, got error: %v", err)
	}

	if result == nil {
		t.Error("DFT should return a result even for non-power-of-2 input")
	}
}

func TestIDFTErrors(t *testing.T) {
	emptyResult := &DFTResult{
		Components:  []FrequencyComponent{},
		SampleCount: 0,
	}

	samples, err := IDFT(emptyResult)
	if err != nil {
		t.Fatalf("IDFT with empty components failed: %v", err)
	}

	if len(samples) != 0 {
		t.Errorf("Expected 0 samples for empty IDFT, got %d", len(samples))
	}
}

func BenchmarkFFT1024(b *testing.B) {
	input := make([]Complex, 1024)
	for i := range input {
		input[i] = Complex{Real: float64(i), Imag: 0}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := FFT(input)
		if err != nil {
			b.Fatalf("FFT failed: %v", err)
		}
	}
}

func BenchmarkDFT1024(b *testing.B) {
	samples := make([]float64, 1024)
	for i := range samples {
		samples[i] = math.Sin(2 * math.Pi * float64(i) / 1024)
	}
	sampleRate := 48000.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DFT(samples, sampleRate)
		if err != nil {
			b.Fatalf("DFT failed: %v", err)
		}
	}
}
