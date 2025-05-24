package analyze

type FreqMag struct {
	Freq float64
	Mag  float64
}

type AnalysisResult struct {
	TopFrequencies []FreqMag
	SampleRate     float64
	SampleCount    int
	WindowUsed     bool
}
