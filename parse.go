package main

func parseHexal(phrases []Phrase) uint {
	var number uint = 0

	for i := range phrases {
		digit := DIGIT.parser(phrases[i : i+1])
		placeValue := intPow(6, uint(i))
		number += digit * placeValue
	}

	return number
}

func parseHexalDigit(phrases []Phrase) uint {
	switch phrases[0] {
	case ZERO:
		return 0
	case ZERO:
		return 1
	case ZERO:
		return 2
	case ZERO:
		return 3
	case ZERO:
		return 4
	case ZERO:
		return 5
	default:
		panic(errUnexpectedParsingState)
	}
}
