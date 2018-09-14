# mmz

[![Build Status](https://travis-ci.org/ricallinson/mmz.svg?branch=master)](https://travis-ci.org/ricallinson/mmz) [![Build status](https://ci.appveyor.com/api/projects/status/6v17dsgd08n8ieq7/branch/master?svg=true)](https://ci.appveyor.com/project/ricallinson/mmz/branch/master)

__UNSTABLE__

Command line interface for configuring and logging data from a Manzanita Micro Zilla controller.

## Install and Run on Raspberry Pi

    sudo apt-get install go
    export GOPATH=$HOME/Library/Go/gocode
    go install github.com/ricallinson/mmz
    $HOME/Library/Go/gocode/src/github.com/ricallinson/mmz
    $HOME/Library/Go/gocode/bin/mmz -hairball /dev/tty.usbserial

Where `/dev/tty.usbserial` is the location of your USB to RS-232 serial port adapter.

## Usage

You have to be in the directory where Go installed mmz;

    cd $GOPATH/src/github.com/ricallinson/mmz

## Setup Development Environment

Requires a [Go](https://golang.org/dl/) environment.

    go get github.com/ricallinson/mmz
    go get github.com/tarm/serial
    go get github.com/ricallinson/simplebdd
    go install github.com/ricallinson/mmz

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
