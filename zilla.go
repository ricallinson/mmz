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
	closed       bool
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
	if err := this.createLogFile(); err != nil {
		return nil, err
	}
	// Start listening on the queue for commands.
	go this.start()
	// Return the Zilla instance.
	return this, nil
}

// Loop on the queue channel sending commands from the queue to the Hairball.
// If there are no commands in the queue log data until a command appears.
func (this *Zilla) start() {
	for {
		select {
		case cmd := <-this.queue:
			for _, b := range cmd.bytes {
				this.writeBytes(b)
				cmd.data = this.readBytes()
			}
			cmd.done <- true
		default:
			// Once Close() has been called there is nothing more to do.
			if this.closed {
				return
			}
			this.writeLogToFile()
		}
	}
}

// Takes a list of commands and sends them to the Zilla in the order received.
func (this *Zilla) writeCommands(types ...interface{}) []byte {
	cmd := newZillaCommand()
	for _, t := range types {
		switch t.(type) {
		case string:
			cmd.sendString(t.(string))
		case int:
			cmd.sendInt(t.(int))
		default:
			log.Println("Cannot only write types string or int")
			return nil
		}
	}
	// We have to stop logging so the command can be picked up.
	this.writeLog = false
	this.queue <- cmd
	<-cmd.done
	return cmd.data
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
func (this *Zilla) readBytes() []byte {
	return this.readBytesTo(0)
}

// Reads bytes directly from the Zilla up to the given byte or EOF.
func (this *Zilla) readBytesTo(to byte) []byte {
	if this.closed {
		log.Print("Read connection to Zilla has been closed.")
		return nil
	}
	// TODO: The mock only returns a byte array not single bytes.
	buff := make([]byte, 1)
	data := make([]byte, 0)
	for {
		i, _ := this.serialPort.Read(buff)
		data = append(data, buff[0])
		if to == buff[0] || i == 0 {
			break
		}
		log.Print(buff[0])
	}
	data = bytes.TrimSpace(data)
	// log.Println(string(data))
	return data
}

// Blocks while writing log file.
func (this *Zilla) writeLogToFile() {
	// Before anything else happens tell the world we are now logging.
	this.writeLog = true
	// We have to write and read each command to be sure the Zilla is responding.
	// The writeCommands() function is not used as it creates a circular dependency.
	this.writeBytes([]byte{27})
	this.readBytes()
	this.writeBytes([]byte{27})
	this.readBytes()
	this.writeBytes([]byte{27})
	this.readBytes()
	this.writeBytes([]byte("p")) // Menu Special
	this.readBytes()
	this.writeBytes([]byte("Q1\r")) // Start logs
	for this.writeLog {
		line := this.readBytesTo('\n')
		// log.Print(string(line))
		logLine := ParseQ1LineFromHairball(line)
		if logLine == nil {
			continue
		}
		if _, err := this.writeLogFile.Write(logLine.ToBytes()); err != nil {
			// If there is a write error it means the log file has been closed.
			this.writeLog = false
			return
		}
		// Sleep for 100ms as the logs are only written 10 times a second.
		time.Sleep(100 * time.Millisecond)
	}
}

// Creates a new log file for this instance to write to.
func (this *Zilla) createLogFile() error {
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

// Ends the connection to the Zilla and stops all logging.
// It's good practice to call Close() once the instance is no longer needed.
// Once called the Zilla instance can no longer be use to send commands.
func (this *Zilla) Close() {
	// Wait for the queue to drain before closing.
	if len(this.queue) > 0 {
		time.Sleep(time.Millisecond)
	}
	this.closed = true
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
	v := this.GetSettings().BatteryAmpLimit
	this.writeCommands("b", "a", val)
	n := this.GetSettings().BatteryAmpLimit
	log.Println("SetBatteryAmpLimit", n)
	return n != v
}

func (this *Zilla) SetLowBatteryVoltageLimit(val int) bool {
	v := this.GetSettings().LowBatteryVoltageLimit
	this.writeCommands("b", "v", val)
	n := this.GetSettings().LowBatteryVoltageLimit
	log.Println("SetLowBatteryVoltageLimit", n)
	return n != v
}

func (this *Zilla) SetLowBatteryVoltageIndicator(val int) bool {
	v := this.GetSettings().LowBatteryVoltageIndicator
	this.writeCommands("b", "i", val)
	n := this.GetSettings().LowBatteryVoltageIndicator
	log.Println("SetLowBatteryVoltageIndicator", n)
	return n != v
}

// Motor Menu

func (this *Zilla) SetNormalMotorAmpLimit(val int) bool {
	v := this.GetSettings().NormalMotorAmpLimit
	this.writeCommands("m", "a", val)
	n := this.GetSettings().NormalMotorAmpLimit
	log.Println("SetNormalMotorAmpLimit", n)
	return n != v
}

func (this *Zilla) SetSeriesMotorVoltageLimit(val int) bool {
	v := this.GetSettings().SeriesMotorVoltageLimit
	this.writeCommands("m", "v", val)
	n := this.GetSettings().SeriesMotorVoltageLimit
	log.Println("SetSeriesMotorVoltageLimit", n)
	return n != v
}

func (this *Zilla) SetReverseMotorAmpLimit(val int) bool {
	v := this.GetSettings().ReverseMotorAmpLimit
	this.writeCommands("m", "i", val)
	n := this.GetSettings().ReverseMotorAmpLimit
	log.Println("SetReverseMotorAmpLimit", n)
	return n != v
}

func (this *Zilla) SetReverseMotorVoltageLimit(val int) bool {
	v := this.GetSettings().ReverseMotorVoltageLimit
	this.writeCommands("m", "r", val)
	n := this.GetSettings().ReverseMotorVoltageLimit
	log.Println("SetReverseMotorVoltageLimit", n)
	return n != v
}

func (this *Zilla) SetParallelMotorAmpLimit(val int) bool {
	v := this.GetSettings().ParallelMotorAmpLimit
	this.writeCommands("m", "c", val)
	n := this.GetSettings().ParallelMotorAmpLimit
	log.Println("SetParallelMotorAmpLimit", n)
	return n != v
}

func (this *Zilla) SetParallelMotorVoltageLimit(val int) bool {
	v := this.GetSettings().ParallelMotorVoltageLimit
	this.writeCommands("m", "p", val)
	n := this.GetSettings().ParallelMotorVoltageLimit
	log.Println("SetParallelMotorVoltageLimit", n)
	return n != v
}

// Speed Menu

func (this *Zilla) SetForwardRpmLimit(val int) bool {
	this.writeCommands("s", "l", val)
	n := this.GetSettings().ForwardRpmLimit
	log.Println("SetForwardRpmLimit", n)
	return n == val
}

func (this *Zilla) SetReverseRpmLimit(val int) bool {
	this.writeCommands("s", "r", val)
	n := this.GetSettings().ReverseRpmLimit
	log.Println("SetReverseRpmLimit", n)
	return n == val
}

func (this *Zilla) SetMaxRpmLimit(val int) bool {
	this.writeCommands("s", "x", val)
	n := this.GetSettings().MaxRpmLimit
	log.Println("SetMaxRpmLimit", n)
	return n == val
}

// Options Menu

func (this *Zilla) ToggleRpmSensorMotorOne() bool {
	v := this.GetSettings().RpmSensorMotorOne
	this.writeCommands("o", "a")
	n := this.GetSettings().RpmSensorMotorOne
	log.Println("ToggleRpmSensorMotorOne", n)
	return n != v
}

func (this *Zilla) ToggleRpmSensorMotorTwo() bool {
	v := this.GetSettings().RpmSensorMotorTwo
	this.writeCommands("o", "b")
	n := this.GetSettings().RpmSensorMotorTwo
	log.Println("ToggleRpmSensorMotorTwo", n)
	return n != v
}

func (this *Zilla) ToggleAutoShiftingSeriesToParallel() bool {
	v := this.GetSettings().AutoShiftingSeriesToParallel
	this.writeCommands("o", "c")
	n := this.GetSettings().AutoShiftingSeriesToParallel
	log.Println("ToggleAutoShiftingSeriesToParallel", n)
	return n != v
}

func (this *Zilla) ToggleStallDetectOn() bool {
	v := this.GetSettings().StallDetectOn
	this.writeCommands("o", "d")
	n := this.GetSettings().StallDetectOn
	log.Println("ToggleStallDetectOn", n)
	return n != v
}

func (this *Zilla) ToggleBatteryLightPolarity() bool {
	v := this.GetSettings().BatteryLightPolarity
	this.writeCommands("o", "e")
	n := this.GetSettings().BatteryLightPolarity
	log.Println("ToggleBatteryLightPolarity", n)
	return n != v
}

func (this *Zilla) ToggleCheckEngineLightPolarity() bool {
	v := this.GetSettings().CheckEngineLightPolarity
	this.writeCommands("o", "f")
	n := this.GetSettings().CheckEngineLightPolarity
	log.Println("ToggleCheckEngineLightPolarity", n)
	return n != v
}

func (this *Zilla) ToggleReversingContactors() bool {
	v := this.GetSettings().ReversingContactors
	this.writeCommands("o", "g")
	n := this.GetSettings().ReversingContactors
	log.Println("ToggleReversingContactors", n)
	return n != v
}

func (this *Zilla) ToggleSeriesParallelContactors() bool {
	v := this.GetSettings().SeriesParallelContactors
	this.writeCommands("o", "h")
	n := this.GetSettings().SeriesParallelContactors
	log.Println("ToggleSeriesParallelContactors", n)
	return n != v
}

func (this *Zilla) ToggleForceParallelInReverse() bool {
	v := this.GetSettings().ForceParallelInReverse
	this.writeCommands("o", "i")
	n := this.GetSettings().ForceParallelInReverse
	log.Println("ToggleForceParallelInReverse", n)
	return n != v
}

func (this *Zilla) ToggleInhibitSeriesParallelShifting() bool {
	v := this.GetSettings().InhibitSeriesParallelShifting
	this.writeCommands("o", "j")
	n := this.GetSettings().InhibitSeriesParallelShifting
	log.Println("ToggleInhibitSeriesParallelShifting", n)
	return n != v
}

func (this *Zilla) ToggleTachometerDisplayMotorAmps() bool {
	v := this.GetSettings().TachometerDisplayMotorAmps
	this.writeCommands("o", "k")
	n := this.GetSettings().TachometerDisplayMotorAmps
	log.Println("ToggleTachometerDisplayMotorAmps", n)
	return n != v
}

func (this *Zilla) ToggleTachometerSixCylinders() bool {
	v := this.GetSettings().TachometerSixCylinders
	this.writeCommands("o", "l")
	n := this.GetSettings().TachometerSixCylinders
	log.Println("ToggleTachometerSixCylinders", n)
	return n != v
}

func (this *Zilla) ToggleReversesPlugInInputPolarity() bool {
	v := this.GetSettings().ReversesPlugInInputPolarity
	this.writeCommands("o", "m")
	n := this.GetSettings().ReversesPlugInInputPolarity
	log.Println("ToggleReversesPlugInInputPolarity", n)
	return n != v
}

func (this *Zilla) ToggleActivateHEPI() bool {
	v := this.GetSettings().ActivateHEPI
	this.writeCommands("o", "n")
	n := this.GetSettings().ActivateHEPI
	log.Println("ToggleActivateHEPI", n)
	return n != v
}

func (this *Zilla) ToggleIsZ2k() bool {
	v := this.GetSettings().IsZ2k
	this.writeCommands("o", "p")
	n := this.GetSettings().IsZ2k
	log.Println("ToggleIsZ2k", n)
	return n != v
}
