package main

import (
	core "github.com/averagestardust/eridol-core-go"

	input "github.com/AverageStardust/simple-io/input"
)

var listeners [core.OctaveCount]*Listener
var speaker *Speaker

func init() {
	for i := range core.OctaveCount {
		listeners[i] = newListener(i)
	}

	speaker = newSpeaker()
}

func main() {
	stoppingCore := false
	doneClaiming := make(chan struct{})
	doneStopping := make(chan struct{})

	go core.Run(func(heard core.Heard) bool {
		for i := range core.OctaveCount {
			listeners[i].RunLoop(heard)
		}

		if stoppingCore {
			doneStopping <- struct{}{}
		}

		return stoppingCore
	})

	speaker.Claim(func() {
		doneClaiming <- struct{}{}
	})
	<-doneClaiming

	for !uiLoop() {
	}

	speaker.Stop(func() {
		doneStopping <- struct{}{}
	})

	<-doneStopping

	stoppingCore = true

	<-doneStopping
}

func uiLoop() (stop bool) {
	t, err := input.Integer("Target?")
	if err != nil {
		println(err.Error())
		return false
	}

	if t < 0 && t >= 6 {
		println("Target out of range 0-5")
		return false
	}

	n, err := input.Integer("Number?")
	if err != nil {
		println(err.Error())
		return false
	}

	if n < 0 && n >= 46656 { // 6 ** 6
		println("Target out of range 0-46655")
		return false
	}

	speaker.PlaySendNumber(uint(t), uint(n))
	return false
}
