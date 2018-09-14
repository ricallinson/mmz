# mmz

[![Build Status](https://travis-ci.org/ricallinson/mmz.svg?branch=master)](https://travis-ci.org/ricallinson/mmz) [![Build status](https://ci.appveyor.com/api/projects/status/6v17dsgd08n8ieq7/branch/master?svg=true)](https://ci.appveyor.com/project/ricallinson/mmz/branch/master)

__UNSTABLE__

Command line interface for configuring and logging data from a Manzanita Micro Zilla controller.


## Usage

Requires a [Go](https://golang.org/dl/) environment.

    go get github.com/ricallinson/mmz
    go install github.com/ricallinson/mmz

## Examples

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -settings

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -realtime

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -cmd ./samples/set_settings.yaml

## Options

### Dongle Location (required)

The path to the USB port where the Hairball is connected.

    mmz -dongle /dev/tty.usbserial-A904RBQ7

### Path to Commands File

The path to the file containing the commands to execute against the bus.

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -cmd ./samples/set_settings.yaml

### Send Raw Command

Send a command as detailed in the [Zilla DC Motor Controller and Hairball 2 Manual](http://www.manzanitamicro.com/downloads/category/1-zilla?download=92%3Ahb202d).

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -raw "XXX"

### Settings

Prints a YAML object of showing the current settings applied to the controller.

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -settings

### Realtime

Prints a YAML object of showing the current state of the controller every 100ms.

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -realtime

## Setup Development Environment

Requires a [Go](https://golang.org/dl/) environment.

    go get github.com/tarm/serial
    go get gopkg.in/yaml.v2
    go get github.com/ricallinson/simplebdd
    go install

## Test

    go test

## Code Coverage Report

    go test -covermode=count -coverprofile=count.out; go tool cover -html=count.out

## Links

* https://www.stewright.me/2013/05/install-and-run-raspbian-from-a-usb-flash-drive/

## Help

### Connect to Hairball

    screen /dev/tty.usbserial

Exit from Hairball `CTRL + A + \`.

## Install and Run on Raspberry Pi

    sudo apt-get install go
    export GOPATH=$HOME/Library/Go/gocode
    go install github.com/ricallinson/mmz
    $HOME/Library/Go/gocode/src/github.com/ricallinson/mmz
    $HOME/Library/Go/gocode/bin/mmz -hairball /dev/tty.usbserial

Where `/dev/tty.usbserial` is the location of your USB to RS-232 serial port adapter.
