package main

import core "github.com/averagestardust/eridol-core-go"

func main() {
	core.Run(func(heard core.Heard) (stop bool) {
		return false
	})
}
