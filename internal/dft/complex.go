package dft

import "math"

type Complex struct {
	Real float64
	Imag float64
}

func (c Complex) Add(other Complex) Complex {
	return Complex{
		Real: c.Real + other.Real,
		Imag: c.Imag + other.Imag,
	}
}

func (c Complex) Sub(other Complex) Complex {
	return Complex{
		Real: c.Real - other.Real,
		Imag: c.Imag - other.Imag,
	}
}

func (c Complex) Mul(other Complex) Complex {
	return Complex{
		Real: c.Real*other.Real - c.Imag*other.Imag,
		Imag: c.Real*other.Imag + c.Imag*other.Real,
	}
}

func (c Complex) Magnitude() float64 {
	return math.Sqrt(c.Real*c.Real + c.Imag*c.Imag)
}

func (c Complex) Phase() float64 {
	return math.Atan2(c.Imag, c.Real)
}
