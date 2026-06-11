package main

import (
	"errors"
)

type Phrase interface {
	check(phrases []Phrase) int
	maxLength() int
}

type Parseable[T any] interface {
	Phrase
	parse(phrases []Phrase) (value T, length int)
}

type SingleParseable[T any] struct {
	pattern Phrase
	parser  func(phrases []Phrase) T
}

type AnyParseable[T any] struct {
	pattern []Phrase
	parser  func(phrases []Phrase) T
}

type SeriesParseable[T any] struct {
	pattern []Phrase
	parser  func(phrases []Phrase) T
}

var errUnexpectedParsingState = errors.New("unexpected parsing state")

func (word *Word) check(phrases []Phrase) int {
	if len(phrases) < 1 {
		return 0
	}

	phraseWord, ok := phrases[0].(*Word)

	if ok && phraseWord == word {
		return 1
	}

	return 0
}

func (word *Word) maxLength() int {
	return 1
}

func (single *SingleParseable[T]) check(phrases []Phrase) int {
	return single.pattern.check(phrases)
}

func (single *SingleParseable[T]) maxLength() int {
	return single.pattern.maxLength()
}

func (single *SingleParseable[T]) parse(phrases []Phrase) (value T, length int) {
	length = single.check(phrases)
	if length == 0 {
		return
	}

	return single.parser(phrases[:length]), length
}

func (any *AnyParseable[T]) check(phrases []Phrase) int {
	for _, phrase := range any.pattern {
		length := phrase.check(phrases)

		if length > 0 {
			return length
		}
	}

	return 0
}

func (any *AnyParseable[T]) maxLength() int {
	maxLength := 0

	for _, phrase := range any.pattern {
		maxLength = max(maxLength, phrase.maxLength())
	}

	return maxLength
}

func (any *AnyParseable[T]) parse(phrases []Phrase) (value T, length int) {
	length = any.check(phrases)
	if length == 0 {
		return
	}

	return any.parser(phrases[:length]), length
}

func (series *SeriesParseable[T]) check(phrases []Phrase) int {
	totalLength := 0
	for _, phrase := range series.pattern {
		length := phrase.check(phrases[totalLength:])

		if length == 0 {
			return 0
		}

		totalLength += length
	}

	return totalLength
}

func (series *SeriesParseable[T]) maxLength() int {
	maxLength := 0

	for _, phrase := range series.pattern {
		maxLength += phrase.maxLength()
	}

	return maxLength
}

func (series *SeriesParseable[T]) parse(phrases []Phrase) (value T, length int) {
	length = series.check(phrases)
	if length == 0 {
		return
	}

	return series.parser(phrases[:length]), length
}
