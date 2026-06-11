package main

import (
	"time"

	core "github.com/averagestardust/eridol-core-go"
)

type Speaker struct {
	state        SpeakerState
	octave       int
	period       time.Duration
	lastSyllable *Syllable
}

type SpeakerState int

const (
	speakerStopped SpeakerState = iota
	speakerClaiming
	speakerRunning
	speakerStopping
)

const avoidedLowerOctaves = 2

func newSpeaker() *Speaker {
	return &Speaker{
		period: time.Millisecond * 250,
	}
}

// Starts speaker begin claiming.
// Returns true if transition is allowed.
func (speaker *Speaker) Claim(callback func()) bool {
	if speaker.state != speakerStopped {
		return false
	}

	speaker.state = speakerClaiming
	speaker.lastSyllable = nil
	speaker.startClaimingAttempts(callback)

	return true
}

// Makes listener begin stopping.
// Returns true if transition is allowed.
func (speaker *Speaker) Stop(callback func()) bool {
	if speaker.state != speakerRunning {
		return false
	}

	speaker.state = speakerStopping

	core.Synth(speaker.octave).OnAllDone(func() {
		listeners[speaker.octave].Disown()
		speaker.state = speakerStopped
		callback()
	})

	return true
}

func (speaker *Speaker) IsRunning() bool {
	return speaker.state == speakerRunning
}

func (speaker *Speaker) startClaimingAttempts(callback func()) {
	speaker.claim([]int{2, 3, 4, 5, 1, 0}, callback)
}

func (speaker *Speaker) claim(attemptOctaves []int, callback func()) {
	octave := attemptOctaves[0]

	core.Synth(octave).PlanNotes(core.Notes{Claim: true}, 500*time.Millisecond)
	handle := core.Synth(octave).PlanDelay(500 * time.Millisecond)

	listeners[octave].OnClaimContested(func() {
		// failed
		if len(attemptOctaves) == 1 {
			// out of attemps, restart attempts
			speaker.startClaimingAttempts(callback)
		} else {
			speaker.claim(attemptOctaves[1:], callback)
		}
	})

	handle.OnDone(func() {
		// success
		speaker.octave = octave
		speaker.state = speakerRunning
		listeners[speaker.octave].Own()
		callback()
	})
}

func (speaker *Speaker) PlayPing() {
	if !speaker.IsRunning() {
		return
	}

	speaker.playWord(PING)
}

func (speaker *Speaker) PlayPongIfQuite() {
	if !speaker.IsRunning() {
		return
	}

	if core.Synth(speaker.octave).PlanTimeRemaining() < speaker.period {
		speaker.playWord(PONG)
	}
}

func (speaker *Speaker) PlaySendNumber(target uint, number uint) {
	if !speaker.IsRunning() {
		return
	}

	speaker.playWord(SEND_NUMBER)
	speaker.playDigit(target)
	speaker.playNumber(number)
}

func (speaker *Speaker) playNumber(number uint) {
	speaker.playHexal(number, 6)
}

func (speaker *Speaker) playHexal(number, length uint) {
	number %= intPow(6, length)

	for i := range length {
		placeValue := intPow(6, length-1-i)
		digit := number / placeValue
		number -= digit * placeValue
		speaker.playDigit(digit)
	}
}

func (speaker *Speaker) playDigit(n uint) {
	if !speaker.IsRunning() {
		return
	}

	switch n {
	case 0:
		speaker.playWord(ZERO)
	case 1:
		speaker.playWord(ONE)
	case 2:
		speaker.playWord(TWO)
	case 3:
		speaker.playWord(THREE)
	case 4:
		speaker.playWord(FOUR)
	case 5:
		speaker.playWord(FIVE)
	}
}

func (speaker *Speaker) playWord(word *Word) {
	for _, syllable := range *word {
		speaker.playSyllable(syllable)
	}
}

func (speaker *Speaker) playSyllable(syllable *Syllable) {
	synth := core.Synth(speaker.octave)

	if speaker.lastSyllable == syllable {
		synth.PlanNotes(core.Notes(*(*STOP)[0]), speaker.period)
	}

	synth.PlanNotes(core.Notes(*syllable), speaker.period)
}
