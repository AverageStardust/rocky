package main

func intPow(n, m uint) uint {
	if m == 0 {
		return 1
	}

	if m == 1 {
		return n
	}

	nPowHalfM := intPow(n, m/2)
	if m%2 == 0 {
		// even
		return nPowHalfM * nPowHalfM
	} else {
		// odd
		return nPowHalfM * nPowHalfM * n
	}
}
