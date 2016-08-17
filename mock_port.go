package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type MockPort struct {
	history                        byte
	update                         byte
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

func copyIntoArray(s []byte, d []byte) {
	for i, _ := range d {
		if i >= len(s) || i >= len(d) {
			return
		}
		d[i] = s[i]
	}
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
	switch this.history {
	case 'Q':
		return this.q1(b)
	case 'd', 'b', 'm', 's', 'o', 'p', 27:
		copyIntoArray(this.mocks[this.history], b)
		return len(this.mocks[this.history]), nil
	}
	return 0, nil
}

func (this *MockPort) Write(b []byte) (int, error) {
	switch this.history {
	case 'b', 'm', 's', 'o':
		this.update = b[0]
		return 1, nil
	}
	switch b[0] {
	case 'd', 'b', 'm', 's', 'o', 'p', 27, 'Q':
		this.history = b[0]
		this.update = 0
		return 1, nil
	default:
		// Value is always an int.
		// Get the mock and replace;
		// string(this.update) + ") int"
		// with;
		// string(this.update) + ")" + string(b)
		fmt.Println(string(b))
		return len(b), nil
	}
}

func (this *MockPort) Flush() error {
	return nil
}

func (this *MockPort) Close() error {
	return nil
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
