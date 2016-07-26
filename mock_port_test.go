package main

import (
	"os"
)

type MockPort struct {
	history byte
}

func (this *MockPort) Read(b []byte) (int, error) {
	var file *os.File
	var err error
	switch this.history {
	case 'd': // Display Settings
		file, err = os.Open("./fixtures/settings.txt")
	case 'b': // Battery Menu
		file, err = os.Open("./fixtures/battery.txt")
	case 'm': // Motor Menu
		file, err = os.Open("./fixtures/motor.txt")
	case 's': // Speed Menu
		file, err = os.Open("./fixtures/speed.txt")
	case 'o': // Options Menu
		file, err = os.Open("./fixtures/options.txt")
	case 'p': // Special Menu
		file, err = os.Open("./fixtures/special.txt")
	case 0, 27: // Home Menu
		file, err = os.Open("./fixtures/home.txt")
	}
	if err != nil {
		return 0, err
	}
	return file.Read(b)
}

func (this *MockPort) Write(b []byte) (int, error) {
	switch b[0] {
	case 'd', 'b', 'm', 's', 'o', 'p', 27:
		this.history = b[0]
	}
	return 0, nil
}

func (this *MockPort) Flush() error {
	return nil
}

func (this *MockPort) Close() error {
	return nil
}
