package main

import (
	"encoding/json"
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
	data := &LiveData{
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
