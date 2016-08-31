package main

import (
	"strconv"
)

type zillaCommand struct {
	bytes [][]byte
	done  chan bool
	data  []byte
}

func newZillaCommand() *zillaCommand {
	this := &zillaCommand{
		bytes: [][]byte{},
		done:  make(chan bool),
	}
	this.sendHome()
	return this
}

func (this *zillaCommand) sendInt(val int) {
	this.sendString(strconv.Itoa(val) + "\r\n")
}

func (this *zillaCommand) sendHome() {
	this.sendBytes([]byte{27})
	this.sendBytes([]byte{27})
	this.sendBytes([]byte{27})
}

func (this *zillaCommand) sendString(s string) {
	this.sendBytes([]byte(s))
}

func (this *zillaCommand) sendBytes(b []byte) {
	this.bytes = append(this.bytes, []byte(b))
}
