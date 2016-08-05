package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	LOGFILE = "./log.dat"
)

type SerialPort interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Flush() error
	Close() error
}

type Zilla struct {
	BatteryAmpLimit               int      // a)BA
	LowBatteryVoltageLimit        int      // v)LBV
	LowBatteryVoltageIndicator    int      // i)LBVI
	NormalMotorAmpLimit           int      // a) Amp
	SeriesMotorVoltageLimit       int      // v) Volt
	ReverseMotorAmpLimit          int      // i) RA
	ReverseMotorVoltageLimit      int      // r) RV
	ParallelMotorAmpLimit         int      // c) PA
	ParallelMotorVoltageLimit     int      // p) PV
	ForwardRpmLimit               int      // l)Norm
	ReverseRpmLimit               int      // r)Rev
	MaxRpmLimit                   int      // x)Max
	RpmSensorMotorOne             bool     // a) On
	RpmSensorMotorTwo             bool     // b) On
	AutoShiftingSeriesToParallel  bool     // c) On
	StallDetectOn                 bool     // d) On
	BatteryLightPolarity          bool     // e) Off
	CheckEngineLightPolarity      bool     // f) On
	ReversingContactors           bool     // g) On
	SeriesParallelContactors      bool     // h) On
	ForceParallelInReverse        bool     // i) Off
	InhibitSeriesParallelShifting bool     // j) Off
	TachometerDisplayMotorAmps    bool     // k) Off
	TachometerSixCylinders        bool     // l) Off
	ReversesPlugInInputPolarity   bool     // m) Off
	ActivateHEPI                  bool     // n) Off
	notUsed                       bool     // o) Off
	IsZ2k                         bool     // p) Off
	CurrentState                  string   // 1311
	Errors                        []string // 1111, 1111, ...
	buffer                        []byte   // byte array of the last Zilla output
	serialPort                    SerialPort
	writeLog                      bool
	readLogFile                   *os.File
	writeLogFile                  *os.File
}

// Return a boolean of "On" or "Off".
func truthy(s string) bool {
	return strings.Contains(s, "On")
}

// Returns a string array of values with white space removed.
func split(s string, sep string) []string {
	values := []string{}
	tokens := strings.Split(s, sep)
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if len(token) > 0 {
			values = append(values, token)
		}
	}
	return values
}

func CreateZilla(p SerialPort) (error, *Zilla) {
	z := &Zilla{serialPort: p}
	z.Errors = make([]string, 0)
	if z.Refresh() == false {
		return errors.New("Could not get data from Hairball."), nil
	}
	var openFileError error
	z.writeLogFile, openFileError = os.OpenFile(LOGFILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if openFileError != nil {
		fmt.Println(openFileError)
		return errors.New("Could not open log file for writing."), nil
	}
	z.readLogFile, openFileError = os.Open(LOGFILE)
	if openFileError != nil {
		fmt.Println(openFileError)
		return errors.New("Could not open log file for reading."), nil
	}
	// z.startLogging()
	return nil, z
}

func (this *Zilla) sendString(s, check string) bool {
	return this.sendBytes([]byte(s), check)
}

func (this *Zilla) sendBytes(b []byte, check string) bool {
	var e error
	_, e = this.serialPort.Write(b)
	if e != nil {
		fmt.Println(e)
		return false
	}
	// Cannot keep sleeping. Need a better solution here.
	time.Sleep(500 * time.Millisecond)
	this.buffer = make([]byte, 512)
	_, e = this.serialPort.Read(this.buffer)
	if e != nil {
		fmt.Println(e)
		return false
	}
	// Debug hairball output.
	this.buffer = bytes.TrimSpace(this.buffer)
	// fmt.Println(string(this.buffer))
	return bytes.Index(this.buffer, []byte(check)) > -1
}

func (this *Zilla) startLogging() {
	// Start logging in a go routine that stops and starts around other function calls into Zilla.
	this.writeLog = true
	go this.log()
}

func (this *Zilla) stopLogging() {
	this.writeLog = false
}

func (this *Zilla) appendToLog(line []byte) []byte {
	return line
}

func (this *Zilla) log() {
	this.menuSpecial()
	_, readError := this.serialPort.Write([]byte("Q1\n")) // The first logging screen.
	if readError != nil {
		fmt.Println(readError)
		return
	}
	// Create a reader buffer.
	buf := bufio.NewReader(this.serialPort)
	// While allowed keep writing bytes to the file.
	for this.writeLog {
		line, readLineError := buf.ReadBytes('\n')
		if readLineError != nil {
			fmt.Println(readLineError)
			return
		}
		_, writeLineError := this.writeLogFile.Write(this.appendToLog(line))
		if writeLineError != nil {
			fmt.Println(writeLineError)
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func (this *Zilla) menuHome() bool {
	this.sendBytes([]byte{27}, "")
	this.sendBytes([]byte{27}, "")
	return this.sendBytes([]byte{27}, "d) Display settings")
}

func (this *Zilla) menuSettings() bool {
	this.menuHome()
	return this.sendString("d", "Display only, change with menu")
}

func (this *Zilla) menuBattery() bool {
	this.menuHome()
	return this.sendString("b", "a)BA, v)LBV, i)LBVI")
}

func (this *Zilla) menuMotor() bool {
	this.menuHome()
	return this.sendString("m", "Motor Settings")
}

func (this *Zilla) menuSpeed() bool {
	this.menuHome()
	return this.sendString("s", "Rev limits")
}

func (this *Zilla) menuOptions() bool {
	this.menuHome()
	return this.sendString("o", "Options: Enter letter to change")
}

func (this *Zilla) menuSpecial() bool {
	this.menuHome()
	return this.sendString("p", "Special Menu:")
}

func (this *Zilla) writeIntValue(id string, val int) bool {
	if this.sendString(id, "XXX") && this.sendString(strconv.Itoa(val), strconv.Itoa(val)) {
		return this.Refresh()
	}
	return false
}

func (this *Zilla) writeToggleValue(id string) bool {
	if this.sendString(id, "XXX") {
		return this.Refresh()
	}
	return false
}

func (this *Zilla) GetLiveData() *LiveData {
	bufSize := 50
	buf := make([]byte, bufSize)
	stat, _ := os.Stat(LOGFILE)
	start := stat.Size() - int64(bufSize)
	i, err := this.readLogFile.ReadAt(buf, start)
	if err != nil {
		fmt.Println("Could not read last line from live data.")
		fmt.Println(LOGFILE)
		fmt.Println(start)
		fmt.Println(err)
	}
	return CreateLiveData(buf[:i])
}

// Refreshes all attributes by reading them from the Zilla Controller.
func (this *Zilla) Refresh() bool {
	if this.menuSettings() == false {
		return false
	}
	// Read all the settings in this struct.
	lines := bytes.Split(this.buffer, []byte{10})
	// Get values for BA, LBV, LBVI
	var values []string
	values = split(string(lines[2]), " ")
	this.BatteryAmpLimit, _ = strconv.Atoi(values[0])
	this.LowBatteryVoltageLimit, _ = strconv.Atoi(values[1])
	this.LowBatteryVoltageIndicator, _ = strconv.Atoi(values[2])
	// Values for Amp, Volt, RA
	values = split(string(lines[4]), " ")
	this.NormalMotorAmpLimit, _ = strconv.Atoi(values[0])
	this.SeriesMotorVoltageLimit, _ = strconv.Atoi(values[1])
	this.ReverseMotorAmpLimit, _ = strconv.Atoi(values[2])
	// Values for RV, PA, PV
	values = split(string(lines[6]), " ")
	this.ReverseMotorVoltageLimit, _ = strconv.Atoi(values[0])
	this.ParallelMotorAmpLimit, _ = strconv.Atoi(values[1])
	this.ParallelMotorVoltageLimit, _ = strconv.Atoi(values[2])
	// Values for Norm, Rev, Max
	values = split(string(lines[8]), " ")
	this.ForwardRpmLimit, _ = strconv.Atoi(values[0])
	this.ReverseRpmLimit, _ = strconv.Atoi(values[1])
	this.MaxRpmLimit, _ = strconv.Atoi(values[2])
	// Values for a, b, c, d
	values = split(string(lines[9]), "\t")
	this.RpmSensorMotorOne = truthy(values[0])
	this.RpmSensorMotorTwo = truthy(values[1])
	this.AutoShiftingSeriesToParallel = truthy(values[2])
	this.StallDetectOn = truthy(values[3])
	// Values for e, f, g, h
	values = split(string(lines[10]), "\t")
	this.BatteryLightPolarity = truthy(values[0])
	this.CheckEngineLightPolarity = truthy(values[1])
	this.ReversingContactors = truthy(values[2])
	this.SeriesParallelContactors = truthy(values[3])
	// Values for i, j, k, l
	values = split(string(lines[11]), "\t")
	this.ForceParallelInReverse = truthy(values[0])
	this.InhibitSeriesParallelShifting = truthy(values[1])
	this.TachometerDisplayMotorAmps = truthy(values[2])
	this.TachometerSixCylinders = truthy(values[3])
	// Values for m, n, o, p
	values = split(string(lines[12]), "\t")
	this.ReversesPlugInInputPolarity = truthy(values[0])
	this.ActivateHEPI = truthy(values[1])
	this.notUsed = truthy(values[2])
	this.IsZ2k = truthy(values[3])
	// Values for errors
	this.Errors = split(string(lines[14]), " ")
	// Values for current state
	values = split(string(lines[15]), " ")
	this.CurrentState = values[1]
	return true
}

func (this *Zilla) SetBatteryAmpLimit(val int) bool {
	this.menuBattery()
	this.writeIntValue("a", val)
	return this.BatteryAmpLimit == val
}

func (this *Zilla) SetLowBatteryVoltageLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetNormalMotorAmpLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetSeriesMotorVoltageLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetReverseMotorAmpLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetReverseMotorVoltageLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetParallelMotorAmpLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetParallelMotorVoltageLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetForwardRpmLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetReverseRpmLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetMaxRpmLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetRpmSensorMotorOne(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetRpmSensorMotorTwo(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetAutoShiftingSeriesToParallel(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetStallDetectOn(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetBatteryLightPolarity(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetCheckEngineLightPolarity(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetReversingContactors(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetSeriesParallelContactors(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetForceParallelInReverse(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetInhibitSeriesParallelShifting(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetTachometerDisplayMotorAmps(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetTachometerSixCylinders(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetReversesPlugInInputPolarity(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetActivateHEPI(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetIsZ2k(val bool) bool {
	return this.Refresh()
}
