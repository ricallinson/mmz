package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path"
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
	closed       bool
	readLogFile  *os.File
	writeLogFile *os.File
}

func NewZilla(p SerialPort) (*Zilla, error) {
	this := &Zilla{
		serialPort: p,
	}
	// Return the Zilla instance.
	return this, nil
}

func (this *Zilla) sendCommand(cmd *zillaCommand) {
	for _, b := range cmd.bytes {
		// log.Println("Command:", string(b))
		this.writeBytes(b)
		cmd.data = this.readAllBytes()
		// log.Print(string(cmd.data))
	}
}

// Takes a list of commands and sends them to the Zilla in the order received.
func (this *Zilla) sendCommands(types ...interface{}) []byte {
	cmd := newZillaCommand()
	for _, t := range types {
		switch t.(type) {
		case string:
			cmd.sendString(t.(string))
		case int:
			cmd.sendInt(t.(int))
		default:
			log.Println("Can only write types 'string' or 'int'")
			return nil
		}
	}
	this.sendCommand(cmd)
	// log.Println("Read data: ", string(cmd.data))
	return cmd.data
}

func (this *Zilla) Raw(types ...interface{}) []byte {
	return this.sendCommands(types)
}

// Sends the given bytes directly to the Zilla.
func (this *Zilla) writeBytes(b []byte) {
	if this.closed {
		log.Print("Write connection to Zilla has been closed.")
		return
	}
	_, e := this.serialPort.Write(b)
	if e != nil {
		log.Println(e)
	}
}

// Reads bytes directly from the Zilla up to EOF.
func (this *Zilla) readAllBytes() []byte {
	return this.readBytes(0)
}

// Reads bytes directly from the Zilla up to the given 'delim' or EOF.
func (this *Zilla) readBytes(delim byte) []byte {
	if this.closed {
		log.Print("Read connection to Zilla has been closed.")
		return nil
	}
	limit := 1000
	buff := make([]byte, 1)
	data := make([]byte, 0)
	for limit > 0 {
		i, _ := this.serialPort.Read(buff)
		// log.Println(i, buff[0], delim)
		if i == 0 || buff[0] == delim {
			break
		}
		data = append(data, buff[0])
		// log.Println(i, string(buff[0]))
		// log.Println(string(data))
		limit--
	}
	data = bytes.TrimSpace(data)
	// log.Println(string(data))
	return data
}

// Creates a new log file for this instance to write to.
func (this *Zilla) createLogFile(logFile string) error {
	// Make sure the directory is created.
	if err := os.MkdirAll(path.Dir(logFile), 0777); err != nil {
		log.Println(err)
		return errors.New("Could not create directory for logs.")
	}
	var openFileError error
	this.writeLogFile, openFileError = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if openFileError != nil {
		log.Println(openFileError)
		return errors.New("Could not open log file for writing.")
	}
	this.readLogFile, openFileError = os.Open(logFile)
	if openFileError != nil {
		log.Println(openFileError)
		return errors.New("Could not open log file for reading.")
	}
	log.Print("Created log file:", logFile)
	return nil
}

func (this *Zilla) appendLineToLog(line []byte) {
	// If the line is too short we can't read it so just return (20 is a arbitrary number).
	if len(line) < 20 {
		return
	}
	logLine := ParseQ1LineFromHairball(line)
	if logLine == nil {
		log.Print("Log line [", string(line), "] could not be parsed.")
		return
	}
	if _, err := this.writeLogFile.Write(logLine.ToBytes()); err != nil {
		// If there is a write error it means the log file has been closed.
		// log.Print("Log file has been closed.")
		return
	}
}

// Blocks while writing log file.
func (this *Zilla) Log(logFile string) {
	// Open log file for reading and writing.
	if err := this.createLogFile(logFile); err != nil {
		return
	}
	// We have to write and read each command to be sure the Zilla is responding.
	// The sendCommands() function is not used as it creates a circular dependency.
	this.writeBytes([]byte{27})
	this.readAllBytes()
	this.writeBytes([]byte{27})
	this.readAllBytes()
	this.writeBytes([]byte{27})
	this.readAllBytes()
	this.writeBytes([]byte("p")) // Menu Special
	this.readAllBytes()
	this.writeBytes([]byte("Q1\r")) // Start logs
	for !this.closed {
		// log.Print("Logging.")
		// Read the log line from the Zilla.
		line := bytes.TrimSpace(this.readBytes('\n'))
		// log.Print(string(line))
		this.appendLineToLog(line)
		// Sleep for 100ms as the logs are only written 10 times a second.
		time.Sleep(100 * time.Millisecond)
	}
}

// Ends the connection to the Zilla and stops all logging.
// It's good practice to call Close() once the instance is no longer needed.
// Once called the Zilla instance can no longer be use to send commands.
func (this *Zilla) Close() {
	this.closed = true
	this.readLogFile.Close()
	this.writeLogFile.Close()
}

func (this *Zilla) GetLiveData() *LiveData {
	if this.writeLogFile == nil {
		log.Println("No log file available.")
		return nil
	}
	bufSize := 1000
	buf := make([]byte, bufSize)
	stat, _ := this.readLogFile.Stat()
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
	data := this.sendCommands("d")
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
	v := this.GetSettings().BatteryAmpLimit
	this.sendCommands("b", "a", val)
	n := this.GetSettings().BatteryAmpLimit
	// log.Println("SetBatteryAmpLimit", n)
	return n != v
}

func (this *Zilla) SetLowBatteryVoltageLimit(val int) bool {
	v := this.GetSettings().LowBatteryVoltageLimit
	this.sendCommands("b", "v", val)
	n := this.GetSettings().LowBatteryVoltageLimit
	// log.Println("SetLowBatteryVoltageLimit", n)
	return n != v
}

func (this *Zilla) SetLowBatteryVoltageIndicator(val int) bool {
	v := this.GetSettings().LowBatteryVoltageIndicator
	this.sendCommands("b", "i", val)
	n := this.GetSettings().LowBatteryVoltageIndicator
	// log.Println("SetLowBatteryVoltageIndicator", n)
	return n != v
}

// Motor Menu

func (this *Zilla) SetNormalMotorAmpLimit(val int) bool {
	v := this.GetSettings().NormalMotorAmpLimit
	this.sendCommands("m", "a", val)
	n := this.GetSettings().NormalMotorAmpLimit
	// log.Println("SetNormalMotorAmpLimit", n)
	return n != v
}

func (this *Zilla) SetSeriesMotorVoltageLimit(val int) bool {
	v := this.GetSettings().SeriesMotorVoltageLimit
	this.sendCommands("m", "v", val)
	n := this.GetSettings().SeriesMotorVoltageLimit
	// log.Println("SetSeriesMotorVoltageLimit", n)
	return n != v
}

func (this *Zilla) SetReverseMotorAmpLimit(val int) bool {
	v := this.GetSettings().ReverseMotorAmpLimit
	this.sendCommands("m", "i", val)
	n := this.GetSettings().ReverseMotorAmpLimit
	// log.Println("SetReverseMotorAmpLimit", n)
	return n != v
}

func (this *Zilla) SetReverseMotorVoltageLimit(val int) bool {
	v := this.GetSettings().ReverseMotorVoltageLimit
	this.sendCommands("m", "r", val)
	n := this.GetSettings().ReverseMotorVoltageLimit
	// log.Println("SetReverseMotorVoltageLimit", n)
	return n != v
}

func (this *Zilla) SetParallelMotorAmpLimit(val int) bool {
	v := this.GetSettings().ParallelMotorAmpLimit
	this.sendCommands("m", "c", val)
	n := this.GetSettings().ParallelMotorAmpLimit
	// log.Println("SetParallelMotorAmpLimit", n)
	return n != v
}

func (this *Zilla) SetParallelMotorVoltageLimit(val int) bool {
	v := this.GetSettings().ParallelMotorVoltageLimit
	this.sendCommands("m", "p", val)
	n := this.GetSettings().ParallelMotorVoltageLimit
	// log.Println("SetParallelMotorVoltageLimit", n)
	return n != v
}

// Speed Menu

func (this *Zilla) SetForwardRpmLimit(val int) bool {
	this.sendCommands("s", "l", val)
	n := this.GetSettings().ForwardRpmLimit
	// log.Println("SetForwardRpmLimit", n)
	return n == val
}

func (this *Zilla) SetReverseRpmLimit(val int) bool {
	this.sendCommands("s", "r", val)
	n := this.GetSettings().ReverseRpmLimit
	// log.Println("SetReverseRpmLimit", n)
	return n == val
}

func (this *Zilla) SetMaxRpmLimit(val int) bool {
	this.sendCommands("s", "x", val)
	n := this.GetSettings().MaxRpmLimit
	// log.Println("SetMaxRpmLimit", n)
	return n == val
}

// Options Menu

func (this *Zilla) ToggleRpmSensorMotorOne() bool {
	v := this.GetSettings().RpmSensorMotorOne
	this.sendCommands("o", "a")
	n := this.GetSettings().RpmSensorMotorOne
	// log.Println("ToggleRpmSensorMotorOne", n)
	return n != v
}

func (this *Zilla) ToggleRpmSensorMotorTwo() bool {
	v := this.GetSettings().RpmSensorMotorTwo
	this.sendCommands("o", "b")
	n := this.GetSettings().RpmSensorMotorTwo
	// log.Println("ToggleRpmSensorMotorTwo", n)
	return n != v
}

func (this *Zilla) ToggleAutoShiftingSeriesToParallel() bool {
	v := this.GetSettings().AutoShiftingSeriesToParallel
	this.sendCommands("o", "c")
	n := this.GetSettings().AutoShiftingSeriesToParallel
	// log.Println("ToggleAutoShiftingSeriesToParallel", n)
	return n != v
}

func (this *Zilla) ToggleStallDetectOn() bool {
	v := this.GetSettings().StallDetectOn
	this.sendCommands("o", "d")
	n := this.GetSettings().StallDetectOn
	// log.Println("ToggleStallDetectOn", n)
	return n != v
}

func (this *Zilla) ToggleBatteryLightPolarity() bool {
	v := this.GetSettings().BatteryLightPolarity
	this.sendCommands("o", "e")
	n := this.GetSettings().BatteryLightPolarity
	// log.Println("ToggleBatteryLightPolarity", n)
	return n != v
}

func (this *Zilla) ToggleCheckEngineLightPolarity() bool {
	v := this.GetSettings().CheckEngineLightPolarity
	this.sendCommands("o", "f")
	n := this.GetSettings().CheckEngineLightPolarity
	// log.Println("ToggleCheckEngineLightPolarity", n)
	return n != v
}

func (this *Zilla) ToggleReversingContactors() bool {
	v := this.GetSettings().ReversingContactors
	this.sendCommands("o", "g")
	n := this.GetSettings().ReversingContactors
	// log.Println("ToggleReversingContactors", n)
	return n != v
}

func (this *Zilla) ToggleSeriesParallelContactors() bool {
	v := this.GetSettings().SeriesParallelContactors
	this.sendCommands("o", "h")
	n := this.GetSettings().SeriesParallelContactors
	// log.Println("ToggleSeriesParallelContactors", n)
	return n != v
}

func (this *Zilla) ToggleForceParallelInReverse() bool {
	v := this.GetSettings().ForceParallelInReverse
	this.sendCommands("o", "i")
	n := this.GetSettings().ForceParallelInReverse
	// log.Println("ToggleForceParallelInReverse", n)
	return n != v
}

func (this *Zilla) ToggleInhibitSeriesParallelShifting() bool {
	v := this.GetSettings().InhibitSeriesParallelShifting
	this.sendCommands("o", "j")
	n := this.GetSettings().InhibitSeriesParallelShifting
	// log.Println("ToggleInhibitSeriesParallelShifting", n)
	return n != v
}

func (this *Zilla) ToggleTachometerDisplayMotorAmps() bool {
	v := this.GetSettings().TachometerDisplayMotorAmps
	this.sendCommands("o", "k")
	n := this.GetSettings().TachometerDisplayMotorAmps
	// log.Println("ToggleTachometerDisplayMotorAmps", n)
	return n != v
}

func (this *Zilla) ToggleTachometerSixCylinders() bool {
	v := this.GetSettings().TachometerSixCylinders
	this.sendCommands("o", "l")
	n := this.GetSettings().TachometerSixCylinders
	// log.Println("ToggleTachometerSixCylinders", n)
	return n != v
}

func (this *Zilla) ToggleReversesPlugInInputPolarity() bool {
	v := this.GetSettings().ReversesPlugInInputPolarity
	this.sendCommands("o", "m")
	n := this.GetSettings().ReversesPlugInInputPolarity
	// log.Println("ToggleReversesPlugInInputPolarity", n)
	return n != v
}

func (this *Zilla) ToggleActivateHEPI() bool {
	v := this.GetSettings().ActivateHEPI
	this.sendCommands("o", "n")
	n := this.GetSettings().ActivateHEPI
	// log.Println("ToggleActivateHEPI", n)
	return n != v
}

func (this *Zilla) ToggleIsZ2k() bool {
	v := this.GetSettings().IsZ2k
	this.sendCommands("o", "p")
	n := this.GetSettings().IsZ2k
	// log.Println("ToggleIsZ2k", n)
	return n != v
}
