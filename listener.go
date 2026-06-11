package main

import core "github.com/averagestardust/eridol-core-go"

type Listener struct {
	octave           int
	owned            bool
	syllables        []*Syllable
	phrases          []Phrase
	onClaimContested func()
	onStop           func()
}

type ListenerState int

func newListener(octave int) *Listener {
	return &Listener{
		octave:    octave,
		syllables: []*Syllable{},
		phrases:   []Phrase{},
	}
}

// Marks octave as owned, claim will be protected as needed.
// Returns true if transition is allowed.
func (listener *Listener) Own() bool {
	if listener.owned {
		return false
	}

	listener.owned = true
	listener.syllables = []*Syllable{}
	listener.phrases = []Phrase{}
	listener.DismisOnClaimContested()

	return true
}

// Marks octave as disowned, claim will no longer be protected.
// Returns true if transition is allowed.
func (listener *Listener) Disown() bool {
	if !listener.owned {
		return false
	}

	synth := core.Synth(listener.octave)
	synth.StopNoteImmediately(core.CounterClaim)

	listener.owned = false

	return true
}

// Runs callback when counter-claim tone is heard, overwrites any previous callbacks.
// Waits until tone is heard, callback is dismissed, or until ownership is taken.
func (listener *Listener) OnClaimContested(callback func()) {
	listener.onClaimContested = callback
}

// Dismisses the callback listening for a counter-claim.
func (listener *Listener) DismisOnClaimContested() {
	listener.onClaimContested = nil
}

func (listener *Listener) RunLoop(heard core.Heard) {
	heardOctave := heard.Octaves[listener.octave]

	if !listener.owned {
		if heardOctave.CounterClaim && listener.onClaimContested != nil {
			listener.onClaimContested()
			listener.DismisOnClaimContested()
		}

		executer := listener.processNotes(heardOctave)
		if executer != nil {
			executer()
		}
	} else {
		synth := core.Synth(listener.octave)
		synth.SetNoteImmediately(core.CounterClaim, heardOctave.Claim)
	}
}

func (listener *Listener) processNotes(notes core.Notes) func() {
	syllable := identifySyllable(notes)

	if syllable != nil {
		if len(listener.syllables) == 0 {
			listener.syllables = append(listener.syllables, syllable)
		} else {
			lastSyllable := listener.syllables[len(listener.syllables)-1]
			if lastSyllable != syllable {
				listener.syllables = append(listener.syllables, syllable)
			}
		}
	}

	listener.formWords(false)
	if len(listener.syllables) >= longestWord {
		listener.syllables = listener.syllables[1:]
		listener.formWords(true)
	}

	executer := listener.formCommands()
	if len(listener.phrases) >= longestCommand {
		listener.phrases = listener.phrases[1:]
		executer = listener.formCommands()
	}

	return executer
}

func (listener *Listener) formWords(identifySubsets bool) (success bool) {
	word := identifyWord(listener.syllables, identifySubsets)

	if word != nil {
		// discard stops
		if word != STOP {
			listener.phrases = append(listener.phrases, word)
		}

		listener.syllables = listener.syllables[len(*word):]
		success = true
	}

	return
}

func (listener *Listener) formCommands() func() {
	executer, length := identifyCommand(listener.phrases)

	if executer != nil {
		listener.phrases = listener.phrases[length:]
	}

	return executer
}
