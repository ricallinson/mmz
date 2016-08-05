package main

import(
	// "fmt"
	"bytes"
	"strings"
	// "encoding/hex"
	// "encoding/binary"
	"strconv"
)

type LiveData struct {
	Timestamp                      int
	MotorKilowatts                 int
	MotorVoltage                   int
	BatteryVoltage                 int
	AverageCurrentOnMotor          int
	AvailableCurrentFromController int
	ControllerTemp                 int
	StoppedState                   bool
	ShiftingInProgress             bool
	MainContactorIsOn              bool
	MotorContactorsAreOn           bool
	DirectionIsReverse             bool
	DirectionIsForward             bool
	MotorsAreInParallel            bool
	MotorsAreInSeries              bool
	MainContactorHasVoltageDrop    bool
	LatestOperation                string
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
		AverageCurrentOnMotor:          getIntFromHex(values[0]),
		AvailableCurrentFromController: getIntFromHex(values[1]),
		// ArmDC
		BatteryVoltage:                 getIntFromHex(values[3]),
		MotorVoltage:                   getIntFromHex(values[4]),
		ControllerTemp:                 getIntFromHex(values[5]),
		// SpiErrorCount
		LatestOperation:                strconv.Itoa(int(getIntFromHex(values[7]))),
		ShiftingInProgress:             false,
		MainContactorIsOn:              true,
		MotorContactorsAreOn:           true,
		DirectionIsReverse:             false,
		DirectionIsForward:             true,
		MotorsAreInParallel:            false,
		MotorsAreInSeries:              false,
		MainContactorHasVoltageDrop:    false,
	}
	data.MotorKilowatts = data.MotorVoltage * data.AverageCurrentOnMotor / 1000
	data.LatestOperation = data.LatestOperation + ": " + Codes[data.LatestOperation]
	return data
}
