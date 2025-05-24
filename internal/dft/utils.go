package dft

func nextPowerOf2(n int) int {
	if n <= 1 {
		return 1
	}

	power := 1
	for power < n {
		power <<= 1
	}
	return power
}

func isPowerOf2(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

func reverseBits(num, numBits int) int {
	result := 0
	for i := 0; i < numBits; i++ {
		if (num & (1 << i)) != 0 {
			result |= 1 << (numBits - 1 - i)
		}
	}
	return result
}
