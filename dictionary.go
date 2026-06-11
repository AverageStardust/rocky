package main

var ZERO = addWord(BD)
var ONE = addWord(BF)
var TWO = addWord(BA)
var THREE = addWord(DF)
var FOUR = addWord(DA)
var FIVE = addWord(FA)
var STOP = addWord(BDFA)
var SEND_NUMBER = addWord(BDF, BD)
var PING = addWord(BFA, BD)
var PONG = addWord(BFA, FA)

var DIGIT = &AnyParseable[uint]{
	pattern: []Phrase{ZERO, ONE, TWO, THREE, FOUR, FIVE},
	parser:  parseHexalDigit,
}
var NUMBER = &SeriesParseable[uint]{
	pattern: []Phrase{DIGIT, DIGIT, DIGIT, DIGIT, DIGIT, DIGIT},
	parser:  parseHexal,
}
