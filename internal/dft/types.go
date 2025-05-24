package dft

import "errors"

var (
	ErrInvalidFFTSize = errors.New("FFT input length must be a power of 2")
	ErrEmptyInput     = errors.New("input cannot be empty")
)

type FrequencyComponent struct {
	Frequency float64
	Magnitude float64
	Phase     float64
	Real      float64
	Imag      float64
}
