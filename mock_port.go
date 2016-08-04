package main

import (
	"os"
	"time"
	"fmt"
)

type MockPort struct {
	history byte
	averageCurrentOnMotor int
	availableCurrentFromController int
	armDC int
	batteryVoltage int
	motorVoltage int
	controllerTemp int
	spiErrorCount int
	currentError int
	operatingStatus string
}

func (this *MockPort) Read(b []byte) (int, error) {
	var file *os.File
	var err error
	switch this.history {
	case 'd': // Display Settings
		file, err = os.Open("./fixtures/settings.txt")
	case 'b': // Battery Menu
		file, err = os.Open("./fixtures/battery.txt")
	case 'm': // Motor Menu
		file, err = os.Open("./fixtures/motor.txt")
	case 's': // Speed Menu
		file, err = os.Open("./fixtures/speed.txt")
	case 'o': // Options Menu
		file, err = os.Open("./fixtures/options.txt")
	case 'p': // Special Menu
		file, err = os.Open("./fixtures/special.txt")
	case 0, 27: // Home Menu
		file, err = os.Open("./fixtures/home.txt")
	case 'Q':
		return this.q1(b)
	}
	if err != nil {
		return 0, err
	}
	return file.Read(b)
}

func (this *MockPort) Write(b []byte) (int, error) {
	switch b[0] {
	case 'd', 'b', 'm', 's', 'o', 'p', 27:
		this.history = b[0]
	case 'Q':
		this.history = b[0]
	}
	return 0, nil
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
	averageCurrentOnMotor := fmt.Sprintf("%x", this.averageCurrentOnMotor)
	availableCurrentFromController := fmt.Sprintf("%x", this.availableCurrentFromController)
	armDC := fmt.Sprintf("%x", this.armDC)
	batteryVoltage := fmt.Sprintf("%x", this.batteryVoltage)
	motorVoltage := fmt.Sprintf("%x", this.motorVoltage)
	controllerTemp := fmt.Sprintf("%x", this.controllerTemp)
	spiErrorCount := fmt.Sprintf("%x", this.spiErrorCount)
	currentError := fmt.Sprintf("%x", this.currentError)
	line := []byte(
		averageCurrentOnMotor + " " +
		availableCurrentFromController + " " +
		armDC + " " +
		batteryVoltage+ " " +
		motorVoltage + " " +
		controllerTemp + " " +
		spiErrorCount + " " +
		currentError + " " +
		this.operatingStatus +
		"\n")
	for i := 0; i < len(line); i++ {
		b[i] = line[i]
	}
	return len(line), nil
}

func (this *MockPort) updateMockData() {
	this.averageCurrentOnMotor++
	if this.averageCurrentOnMotor > 2000{
		this.averageCurrentOnMotor = 0;
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
	if this.batteryVoltage > 285 {
		this.batteryVoltage = 285
	}
	if this.batteryVoltage > 290 {
		this.batteryVoltage = 285
	}
	this.motorVoltage++
	if this.motorVoltage < 285 {
		this.motorVoltage = 285
	}
	if this.motorVoltage > 290 {
		this.motorVoltage = 285
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
	this.operatingStatus = "S"
}
