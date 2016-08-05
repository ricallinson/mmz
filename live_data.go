package main

import (
	"bytes"
	"strings"
	"strconv"
)

type LiveData struct {
	Timestamp                      int
	RxCtrlFlagByte int
	AverageCurrentOnMotor          int
	AvailableCurrentFromController int
	ArmDC                          int
	BatteryVoltage                 int
	MotorVoltage                   int
	ControllerTemp                 int
	SpiErrorCount                  int
	CurrentError                   string
	OperatingStatus                int
	MotorKilowatts                 int
	StoppedState                   bool
	ShiftingInProgress             bool
	MainContactorIsOn              bool
	MotorContactorsAreOn           bool
	DirectionIsReverse             bool
	DirectionIsForward             bool
	MotorsAreInParallel            bool
	MotorsAreInSeries              bool
	MainContactorHasVoltageDrop    bool
}

func getIntFromHex(s string) int {
	v, _ := strconv.ParseInt(s, 16, 16)
	return int(v)
}

func CreateLiveData(b []byte) *LiveData {
	line := string(b[bytes.Index(b, []byte{10}):])
	values := strings.Split(strings.TrimSpace(line), " ")
	data := &LiveData{
		Timestamp:                      0,
		// RxCtrlFlagByte
		AverageCurrentOnMotor:          getIntFromHex(values[1]),
		AvailableCurrentFromController: getIntFromHex(values[2]),
		// ArmDC
		BatteryVoltage: getIntFromHex(values[4]),
		MotorVoltage:   getIntFromHex(values[5]),
		ControllerTemp: getIntFromHex(values[6]),
		// SpiErrorCount
		CurrentError: strconv.Itoa(int(getIntFromHex(values[8]))),
		// OperatingStatus
		ShiftingInProgress:          false,
		MainContactorIsOn:           true,
		MotorContactorsAreOn:        true,
		DirectionIsReverse:          false,
		DirectionIsForward:          true,
		MotorsAreInParallel:         false,
		MotorsAreInSeries:           false,
		MainContactorHasVoltageDrop: false,
	}
	data.MotorKilowatts = data.MotorVoltage * data.AverageCurrentOnMotor / 1000
	data.CurrentError = data.CurrentError + ": " + Codes[data.CurrentError]
	return data
}
