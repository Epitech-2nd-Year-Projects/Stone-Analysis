package dft

type DFTResult struct {
	Components     []FrequencyComponent
	Nyquist        Complex
	SampleRate     float64
	SampleCount    int
	FreqResolution float64
}

func (r *DFTResult) GetTopFrequencies(n int) []FrequencyComponent {
	if n <= 0 || len(r.Components) == 0 {
		return []FrequencyComponent{}
	}

	components := make([]FrequencyComponent, len(r.Components))
	copy(components, r.Components)

	for i := 0; i < len(components)-1; i++ {
		for j := 0; j < len(components)-i-1; j++ {
			if components[j].Magnitude < components[j+1].Magnitude {
				components[j], components[j+1] = components[j+1], components[j]
			}
		}
	}

	if n > len(components) {
		n = len(components)
	}
	return components[:n]
}
