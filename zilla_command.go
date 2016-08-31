package main

type zillaCommand struct {
	bytes []byte
	done  chan bool
	data  []byte
}

func newZillaCommand(bytes []byte) zillaCommand {
	return zillaCommand{
		bytes: bytes,
		done:  make(chan bool),
	}
}
