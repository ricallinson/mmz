package main

import (
	"bufio"
	"bytes"
	"errors"
	"log"
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

type ZillaSettings struct {
	BatteryAmpLimit               int      // a) BA
	LowBatteryVoltageLimit        int      // v) LBV
	LowBatteryVoltageIndicator    int      // i) LBVI
	NormalMotorAmpLimit           int      // a) Amp
	SeriesMotorVoltageLimit       int      // v) Volt
	ReverseMotorAmpLimit          int      // i) RA
	ReverseMotorVoltageLimit      int      // r) RV
	ParallelMotorAmpLimit         int      // c) PA
	ParallelMotorVoltageLimit     int      // p) PV
	ForwardRpmLimit               int      // l) Norm
	ReverseRpmLimit               int      // r) Rev
	MaxRpmLimit                   int      // x) Max
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
}

type Zilla struct {
	serialPort   SerialPort
	queue        chan *zillaCommand
	logFile      string
	writeLog     bool
	readLogFile  *os.File
	writeLogFile *os.File
}

func NewZilla(p SerialPort) (*Zilla, error) {
	this := &Zilla{
		serialPort: p,
		queue:      make(chan *zillaCommand, 100),
	}
	// Open log file for reading and writing.
	if err := this.openLogFile(); err != nil {
		return nil, err
	}
	// Start listening on the queue.
	go this.start()
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
		default:
			this.writeLogToFile()
		}
	}
}

func (this *Zilla) writeCommands(types ...interface{}) []byte {
	cmd := newZillaCommand()
	for _, t := range types {
		switch t.(type) {
		case string:
			cmd.sendString(t.(string))
		case int:
			cmd.sendInt(t.(int))
		default:
			log.Println("Cannot send type")
			return nil
		}
	}
	this.writeCommand(cmd)
	return cmd.data
}

func (this *Zilla) writeCommand(cmd *zillaCommand) {
	this.queue <- cmd
	this.writeLog = false
	<-cmd.done
}

func (this *Zilla) writeBytes(b []byte) []byte {
	var e error
	_, e = this.serialPort.Write(b)
	// log.Println(string(b))
	if e != nil {
		log.Println(e)
		return nil
	}
	// Cannot keep sleeping. Need a better solution here.
	if reflect.TypeOf(this.serialPort).String() != "*main.MockPort" {
		time.Sleep(500 * time.Millisecond)
	}
	data := make([]byte, 1000)
	_, e = this.serialPort.Read(data)
	if e != nil {
		log.Println(e)
		return nil
	}
	data = bytes.TrimSpace(data)
	// log.Println(string(data))
	return data
}

func (this *Zilla) writeLogToFile() {
	this.serialPort.Write([]byte{27})  // Esc
	this.serialPort.Write([]byte{27})  // Esc
	this.serialPort.Write([]byte{27})  // Esc
	this.serialPort.Write([]byte("p")) // Menu Special
	_, readError := this.serialPort.Write([]byte("Q1\r"))
	if readError != nil {
		log.Println("Could not read logs from Hairball.")
		return
	}
	input := bufio.NewReader(this.serialPort)
	this.writeLog = true
	for this.writeLog {
		// This could be the real world problem.
		// The code could be waiting here for bytes and get the ones meant for the menu change.
		line, readLineError := input.ReadBytes('\n')
		if readLineError != nil {
			log.Println("Could read log line from Hairball.")
			log.Println(readLineError)
			return
		}
		logLine := ParseQ1LineFromHairball(line)
		if logLine == nil {
			log.Println("Could not parse Hairball log line.")
			return
		}
		if _, err := this.writeLogFile.Write(logLine.ToBytes()); err != nil {
			// If there is a write error it means the file has been closed.
			this.writeLog = false
			return
		}
		// Sleep for 100ms as the logs are only written 10 times a second.
		time.Sleep(100 * time.Millisecond)
	}
}

func (this *Zilla) openLogFile() error {
	// Set the log file for this session.
	this.logFile = "./logs/" + strconv.FormatInt(time.Now().Unix(), 10) + ".dat"
	// Make sure the directory is created.
	if err := os.MkdirAll(path.Dir(this.logFile), 0777); err != nil {
		log.Println(err)
		return errors.New("Could not create directory for logs.")
	}
	var openFileError error
	this.writeLogFile, openFileError = os.OpenFile(this.logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if openFileError != nil {
		log.Println(openFileError)
		return errors.New("Could not open log file for writing.")
	}
	this.readLogFile, openFileError = os.Open(this.logFile)
	if openFileError != nil {
		log.Println(openFileError)
		return errors.New("Could not open log file for reading.")
	}
	log.Print("Created log file:", this.logFile)
	return nil
}

func (this *Zilla) CloseLogFile() {
	this.writeLog = false
	this.readLogFile.Close()
	this.writeLogFile.Close()
}

func (this *Zilla) GetLiveData() *LiveData {
	if this.logFile == "" {
		log.Println("No log file available.")
		return nil
	}
	bufSize := 1000
	buf := make([]byte, bufSize)
	stat, _ := os.Stat(this.logFile)
	start := stat.Size() - int64(bufSize)
	i, err := this.readLogFile.ReadAt(buf, start)
	if err != nil {
		log.Println("Could not read last line from live data.")
		return nil
	}
	buf = buf[:i]                              // Get the bytes written to the buffer.
	buf = buf[bytes.Index(buf, []byte{10})+1:] // Remove all bytes before the first line feed.
	buf = buf[:bytes.Index(buf, []byte{10})]   // Remove all bytes after the next line feed.
	return CreateLiveData(buf)
}

// Refreshes all attributes by reading them from the Zilla Controller.
func (this *Zilla) GetSettings() *ZillaSettings {
	data := this.writeCommands("d")
	settings := &ZillaSettings{}
	// Read all the settings line by line.
	lines := bytes.Split(data, []byte{10})
	// Get values for BA, LBV, LBVI
	var values []string
	values = split(string(lines[2]), " ")
	settings.BatteryAmpLimit, _ = strconv.Atoi(values[0])
	settings.LowBatteryVoltageLimit, _ = strconv.Atoi(values[1])
	settings.LowBatteryVoltageIndicator, _ = strconv.Atoi(values[2])
	// Values for Amp, Volt, RA
	values = split(string(lines[4]), " ")
	settings.NormalMotorAmpLimit, _ = strconv.Atoi(values[0])
	settings.SeriesMotorVoltageLimit, _ = strconv.Atoi(values[1])
	settings.ReverseMotorAmpLimit, _ = strconv.Atoi(values[2])
	// Values for RV, PA, PV
	values = split(string(lines[6]), " ")
	settings.ReverseMotorVoltageLimit, _ = strconv.Atoi(values[0])
	settings.ParallelMotorAmpLimit, _ = strconv.Atoi(values[1])
	settings.ParallelMotorVoltageLimit, _ = strconv.Atoi(values[2])
	// Values for Norm, Rev, Max
	values = split(string(lines[8]), " ")
	settings.ForwardRpmLimit, _ = strconv.Atoi(values[0])
	settings.ReverseRpmLimit, _ = strconv.Atoi(values[1])
	settings.MaxRpmLimit, _ = strconv.Atoi(values[2])
	// Values for a, b, c, d
	values = split(string(lines[9]), "\t")
	settings.RpmSensorMotorOne = truthy(values[0])
	settings.RpmSensorMotorTwo = truthy(values[1])
	settings.AutoShiftingSeriesToParallel = truthy(values[2])
	settings.StallDetectOn = truthy(values[3])
	// Values for e, f, g, h
	values = split(string(lines[10]), "\t")
	settings.BatteryLightPolarity = truthy(values[0])
	settings.CheckEngineLightPolarity = truthy(values[1])
	settings.ReversingContactors = truthy(values[2])
	settings.SeriesParallelContactors = truthy(values[3])
	// Values for i, j, k, l
	values = split(string(lines[11]), "\t")
	settings.ForceParallelInReverse = truthy(values[0])
	settings.InhibitSeriesParallelShifting = truthy(values[1])
	settings.TachometerDisplayMotorAmps = truthy(values[2])
	settings.TachometerSixCylinders = truthy(values[3])
	// Values for m, n, o, p
	values = split(string(lines[12]), "\t")
	settings.ReversesPlugInInputPolarity = truthy(values[0])
	settings.ActivateHEPI = truthy(values[1])
	settings.notUsed = truthy(values[2])
	settings.IsZ2k = truthy(values[3])
	// Values for errors
	settings.Errors = split(string(lines[14]), " ")
	// Values for current state
	values = split(string(lines[15]), " ")
	settings.CurrentState = values[1]
	return settings
}

// Battery Menu

func (this *Zilla) SetBatteryAmpLimit(val int) bool {
	this.writeCommands("b", "a", val)
	log.Println("SetBatteryAmpLimit", val)
	return this.GetSettings().BatteryAmpLimit == val
}

func (this *Zilla) SetLowBatteryVoltageLimit(val int) bool {
	v := this.GetSettings().LowBatteryVoltageLimit
	this.writeCommands("b", "v", val)
	log.Println("SetLowBatteryVoltageLimit", val)
	return this.GetSettings().LowBatteryVoltageLimit != v
}

func (this *Zilla) SetLowBatteryVoltageIndicator(val int) bool {
	v := this.GetSettings().LowBatteryVoltageIndicator
	this.writeCommands("b", "i", val)
	log.Println("SetLowBatteryVoltageIndicator", val)
	return this.GetSettings().LowBatteryVoltageIndicator != v
}

// Motor Menu

func (this *Zilla) SetNormalMotorAmpLimit(val int) bool {
	v := this.GetSettings().NormalMotorAmpLimit
	this.writeCommands("m", "a", val)
	log.Println("SetNormalMotorAmpLimit", val)
	return this.GetSettings().NormalMotorAmpLimit != v
}

func (this *Zilla) SetSeriesMotorVoltageLimit(val int) bool {
	v := this.GetSettings().SeriesMotorVoltageLimit
	this.writeCommands("m", "v", val)
	log.Println("SetSeriesMotorVoltageLimit", val)
	return this.GetSettings().SeriesMotorVoltageLimit != v
}

func (this *Zilla) SetReverseMotorAmpLimit(val int) bool {
	v := this.GetSettings().ReverseMotorAmpLimit
	this.writeCommands("m", "i", val)
	log.Println("SetReverseMotorAmpLimit", val)
	return this.GetSettings().ReverseMotorAmpLimit != v
}

func (this *Zilla) SetReverseMotorVoltageLimit(val int) bool {
	v := this.GetSettings().ReverseMotorVoltageLimit
	this.writeCommands("m", "r", val)
	log.Println("SetReverseMotorVoltageLimit", val)
	return this.GetSettings().ReverseMotorVoltageLimit != v
}

func (this *Zilla) SetParallelMotorAmpLimit(val int) bool {
	v := this.GetSettings().ParallelMotorAmpLimit
	this.writeCommands("m", "c", val)
	log.Println("SetParallelMotorAmpLimit", val)
	return this.GetSettings().ParallelMotorAmpLimit != v
}

func (this *Zilla) SetParallelMotorVoltageLimit(val int) bool {
	v := this.GetSettings().ParallelMotorVoltageLimit
	this.writeCommands("m", "p", val)
	log.Println("SetParallelMotorVoltageLimit", val)
	return this.GetSettings().ParallelMotorVoltageLimit != v
}

// Speed Menu

func (this *Zilla) SetForwardRpmLimit(val int) bool {
	this.writeCommands("s", "l", val)
	log.Println("SetForwardRpmLimit", val)
	return this.GetSettings().ForwardRpmLimit == val
}

func (this *Zilla) SetReverseRpmLimit(val int) bool {
	this.writeCommands("s", "r", val)
	log.Println("SetReverseRpmLimit", val)
	return this.GetSettings().ReverseRpmLimit == val
}

func (this *Zilla) SetMaxRpmLimit(val int) bool {
	this.writeCommands("s", "x", val)
	log.Println("SetMaxRpmLimit", val)
	return this.GetSettings().MaxRpmLimit == val
}

// Options Menu

func (this *Zilla) ToggleRpmSensorMotorOne() bool {
	v := this.GetSettings().RpmSensorMotorOne
	this.writeCommands("o", "a")
	log.Println("ToggleRpmSensorMotorOne")
	return this.GetSettings().RpmSensorMotorOne != v
}

func (this *Zilla) ToggleRpmSensorMotorTwo() bool {
	v := this.GetSettings().RpmSensorMotorTwo
	this.writeCommands("o", "b")
	log.Println("ToggleRpmSensorMotorTwo")
	return this.GetSettings().RpmSensorMotorTwo != v
}

func (this *Zilla) ToggleAutoShiftingSeriesToParallel() bool {
	v := this.GetSettings().AutoShiftingSeriesToParallel
	this.writeCommands("o", "c")
	log.Println("ToggleAutoShiftingSeriesToParallel")
	return this.GetSettings().AutoShiftingSeriesToParallel != v
}

func (this *Zilla) ToggleStallDetectOn() bool {
	v := this.GetSettings().StallDetectOn
	this.writeCommands("o", "d")
	log.Println("ToggleStallDetectOn")
	return this.GetSettings().StallDetectOn != v
}

func (this *Zilla) ToggleBatteryLightPolarity() bool {
	v := this.GetSettings().BatteryLightPolarity
	this.writeCommands("o", "e")
	log.Println("ToggleBatteryLightPolarity")
	return this.GetSettings().BatteryLightPolarity != v
}

func (this *Zilla) ToggleCheckEngineLightPolarity() bool {
	v := this.GetSettings().CheckEngineLightPolarity
	this.writeCommands("o", "f")
	log.Println("ToggleCheckEngineLightPolarity")
	return this.GetSettings().CheckEngineLightPolarity != v
}

func (this *Zilla) ToggleReversingContactors() bool {
	v := this.GetSettings().ReversingContactors
	this.writeCommands("o", "g")
	log.Println("ToggleReversingContactors")
	return this.GetSettings().ReversingContactors != v
}

func (this *Zilla) ToggleSeriesParallelContactors() bool {
	v := this.GetSettings().SeriesParallelContactors
	this.writeCommands("o", "h")
	log.Println("ToggleSeriesParallelContactors")
	return this.GetSettings().SeriesParallelContactors != v
}

func (this *Zilla) ToggleForceParallelInReverse() bool {
	v := this.GetSettings().ForceParallelInReverse
	this.writeCommands("o", "i")
	log.Println("ToggleForceParallelInReverse")
	return this.GetSettings().ForceParallelInReverse != v
}

func (this *Zilla) ToggleInhibitSeriesParallelShifting() bool {
	v := this.GetSettings().InhibitSeriesParallelShifting
	this.writeCommands("o", "j")
	log.Println("ToggleInhibitSeriesParallelShifting")
	return this.GetSettings().InhibitSeriesParallelShifting != v
}

func (this *Zilla) ToggleTachometerDisplayMotorAmps() bool {
	v := this.GetSettings().TachometerDisplayMotorAmps
	this.writeCommands("o", "k")
	log.Println("ToggleTachometerDisplayMotorAmps")
	return this.GetSettings().TachometerDisplayMotorAmps != v
}

func (this *Zilla) ToggleTachometerSixCylinders() bool {
	v := this.GetSettings().TachometerSixCylinders
	this.writeCommands("o", "l")
	log.Println("ToggleTachometerSixCylinders")
	return this.GetSettings().TachometerSixCylinders != v
}

func (this *Zilla) ToggleReversesPlugInInputPolarity() bool {
	v := this.GetSettings().ReversesPlugInInputPolarity
	this.writeCommands("o", "m")
	log.Println("ToggleReversesPlugInInputPolarity")
	return this.GetSettings().ReversesPlugInInputPolarity != v
}

func (this *Zilla) ToggleActivateHEPI() bool {
	v := this.GetSettings().ActivateHEPI
	this.writeCommands("o", "n")
	log.Println("ToggleActivateHEPI")
	return this.GetSettings().ActivateHEPI != v
}

func (this *Zilla) ToggleIsZ2k() bool {
	v := this.GetSettings().IsZ2k
	this.writeCommands("o", "p")
	log.Println("ToggleIsZ2k")
	return this.GetSettings().IsZ2k != v
}
