package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type MockPort struct {
	menu                           byte
	update                         string
	averageCurrentOnMotor          int
	availableCurrentFromController int
	armDC                          int
	batteryVoltage                 int
	motorVoltage                   int
	controllerTemp                 int
	spiErrorCount                  int
	currentError                   int
	operatingStatus                int
	mocks                          map[byte][]byte
}

func NewMockPort() *MockPort {
	this := &MockPort{
		mocks: make(map[byte][]byte),
	}
	var file *os.File
	// Display Settings.
	file, _ = os.Open("./fixtures/settings.txt")
	this.mocks['d'], _ = ioutil.ReadAll(file)
	// Battery Menu.
	file, _ = os.Open("./fixtures/battery.txt")
	this.mocks['b'], _ = ioutil.ReadAll(file)
	// Motor Menu.
	file, _ = os.Open("./fixtures/motor.txt")
	this.mocks['m'], _ = ioutil.ReadAll(file)
	// Speed Menu.
	file, _ = os.Open("./fixtures/speed.txt")
	this.mocks['s'], _ = ioutil.ReadAll(file)
	// Options Menu.
	file, _ = os.Open("./fixtures/options.txt")
	this.mocks['o'], _ = ioutil.ReadAll(file)
	// Special Menu.
	file, _ = os.Open("./fixtures/special.txt")
	this.mocks['p'], _ = ioutil.ReadAll(file)
	// Home Menu.
	file, _ = os.Open("./fixtures/home.txt")
	this.mocks[27], _ = ioutil.ReadAll(file)
	// All done.
	return this
}

func (this *MockPort) Read(b []byte) (int, error) {
	switch this.menu {
	case 'Q':
		return this.q1(b)
	case 'd', 'b', 'm', 's', 'o', 'p', 27:
		copyIntoArray(this.mocks[this.menu], b)
		return len(this.mocks[this.menu]), nil
	}
	return 0, nil
}

func (this *MockPort) Write(b []byte) (int, error) {
	// Are we in a menu?
	switch this.menu {
	case 'b', 'm', 's', 'o':
		if b[0] != 27 {
			this.changeSettingsValue(string(b))
			return len(b), nil
		}
	}
	// If not then are we going into a menu?
	switch b[0] {
	case 'd', 'b', 'm', 's', 'o', 'p', 27, 'Q':
		this.menu = b[0]
		this.update = ""
		return 1, nil
	}
	return 0, nil
}

func (this *MockPort) Flush() error {
	return nil
}

func (this *MockPort) Close() error {
	return nil
}

func (this *MockPort) changeSettingsValue(value string) {
	if this.menu == 'o' {
		this.changeSettingsToggleValue(value)
	} else if this.update != "" {
		this.changeSettingsIntValue(string(this.menu), this.update, value)
		this.update = ""
	} else {
		this.update = value
	}
}

func (this *MockPort) changeSettingsIntValue(menu string, option string, value string) {
	b := this.mocks['d']
	var oldLine []byte
	var values []string
	var index int
	switch menu {
	case "b":
		// First find the start and end to change based on the menu.
		oldLine = findValueLine(b, "LBVI")
		// Next find the bytes to change based on start and end and option.
		values = split(string(oldLine), " ")
		switch option {
		case "a":
			index = 0
		case "v":
			index = 1
		case "i":
			index = 2
		}
	case "m":
		// First find the start and end to change based on the menu.
		// Next find the bytes to change based on start and end and option.
		switch option {
		case "a":
			oldLine = findValueLine(b, "RA")
			index = 0
		case "v":
			oldLine = findValueLine(b, "RA")
			index = 1
		case "i":
			oldLine = findValueLine(b, "RA")
			index = 2
		case "r":
			oldLine = findValueLine(b, "PV")
			index = 0
		case "c":
			oldLine = findValueLine(b, "PV")
			index = 1
		case "p":
			oldLine = findValueLine(b, "PV")
			index = 2
		}
		values = split(string(oldLine), " ")
	case "s":
		// First find the start and end to change based on the menu.
		oldLine = findValueLine(b, "Max")
		// Next find the bytes to change based on start and end and option.
		values = split(string(oldLine), " ")
		switch option {
		case "l":
			index = 0
		case "r":
			index = 1
		case "x":
			index = 2
		}
	}
	// Replace the old value with the new value.
	newLine := bytes.Replace(oldLine, []byte(values[index]), []byte(value), 1)
	this.mocks['d'] = bytes.Replace(b, oldLine, newLine, 1)
	fmt.Println(string(this.mocks['d']))
}

func (this *MockPort) changeSettingsToggleValue(option string) {
	b := this.mocks['d']
	if bytes.Contains(b, []byte(option+") On")) {
		b = bytes.Replace(b, []byte(option+") On"), []byte(option+") Off"), 1)
	} else {
		b = bytes.Replace(b, []byte(option+") Off"), []byte(option+") On"), 1)
	}
	this.mocks['d'] = b
}

func (this *MockPort) q1(b []byte) (int, error) {
	time.Sleep(100 * time.Millisecond)
	this.updateMockData()
	averageCurrentOnMotor := fmt.Sprintf("%X", this.averageCurrentOnMotor)
	availableCurrentFromController := fmt.Sprintf("%X", this.availableCurrentFromController)
	armDC := fmt.Sprintf("%X", this.armDC)
	batteryVoltage := fmt.Sprintf("%X", this.batteryVoltage)
	motorVoltage := fmt.Sprintf("%X", this.motorVoltage)
	controllerTemp := fmt.Sprintf("%X", this.controllerTemp)
	spiErrorCount := fmt.Sprintf("%X", this.spiErrorCount)
	currentError := fmt.Sprintf("%X", this.currentError)
	operatingStatus := fmt.Sprintf("%X", this.operatingStatus)
	line := []byte(
		"00 " +
			averageCurrentOnMotor + " " +
			availableCurrentFromController + " " +
			armDC + " " +
			batteryVoltage + " " +
			motorVoltage + " " +
			controllerTemp + " " +
			spiErrorCount + " " +
			currentError + " " +
			operatingStatus + " " +
			"SFS\n")
	for i := 0; i < len(line); i++ {
		b[i] = line[i]
	}
	return len(line), nil
}

func (this *MockPort) updateMockData() {
	this.averageCurrentOnMotor++
	if this.averageCurrentOnMotor > 2000 {
		this.averageCurrentOnMotor = 0
	}
	this.availableCurrentFromController++
	if this.availableCurrentFromController > 2000 {
		this.availableCurrentFromController = 100
	}
	this.armDC++
	if this.armDC > 10 {
		this.armDC = 0
	}
	this.batteryVoltage++
	if this.batteryVoltage < 285 {
		this.batteryVoltage = 285
	}
	if this.batteryVoltage > 290 {
		this.batteryVoltage = 285
	}
	this.motorVoltage++
	if this.motorVoltage > 290 {
		this.motorVoltage = 150
	}
	this.controllerTemp++
	if this.controllerTemp > 120 {
		this.controllerTemp = 80
	}
	this.spiErrorCount++
	if this.spiErrorCount > 1000 {
		this.spiErrorCount = 0
	}
	this.currentError++
	if this.currentError < 1111 {
		this.currentError = 1111
	}
	if this.currentError > 1500 {
		this.currentError = 1111
	}
}
