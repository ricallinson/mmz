package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"
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
	queue                         chan *zillaCommand
	buffer                        []byte // byte array of the last Zilla output
	serialPort                    SerialPort
	writeLog                      bool
	logFile                       string
	readLogFile                   *os.File
	writeLogFile                  *os.File
}

func NewZilla(p SerialPort) (*Zilla, error) {
	this := &Zilla{
		Errors:     make([]string, 0),
		serialPort: p,
		queue:      make(chan *zillaCommand, 100),
	}
	// Open log file for reading and writing.
	if err := this.OpenLog(); err != nil {
		return nil, err
	}
	go this.start()
	// Update the Zilla object with it's values.
	if this.Refresh() == false {
		return nil, errors.New("Could not get data from Hairball.")
	}
	return this, nil
}

// Loop on the queue channel sending commands from the queue to the Hairball.
// If there are no commands in the queue log data until a command appears.
func (this *Zilla) start() {
	for {
		select {
		case cmd := <-this.queue:
			for _, b := range cmd.bytes {
				cmd.data = this.writeBytes(b)
			}
			cmd.done <- true
		}
	}
}

func (this *Zilla) writeCommand(cmd *zillaCommand) {
	this.queue <- cmd
	<-cmd.done
}

func (this *Zilla) writeBytes(b []byte) []byte {
	var e error
	_, e = this.serialPort.Write(b)
	// fmt.Println(string(b))
	if e != nil {
		fmt.Println(e)
		return nil
	}
	// Cannot keep sleeping. Need a better solution here.
	if reflect.TypeOf(this.serialPort).String() != "*main.MockPort" {
		time.Sleep(500 * time.Millisecond)
	}
	data := make([]byte, 1000)
	_, e = this.serialPort.Read(data)
	if e != nil {
		fmt.Println(e)
		return nil
	}
	data = bytes.TrimSpace(data)
	// fmt.Println(string(data))
	return data
}

func (this *Zilla) sendString(s, check string) bool {
	return this.sendBytes([]byte(s), check)
}

func (this *Zilla) sendIntValue(id string, val int) bool {
	if this.sendString(id, "") && this.sendString(strconv.Itoa(val)+"\r\n", "") {
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
	// Debug sent data.
	// fmt.Println(string(b))
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
	// Debug returned data.
	// fmt.Println(string(this.buffer))
	if e != nil {
		fmt.Println(e)
		return false
	}
	this.buffer = bytes.TrimSpace(this.buffer)
	// Debug hairball output.
	// fmt.Println(string(this.buffer))
	return bytes.Index(this.buffer, []byte(check)) > -1
}

func (this *Zilla) OpenLog() error {
	// Set the log file for this session.
	this.logFile = "./logs/" + strconv.FormatInt(time.Now().Unix(), 10) + ".dat"
	// Make sure the directory is created.
	if err := os.MkdirAll(path.Dir(this.logFile), 0777); err != nil {
		fmt.Println(err)
		return errors.New("Could not create directory for logs.")
	}
	var openFileError error
	this.writeLogFile, openFileError = os.OpenFile(this.logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if openFileError != nil {
		fmt.Println(openFileError)
		return errors.New("Could not open log file for writing.")
	}
	this.readLogFile, openFileError = os.Open(this.logFile)
	if openFileError != nil {
		fmt.Println(openFileError)
		return errors.New("Could not open log file for reading.")
	}
	return nil
}

func (this *Zilla) CloseLog() {
	this.writeLog = false
	this.readLogFile.Close()
	this.writeLogFile.Close()
}

func (this *Zilla) startLogging() {
	this.menuHome()
	if this.sendString("p", "Special Menu:") == false {
		fmt.Println("Could not start logging.")
		return
	}
	// Open the first logging stream.
	_, readError := this.serialPort.Write([]byte("Q1\r"))
	if readError != nil {
		fmt.Println("Could not read logs from Hairball.")
		fmt.Println(readError)
		return
	}
	// Create a reader buffer.
	buf := bufio.NewReader(this.serialPort)
	// While allowed keep writing bytes to the file.
	this.writeLog = true
	// Now we have a buffer reading from the log keep reading and writing.
	go func(zilla *Zilla, input *bufio.Reader) {
		for zilla.writeLog {
			// This could be the real world problem.
			// The code could be waiting here for bytes and get the ones meant for the menu change.
			line, readLineError := input.ReadBytes('\n')
			if readLineError != nil {
				fmt.Println("Could read log line from Hairball.")
				fmt.Println(readLineError)
				return
			}
			logLine := ParseQ1LineFromHairball(line)
			if logLine == nil {
				fmt.Println("Could not parse Hairball log line.")
				return
			}
			if _, err := zilla.writeLogFile.Write(logLine.ToBytes()); err != nil {
				// If there is a write error it means the file has been closed.
				return
			}
			// Sleep for 100ms as the logs are only written 10 times a second.
			time.Sleep(100 * time.Millisecond)
		}
	}(this, buf)
}

func (this *Zilla) menuHome() bool {
	// Is there a better way to stop logging?
	this.writeLog = false
	// Is the result of this sendBytes() read by the startLogging() function?
	this.sendBytes([]byte{27}, "")
	// We should have stopped logging by now.
	this.sendBytes([]byte{27}, "")
	this.sendBytes([]byte{27}, "")
	// This sendBytes() should put us in the home menu with no logging.
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

func (this *Zilla) GetLiveData() *LiveData {
	if this.logFile == "" {
		fmt.Println("No log file available.")
		return nil
	}
	bufSize := 1000
	buf := make([]byte, bufSize)
	stat, _ := os.Stat(this.logFile)
	start := stat.Size() - int64(bufSize)
	i, err := this.readLogFile.ReadAt(buf, start)
	if err != nil {
		fmt.Println("Could not read last line from live data.")
		return nil
	}
	buf = buf[:i]                              // Get the bytes written to the buffer.
	buf = buf[bytes.Index(buf, []byte{10})+1:] // Remove all bytes before the first line feed.
	buf = buf[:bytes.Index(buf, []byte{10})]   // Remove all bytes after the next line feed.
	return CreateLiveData(buf)
}

// Refreshes all attributes by reading them from the Zilla Controller.
func (this *Zilla) GetSettings() bool {
	// Add Zilla command sequence to a queue.
	// Wait for the command to execute and return.
	// Parse the returned buffer.
	return false
}

// Refreshes all attributes by reading them from the Zilla Controller.
func (this *Zilla) Refresh() bool {
	cmd := newZillaCommand()
	cmd.sendHome()
	cmd.sendString("d") // Menu Settings
	this.writeCommand(cmd)
	// Read all the settings line by line.
	lines := bytes.Split(cmd.data, []byte{10})
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
	v := this.LowBatteryVoltageLimit
	this.sendIntValue("v", val)
	return this.LowBatteryVoltageLimit != v
}

func (this *Zilla) SetLowBatteryVoltageIndicator(val int) bool {
	this.menuBattery()
	v := this.LowBatteryVoltageIndicator
	this.sendIntValue("i", val)
	return this.LowBatteryVoltageIndicator != v
}

// Motor Menu

func (this *Zilla) SetNormalMotorAmpLimit(val int) bool {
	this.menuMotor()
	v := this.NormalMotorAmpLimit
	this.sendIntValue("a", val)
	return this.NormalMotorAmpLimit != v
}

func (this *Zilla) SetSeriesMotorVoltageLimit(val int) bool {
	this.menuMotor()
	v := this.SeriesMotorVoltageLimit
	this.sendIntValue("v", val)
	return this.SeriesMotorVoltageLimit != v
}

func (this *Zilla) SetReverseMotorAmpLimit(val int) bool {
	this.menuMotor()
	v := this.ReverseMotorAmpLimit
	this.sendIntValue("i", val)
	return this.ReverseMotorAmpLimit != v
}

func (this *Zilla) SetReverseMotorVoltageLimit(val int) bool {
	this.menuMotor()
	v := this.ReverseMotorVoltageLimit
	this.sendIntValue("r", val)
	return this.ReverseMotorVoltageLimit != v
}

func (this *Zilla) SetParallelMotorAmpLimit(val int) bool {
	this.menuMotor()
	v := this.ParallelMotorAmpLimit
	this.sendIntValue("c", val)
	return this.ParallelMotorAmpLimit != v
}

func (this *Zilla) SetParallelMotorVoltageLimit(val int) bool {
	this.menuMotor()
	v := this.ParallelMotorVoltageLimit
	this.sendIntValue("p", val)
	return this.ParallelMotorVoltageLimit != v
}

// Speed Menu

func (this *Zilla) SetForwardRpmLimit(val int) bool {
	this.menuSpeed()
	this.sendIntValue("l", val)
	return this.ForwardRpmLimit == val
}

func (this *Zilla) SetReverseRpmLimit(val int) bool {
	this.menuSpeed()
	return this.sendIntValue("r", val)
	return this.ReverseRpmLimit == val
}

func (this *Zilla) SetMaxRpmLimit(val int) bool {
	this.menuSpeed()
	return this.sendIntValue("x", val)
	return this.MaxRpmLimit == val
}

// Options Menu

func (this *Zilla) ToggleRpmSensorMotorOne() bool {
	this.menuOptions()
	v := this.RpmSensorMotorOne
	this.sendToggleValue("a")
	return this.RpmSensorMotorOne != v
}

func (this *Zilla) ToggleRpmSensorMotorTwo() bool {
	this.menuOptions()
	v := this.RpmSensorMotorTwo
	this.sendToggleValue("b")
	return this.RpmSensorMotorTwo != v
}

func (this *Zilla) ToggleAutoShiftingSeriesToParallel() bool {
	this.menuOptions()
	v := this.AutoShiftingSeriesToParallel
	this.sendToggleValue("c")
	return this.AutoShiftingSeriesToParallel != v
}

func (this *Zilla) ToggleStallDetectOn() bool {
	this.menuOptions()
	v := this.StallDetectOn
	this.sendToggleValue("d")
	return this.StallDetectOn != v
}

func (this *Zilla) ToggleBatteryLightPolarity() bool {
	this.menuOptions()
	v := this.BatteryLightPolarity
	this.sendToggleValue("e")
	return this.BatteryLightPolarity != v
}

func (this *Zilla) ToggleCheckEngineLightPolarity() bool {
	this.menuOptions()
	v := this.CheckEngineLightPolarity
	this.sendToggleValue("f")
	return this.CheckEngineLightPolarity != v
}

func (this *Zilla) ToggleReversingContactors() bool {
	this.menuOptions()
	v := this.ReversingContactors
	this.sendToggleValue("g")
	return this.ReversingContactors != v
}

func (this *Zilla) ToggleSeriesParallelContactors() bool {
	this.menuOptions()
	v := this.SeriesParallelContactors
	this.sendToggleValue("h")
	return this.SeriesParallelContactors != v
}

func (this *Zilla) ToggleForceParallelInReverse() bool {
	this.menuOptions()
	v := this.ForceParallelInReverse
	this.sendToggleValue("i")
	return this.ForceParallelInReverse != v
}

func (this *Zilla) ToggleInhibitSeriesParallelShifting() bool {
	this.menuOptions()
	v := this.InhibitSeriesParallelShifting
	this.sendToggleValue("j")
	return this.InhibitSeriesParallelShifting != v
}

func (this *Zilla) ToggleTachometerDisplayMotorAmps() bool {
	this.menuOptions()
	v := this.TachometerDisplayMotorAmps
	this.sendToggleValue("k")
	return this.TachometerDisplayMotorAmps != v
}

func (this *Zilla) ToggleTachometerSixCylinders() bool {
	this.menuOptions()
	v := this.TachometerSixCylinders
	this.sendToggleValue("l")
	return this.TachometerSixCylinders != v
}

func (this *Zilla) ToggleReversesPlugInInputPolarity() bool {
	this.menuOptions()
	v := this.ReversesPlugInInputPolarity
	this.sendToggleValue("m")
	return this.ReversesPlugInInputPolarity != v
}

func (this *Zilla) ToggleActivateHEPI() bool {
	this.menuOptions()
	v := this.ActivateHEPI
	this.sendToggleValue("n")
	return this.ActivateHEPI != v
}

func (this *Zilla) ToggleIsZ2k() bool {
	this.menuOptions()
	v := this.IsZ2k
	this.sendToggleValue("p")
	return this.IsZ2k != v
}
