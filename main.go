package main

import (
	"flag"
	"fmt"
	"github.com/tarm/serial"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port number to run the service on")
	var hairball string
	flag.StringVar(&hairball, "hairball", "", "Serial port that's connected to the Zilla Hairball")
	flag.Parse()
	var serialPort SerialPort
	var serialError error
	if hairball == "" {
		serialPort = &MockPort{}
	} else {
		serialError, serialPort = connectToHairball(hairball)
	}
	if serialError != nil {
		fmt.Println(serialError)
		return
	}
	zilla, zillaError := CreateZilla(serialPort)
	if zillaError != nil {
		fmt.Println(zillaError)
		return
	}
	StartWebServer(port, zilla)
}

func connectToHairball(path string) (error, *serial.Port) {
	c := &serial.Config{
		Name:     path,
		Baud:     9600,
		Size:     8,
		Parity:   serial.ParityNone,
		StopBits: serial.Stop1,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		return err, nil
	}
	return nil, s
}
