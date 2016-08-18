package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
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

func CreateZilla(p SerialPort) (*Zilla, error) {
	this := &Zilla{serialPort: p}
	this.Errors = make([]string, 0)
	// Open log file for reading and writing.
	var openFileError error
	this.writeLogFile, openFileError = os.OpenFile(LOGFILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if openFileError != nil {
		fmt.Println(openFileError)
		return nil, errors.New("Could not open log file for writing.")
	}
	this.readLogFile, openFileError = os.Open(LOGFILE)
	if openFileError != nil {
		fmt.Println(openFileError)
		return nil, errors.New("Could not open log file for reading.")
	}
	// Update the Zilla object with it's values.
	if this.Refresh() == false {
		return nil, errors.New("Could not get data from Hairball.")
	}
	return this, nil
}

func (this *Zilla) sendString(s, check string) bool {
	return this.sendBytes([]byte(s), check)
}

func (this *Zilla) sendIntValue(id string, val int) bool {
	if this.sendString(id, "") && this.sendString(strconv.Itoa(val), strconv.Itoa(val)) {
		return this.Refresh()
	}
	return false
}

func (this *Zilla) sendToggleValue(id string) bool {
	if this.sendString(id, "") {
		return this.Refresh()
	}
	return false
}

func (this *Zilla) sendBytes(b []byte, check string) bool {
	var e error
	_, e = this.serialPort.Write(b)
	if e != nil {
		fmt.Println(e)
		return false
	}
	// Cannot keep sleeping. Need a better solution here.
	if reflect.TypeOf(this.serialPort).String() != "*main.MockPort" {
		time.Sleep(500 * time.Millisecond)
	}
	this.buffer = make([]byte, 1000)
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

func (this *Zilla) menuHome() bool {
	this.writeLog = false
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

func (this *Zilla) startLogging() {
	go func() {
		this.menuHome()
		if this.sendString("p", "Special Menu:") == false {
			fmt.Println("Could not start logging.")
			return
		}
		// Open the first logging stream.
		_, readError := this.serialPort.Write([]byte("Q1\r"))
		if readError != nil {
			fmt.Println("Could read from logs from Hairball.")
			fmt.Println(readError)
			return
		}
		// Create a reader buffer.
		buf := bufio.NewReader(this.serialPort)
		// While allowed keep writing bytes to the file.
		this.writeLog = true
		for this.writeLog {
			line, readLineError := buf.ReadBytes('\n')
			if readLineError != nil {
				fmt.Println("Could read log line from Hairball.")
				fmt.Println(readLineError)
				return
			}
			_, writeLineError := this.writeLogFile.Write(ParseQ1LineFromHairball(line).ToBytes())
			if writeLineError != nil {
				fmt.Println(writeLineError)
				return
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()
}

func (this *Zilla) GetLiveData() *LiveData {
	bufSize := 1000
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
	buf = buf[:i]                              // Get the bytes written to the buffer.
	buf = buf[bytes.Index(buf, []byte{10})+1:] // Remove all bytes before the first line feed.
	buf = buf[:bytes.Index(buf, []byte{10})]   // Remove all bytes after the next line feed.
	return CreateLiveData(buf)
}

// Refreshes all attributes by reading them from the Zilla Controller.
func (this *Zilla) Refresh() bool {
	if this.menuSettings() == false {
		return false
	}
	defer this.startLogging()
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

// Battery Menu

func (this *Zilla) SetBatteryAmpLimit(val int) bool {
	this.menuBattery()
	this.sendIntValue("a", val)
	return this.BatteryAmpLimit == val
}

func (this *Zilla) SetLowBatteryVoltageLimit(val int) bool {
	this.menuBattery()
	return this.sendIntValue("v", val)
}

func (this *Zilla) SetLowBatteryVoltageIndicator(val int) bool {
	this.menuBattery()
	return this.sendIntValue("i", val)
}

// Motor Menu

func (this *Zilla) SetNormalMotorAmpLimit(val int) bool {
	this.menuMotor()
	return this.sendIntValue("a", val)
}

func (this *Zilla) SetSeriesMotorVoltageLimit(val int) bool {
	this.menuMotor()
	return this.sendIntValue("v", val)
}

func (this *Zilla) SetReverseMotorAmpLimit(val int) bool {
	this.menuMotor()
	return this.sendIntValue("i", val)
}

func (this *Zilla) SetReverseMotorVoltageLimit(val int) bool {
	this.menuMotor()
	return this.sendIntValue("r", val)
}

func (this *Zilla) SetParallelMotorAmpLimit(val int) bool {
	this.menuMotor()
	return this.sendIntValue("c", val)
}

func (this *Zilla) SetParallelMotorVoltageLimit(val int) bool {
	this.menuMotor()
	return this.sendIntValue("p", val)
}

// Speed Menu

func (this *Zilla) SetForwardRpmLimit(val int) bool {
	this.menuSpeed()
	return this.sendIntValue("l", val)
}

func (this *Zilla) SetReverseRpmLimit(val int) bool {
	this.menuSpeed()
	return this.sendIntValue("r", val)
}

func (this *Zilla) SetMaxRpmLimit(val int) bool {
	this.menuSpeed()
	return this.sendIntValue("x", val)
}

// Options Menu

func (this *Zilla) ToggleRpmSensorMotorOne() bool {
	this.menuOptions()
	return this.sendToggleValue("a")
}

func (this *Zilla) SetRpmSensorMotorTwo() bool {
	this.menuOptions()
	return this.sendToggleValue("b")
}

func (this *Zilla) ToggleAutoShiftingSeriesToParallel() bool {
	this.menuOptions()
	return this.sendToggleValue("c")
}

func (this *Zilla) ToggleStallDetectOn() bool {
	this.menuOptions()
	return this.sendToggleValue("d")
}

func (this *Zilla) ToggleBatteryLightPolarity() bool {
	this.menuOptions()
	return this.sendToggleValue("e")
}

func (this *Zilla) ToggleCheckEngineLightPolarity() bool {
	this.menuOptions()
	return this.sendToggleValue("f")
}

func (this *Zilla) ToggleReversingContactors() bool {
	this.menuOptions()
	return this.sendToggleValue("g")
}

func (this *Zilla) ToggleSeriesParallelContactors() bool {
	this.menuOptions()
	return this.sendToggleValue("h")
}

func (this *Zilla) SetForceParallelInReverse() bool {
	this.menuOptions()
	return this.sendToggleValue("i")
}

func (this *Zilla) ToggleInhibitSeriesParallelShifting() bool {
	this.menuOptions()
	return this.sendToggleValue("j")
}

func (this *Zilla) ToggleTachometerDisplayMotorAmps() bool {
	this.menuOptions()
	return this.sendToggleValue("k")
}

func (this *Zilla) ToggleTachometerSixCylinders() bool {
	this.menuOptions()
	return this.sendToggleValue("l")
}

func (this *Zilla) ToggleReversesPlugInInputPolarity() bool {
	this.menuOptions()
	return this.sendToggleValue("m")
}

func (this *Zilla) ToggleActivateHEPI() bool {
	this.menuOptions()
	return this.sendToggleValue("n")
}

func (this *Zilla) ToggleIsZ2k() bool {
	this.menuOptions()
	return this.sendToggleValue("p")
}
