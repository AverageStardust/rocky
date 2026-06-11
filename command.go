package main

var commands = []Parseable[func()]{}
var longestCommand int

func init() {
	addCommand(&SingleParseable[func()]{
		pattern: PING,
		parser: func(_ []Phrase) func() {
			return func() {
				speaker.PlayPongIfQuite()
			}
		},
	})

	addCommand(&SingleParseable[func()]{
		pattern: PONG,
		parser: func(_ []Phrase) func() {
			return nil
		},
	})

	addCommand(&SeriesParseable[func()]{
		pattern: []Phrase{SEND_NUMBER, NUMBER},
		parser: func(phrases []Phrase) func() {
			target, _ := DIGIT.parse(phrases[1:])
			number, _ := NUMBER.parse(phrases[1:])
			return func() {
				if speaker.octave == int(target) {
					println("Received:", number)
				}
			}
		},
	})
}

func addCommand(command Parseable[func()]) {
	commands = append(commands, command)
	longestCommand = max(longestCommand, command.maxLength())
}

func identifyCommand(phrases []Phrase) (executer func(), length int) {
	for _, command := range commands {
		executer, length = command.parse(phrases)
		if length > 0 {
			return
		}
	}

	return nil, 0
}
