package main

import core "github.com/averagestardust/eridol-core-go"

type Syllable core.Notes

var BD = &Syllable{B: true, Ds: true}
var BF = &Syllable{B: true, Fs: true}
var BA = &Syllable{B: true, A: true}
var DF = &Syllable{Ds: true, Fs: true}
var DA = &Syllable{Ds: true, A: true}
var FA = &Syllable{Fs: true, A: true}
var BDF = &Syllable{B: true, Ds: true, Fs: true}
var BDA = &Syllable{B: true, Ds: true, A: true}
var BFA = &Syllable{B: true, Fs: true, A: true}
var DFA = &Syllable{Ds: true, Fs: true, A: true}
var BDFA = &Syllable{B: true, Ds: true, Fs: true, A: true}

var syllables = []*Syllable{BD, BF, BA, DF, DA, FA, BDF, BDA, BFA, DFA, BDFA}

func identifySyllable(notes core.Notes) *Syllable {
	for _, syllable := range syllables {
		if (*core.Notes)(syllable).IsEqualWithoutTones(notes) {
			return syllable
		}
	}

	return nil
}
