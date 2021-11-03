package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type LiveData struct {
	Timestamp                      int64
	RxCtrlFlagByte                 int
	DrivePot                       int
	SpeedOne                       int
	AverageCurrentOnMotor          int
	AvailableCurrentFromController int
	ArmDC                          int
	BatteryVoltage                 int
	MotorVoltage                   int
	ControllerTemp                 int
	SpiErrorCount                  int
	CurrentError                   string
	OperatingStatus                int
	MotorKilowatts                 float64
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
	v, _ := strconv.ParseUint(s, 32, 32)
	return int(v)
}

func GetRealtimeValuesQ1(b []byte) *LiveData {
	values := strings.Split(string(b), " ")
	if len(values) < 9 {
		return nil
	}
	ld := &LiveData{
		Timestamp:                      time.Now().Unix(),
		RxCtrlFlagByte:                 0,
		AverageCurrentOnMotor:          getIntFromHex(values[1]),
		AvailableCurrentFromController: getIntFromHex(values[2]),
		ArmDC:                          0,
		BatteryVoltage:                 getIntFromHex(values[4]),
		MotorVoltage:                   getIntFromHex(values[5]),
		ControllerTemp:                 getIntFromHex(values[6]),
		SpiErrorCount:                  0,
		CurrentError:                   strconv.Itoa(int(getIntFromHex(values[8]))),
		OperatingStatus:                0,
	}
	// States
	setStates(values[10], ld)
	ld.MotorKilowatts = float64(ld.MotorVoltage*ld.AverageCurrentOnMotor) / 1000
	ld.CurrentError = ld.CurrentError + ": " + Codes[ld.CurrentError]
	return ld
}

func GetRealtimeValuesQ4(b []byte) *LiveData {
	values := strings.Split(string(b), " ")
	if len(values) < 9 {
		return nil
	}
	ld := &LiveData{
		Timestamp:                      time.Now().Unix(),
		DrivePot:                       getIntFromHex(values[0]),
		SpeedOne:                       getIntFromHex(values[1]),
		AverageCurrentOnMotor:          getIntFromHex(values[2]),
		AvailableCurrentFromController: getIntFromHex(values[3]),
		ArmDC:                          0,
		BatteryVoltage:                 getIntFromHex(values[5]),
		MotorVoltage:                   getIntFromHex(values[6]),
		ControllerTemp:                 getIntFromHex(values[7]),
		OperatingStatus:                0,
	}
	// States
	setStates(values[10], ld)
	ld.MotorKilowatts = float64(ld.MotorVoltage*ld.AverageCurrentOnMotor) / 1000
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
