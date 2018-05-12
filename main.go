package main

import (
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	var hairball string
	flag.StringVar(&hairball, "hairball", "", "Serial port that's connected to the Zilla Hairball.")
	var raw string
	flag.StringVar(&raw, "raw", "", "Send a raw command to the Hairball.")
	var commands string
	flag.StringVar(&commands, "cmd", "", "Path to the YAML configuration file of commands to execute.")
	var settings bool
	flag.BoolVar(&settings, "settings", false, "List all settings and their current values.")
	var logging bool
	flag.BoolVar(&logging, "log", false, "Log live Zilla data.")
	var logFile string
	flag.StringVar(&logFile, "log-file", "", "Log live Zilla data to give file.")
	flag.Parse()
	var serialPort SerialPort
	var serialError error
	if hairball == "" {
		serialPort = NewMockPort()
	} else {
		serialError, serialPort = connectToHairball(hairball)
	}
	if serialError != nil {
		log.Println(serialError)
		return
	}
	if logging && logFile == "" {
		logFile = "./logs/" + strconv.FormatInt(time.Now().Unix(), 10) + ".dat"
	}
	zilla, zillaError := NewZilla(serialPort, logFile)
	if zillaError != nil {
		log.Println(zillaError)
		return
	}
	defer zilla.Close()
	// If logging is requested keep running until the process is ended.
	for logging || logFile != "" {
		time.Sleep(100 * time.Millisecond)
	}
	// Process CLI Options.
	if raw != "" {
		sendRawCommand(zilla, raw)
		return
	}
	if settings {
		listSettings(zilla)
		return
	}
	if commands == "" {
		log.Println("You must provide a path to YAML file with the commands to execute")
		return
	}
	// Execute the given command object.
	e := NewExecutor(zilla)
	e.Commands = readYamlFileToExecutorCommands(commands)
	var yaml []byte
	yaml = interfaceToYaml(e.ExecuteCommands())
	// Prints to standard out the result.
	log.Println(string(yaml))
}

func connectToHairball(path string) (error, *serial.Port) {
	c := &serial.Config{
		Name:        path,
		Baud:        9600,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
		ReadTimeout: time.Millisecond,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		return err, nil
	}
	return nil, s
}

// Prints to standard out the result.
func sendRawCommand(zilla *Zilla, s string) {
	fmt.Println(string(zilla.Raw(strings.Split(s, " "))))
}

// Prints to standard out the result.
func listSettings(zilla *Zilla) {
	fmt.Println(string(interfaceToYaml(zilla.GetSettings())))
}
