package main

type Word []*Syllable

var words = []*Word{}
var longestWord int

func addWord(syllables ...*Syllable) *Word {
	word := (*Word)(&syllables)

	words = append(words, word)
	longestWord = max(longestWord, len(syllables))

	return word
}

func identifyWord(syllables []*Syllable, allowSubsets bool) *Word {
wordLoop:
	for _, word := range words {
		if allowSubsets {
			if len(*word) > len(syllables) {
				continue
			}
		} else {
			if len(*word) != len(syllables) {
				continue
			}
		}

		for i, wordSyllable := range *word {
			if wordSyllable != syllables[i] {
				continue wordLoop
			}
		}

		return word
	}

	return nil
}
