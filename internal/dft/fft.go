package dft

import "math"

func FFT(input []Complex) ([]Complex, error) {
	n := len(input)

	if n == 0 {
		return nil, ErrEmptyInput
	}

	if n == 1 {
		return input, nil
	}

	if !isPowerOf2(n) {
		return nil, ErrInvalidFFTSize
	}

	output := make([]Complex, n)
	numBits := int(math.Log2(float64(n)))

	for i := 0; i < n; i++ {
		j := reverseBits(i, numBits)
		output[j] = input[i]
	}

	for size := 2; size <= n; size <<= 1 {
		halfSize := size / 2
		step := 2 * math.Pi / float64(size)

		for i := 0; i < n; i += size {
			for j := 0; j < halfSize; j++ {
				u := output[i+j]
				angle := -step * float64(j)
				twiddle := Complex{
					Real: math.Cos(angle),
					Imag: math.Sin(angle),
				}
				v := output[i+j+halfSize].Mul(twiddle)

				output[i+j] = u.Add(v)
				output[i+j+halfSize] = u.Sub(v)
			}
		}
	}

	return output, nil
}

func IFFT(input []Complex) ([]Complex, error) {
	n := len(input)

	if n == 0 {
		return nil, ErrEmptyInput
	}

	if n == 1 {
		return input, nil
	}

	if !isPowerOf2(n) {
		return nil, ErrInvalidFFTSize
	}

	conjugated := make([]Complex, n)
	for i, c := range input {
		conjugated[i] = Complex{Real: c.Real, Imag: -c.Imag}
	}

	result, err := FFT(conjugated)
	if err != nil {
		return nil, err
	}

	for i := range result {
		result[i] = Complex{
			Real: result[i].Real / float64(n),
			Imag: -result[i].Imag / float64(n),
		}
	}

	return result, nil
}
