package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type LiveData struct {
	Timestamp                      int
	RxCtrlFlagByte                 int
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

func (this *LiveData) ToBytes() []byte {
	b, _ := json.Marshal(this)
	return append(b, 10)
}

func CreateLiveData(b []byte) *LiveData {
	var d LiveData
	json.Unmarshal(b, &d)
	return &d
}

func getIntFromHex(s string) int {
	v, _ := strconv.ParseInt(s, 16, 16)
	return int(v)
}

func ParseQ1LineFromHairball(b []byte) *LiveData {
	line := string(b)
	values := strings.Split(strings.TrimSpace(line), " ")
	ld := &LiveData{
		Timestamp: 0,
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
	}
	// States
	setStates(values[10], ld)
	ld.MotorKilowatts = ld.MotorVoltage * ld.AverageCurrentOnMotor / 1000
	ld.CurrentError = ld.CurrentError + ": " + Codes[ld.CurrentError]
	return ld
}

func setStates(s string, ld *LiveData) {
	for i, v := range s {
		switch v {
		case 'S':
			if i == 0 {
				ld.StoppedState = true
			} else {
				ld.MotorsAreInSeries = true
			}
		case 'G':
			ld.ShiftingInProgress = true
		case 'O':
			ld.MainContactorIsOn = true
		case 'M':
			ld.MotorContactorsAreOn = true
		case 'R':
			ld.DirectionIsReverse = true
		case 'F':
			ld.DirectionIsForward = true
		case 'P':
			ld.MotorsAreInParallel = true
		case 'V':
			ld.MainContactorHasVoltageDrop = true
		}
	}
}
