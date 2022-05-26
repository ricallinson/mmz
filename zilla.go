package main

import (
	"bytes"
	"fmt"
	"log"
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
	IsZ1k                         bool     // p) Off
	CurrentState                  string   // 1311
	Errors                        []string // 1111, 1111, ...
}

type Zilla struct {
	serialPort SerialPort
	closed     bool
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

// Blocks while writing to console out.
func (this *Zilla) RealtimeValues(queue string, format string) {
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
	this.writeBytes([]byte(queue + "\r")) // Start data acquisition.
	var fn func(interface{}) []byte
	switch format {
	case "raw":
		fn = interfaceToBytes
	case "json":
		fn = interfaceToJson
	default:
		fn = interfaceToYaml
	}
	for !this.closed {
		// Read the log line from the Zilla.
		line := bytes.TrimSpace(this.readBytes('\n'))
		// Select which data acquisition queue to read from.
		switch queue {
		case "Q4":
			fmt.Println(fn(GetRealtimeValuesQ4(line)))
		default:
			fmt.Println(fn(GetRealtimeValuesQ1(line)))
		}
		// Sleep for 100ms as the logs are only written 10 times a second.
		time.Sleep(100 * time.Millisecond)
	}
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
	settings.IsZ1k = truthy(values[3])
	// Values for errors
	settings.Errors = split(string(lines[14]), " ")
	for i := range settings.Errors {
		settings.Errors[i] = settings.Errors[i] + ": " + Codes[settings.Errors[i]]
	}
	// Values for current state
	values = split(string(lines[15]), " ")
	settings.CurrentState = values[1]
	return settings
}

// Battery Menu

func (this *Zilla) SetBatteryAmpLimit(val int) bool {
	this.sendCommands("b", "a", val)
	n := this.GetSettings().BatteryAmpLimit
	// log.Println("SetBatteryAmpLimit", n)
	return n == val
}

func (this *Zilla) SetLowBatteryVoltageLimit(val int) bool {
	this.sendCommands("b", "v", val)
	n := this.GetSettings().LowBatteryVoltageLimit
	// log.Println("SetLowBatteryVoltageLimit", n)
	return n == val
}

func (this *Zilla) SetLowBatteryVoltageIndicator(val int) bool {
	this.sendCommands("b", "i", val)
	n := this.GetSettings().LowBatteryVoltageIndicator
	// log.Println("SetLowBatteryVoltageIndicator", n)
	return n == val
}

// Motor Menu

func (this *Zilla) SetNormalMotorAmpLimit(val int) bool {
	this.sendCommands("m", "a", val)
	n := this.GetSettings().NormalMotorAmpLimit
	// log.Println("SetNormalMotorAmpLimit", n)
	return n == val
}

func (this *Zilla) SetSeriesMotorVoltageLimit(val int) bool {
	this.sendCommands("m", "v", val)
	n := this.GetSettings().SeriesMotorVoltageLimit
	// log.Println("SetSeriesMotorVoltageLimit", n)
	return n == val
}

func (this *Zilla) SetReverseMotorAmpLimit(val int) bool {
	this.sendCommands("m", "i", val)
	n := this.GetSettings().ReverseMotorAmpLimit
	// log.Println("SetReverseMotorAmpLimit", n)
	return n == val
}

func (this *Zilla) SetReverseMotorVoltageLimit(val int) bool {
	this.sendCommands("m", "r", val)
	n := this.GetSettings().ReverseMotorVoltageLimit
	// log.Println("SetReverseMotorVoltageLimit", n)
	return n == val
}

func (this *Zilla) SetParallelMotorAmpLimit(val int) bool {
	this.sendCommands("m", "c", val)
	n := this.GetSettings().ParallelMotorAmpLimit
	// log.Println("SetParallelMotorAmpLimit", n)
	return n == val
}

func (this *Zilla) SetParallelMotorVoltageLimit(val int) bool {
	this.sendCommands("m", "p", val)
	n := this.GetSettings().ParallelMotorVoltageLimit
	// log.Println("SetParallelMotorVoltageLimit", n)
	return n == val
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

func (this *Zilla) SetRpmSensorMotorOne(b bool) bool {
	v := this.GetSettings().RpmSensorMotorOne
	if b == v {
		return true
	}
	this.sendCommands("o", "a")
	n := this.GetSettings().RpmSensorMotorOne
	// log.Println("SetRpmSensorMotorOne", n)
	return n == b
}

func (this *Zilla) SetRpmSensorMotorTwo(b bool) bool {
	v := this.GetSettings().RpmSensorMotorTwo
	if b == v {
		return true
	}
	this.sendCommands("o", "b")
	n := this.GetSettings().RpmSensorMotorTwo
	// log.Println("SetRpmSensorMotorTwo", n)
	return n == b
}

func (this *Zilla) SetAutoShiftingSeriesToParallel(b bool) bool {
	v := this.GetSettings().AutoShiftingSeriesToParallel
	if b == v {
		return true
	}
	this.sendCommands("o", "c")
	n := this.GetSettings().AutoShiftingSeriesToParallel
	// log.Println("SetAutoShiftingSeriesToParallel", n)
	return n == b
}

func (this *Zilla) SetStallDetectOn(b bool) bool {
	v := this.GetSettings().StallDetectOn
	if b == v {
		return true
	}
	this.sendCommands("o", "d")
	n := this.GetSettings().StallDetectOn
	// log.Println("SetStallDetectOn", n)
	return n == b
}

func (this *Zilla) SetBatteryLightPolarity(b bool) bool {
	v := this.GetSettings().BatteryLightPolarity
	if b == v {
		return true
	}
	this.sendCommands("o", "e")
	n := this.GetSettings().BatteryLightPolarity
	// log.Println("SetBatteryLightPolarity", n)
	return n == b
}

func (this *Zilla) SetCheckEngineLightPolarity(b bool) bool {
	v := this.GetSettings().CheckEngineLightPolarity
	if b == v {
		return true
	}
	this.sendCommands("o", "f")
	n := this.GetSettings().CheckEngineLightPolarity
	// log.Println("SetCheckEngineLightPolarity", n)
	return n == b
}

func (this *Zilla) SetReversingContactors(b bool) bool {
	v := this.GetSettings().ReversingContactors
	if b == v {
		return true
	}
	this.sendCommands("o", "g")
	n := this.GetSettings().ReversingContactors
	// log.Println("SetReversingContactors", n)
	return n == b
}

func (this *Zilla) SetSeriesParallelContactors(b bool) bool {
	v := this.GetSettings().SeriesParallelContactors
	if b == v {
		return true
	}
	this.sendCommands("o", "h")
	n := this.GetSettings().SeriesParallelContactors
	// log.Println("SetSeriesParallelContactors", n)
	return n == b
}

func (this *Zilla) SetForceParallelInReverse(b bool) bool {
	v := this.GetSettings().ForceParallelInReverse
	if b == v {
		return true
	}
	this.sendCommands("o", "i")
	n := this.GetSettings().ForceParallelInReverse
	// log.Println("SetForceParallelInReverse", n)
	return n == b
}

func (this *Zilla) SetInhibitSeriesParallelShifting(b bool) bool {
	v := this.GetSettings().InhibitSeriesParallelShifting
	if b == v {
		return true
	}
	this.sendCommands("o", "j")
	n := this.GetSettings().InhibitSeriesParallelShifting
	// log.Println("SetInhibitSeriesParallelShifting", n)
	return n == b
}

func (this *Zilla) SetTachometerDisplayMotorAmps(b bool) bool {
	v := this.GetSettings().TachometerDisplayMotorAmps
	if b == v {
		return true
	}
	this.sendCommands("o", "k")
	n := this.GetSettings().TachometerDisplayMotorAmps
	// log.Println("SetTachometerDisplayMotorAmps", n)
	return n == b
}

func (this *Zilla) SetTachometerSixCylinders(b bool) bool {
	v := this.GetSettings().TachometerSixCylinders
	if b == v {
		return true
	}
	this.sendCommands("o", "l")
	n := this.GetSettings().TachometerSixCylinders
	// log.Println("SetTachometerSixCylinders", n)
	return n == b
}

func (this *Zilla) SetReversesPlugInInputPolarity(b bool) bool {
	v := this.GetSettings().ReversesPlugInInputPolarity
	if b == v {
		return true
	}
	this.sendCommands("o", "m")
	n := this.GetSettings().ReversesPlugInInputPolarity
	// log.Println("SetReversesPlugInInputPolarity", n)
	return n == b
}

func (this *Zilla) SetActivateHEPI(b bool) bool {
	v := this.GetSettings().ActivateHEPI
	if b == v {
		return true
	}
	this.sendCommands("o", "n")
	n := this.GetSettings().ActivateHEPI
	// log.Println("SetActivateHEPI", n)
	return n == b
}

func (this *Zilla) SetIsZ1k(b bool) bool {
	v := this.GetSettings().IsZ1k
	if b == v {
		return true
	}
	this.sendCommands("o", "p")
	n := this.GetSettings().IsZ1k
	// log.Println("SetIsZ1k", n)
	return n == b
}
