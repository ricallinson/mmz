package main

import (
	"flag"
	"github.com/tarm/serial"
	"log"
	"time"
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
		serialPort = NewMockPort()
	} else {
		serialError, serialPort = connectToHairball(hairball)
	}
	if serialError != nil {
		log.Println(serialError)
		return
	}
	zilla, zillaError := NewZilla(serialPort)
	if zillaError != nil {
		log.Println(zillaError)
		return
	}
	StartWebServer(port, zilla)
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
