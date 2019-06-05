# mmz

[![Build Status](https://travis-ci.org/ricallinson/mmz.svg?branch=master)](https://travis-ci.org/ricallinson/mmz) [![Build status](https://ci.appveyor.com/api/projects/status/6v17dsgd08n8ieq7/branch/master?svg=true)](https://ci.appveyor.com/project/ricallinson/mmz/branch/master)

__UNSTABLE__

Command line interface to configure and log data for a Manzanita Micro Zilla controller.

## Usage

Download the executable for your chosen platform from the [releases](https://github.com/ricallinson/mmz/releases/tag/v1.0) page. Feel free to rename the executable to `mmz` or `mmz.exe` depending on your chosen platform.

You will need to know the location of the USB serial port which the dongle is plugged into. The [MK3 Digital Perl Scanner Software](http://www.manzanitamicro.com/downloads/category/5-bms2?download=93%3Aperlscanner) documentation describes how to find this for Windows as a COM port number. For Unix based systems you can use [dmesg | grep tty](https://www.cyberciti.biz/faq/find-out-linux-serial-ports-with-setserial/) as described in the link.

## Examples

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -settings

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -realtime

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -cmd ./samples/set_settings.yaml

## Options

### Dongle Location (required)

The path to the USB port where the Hairball is connected.

    mmz -dongle /dev/tty.usbserial-A904RBQ7

### Path to Commands File

The path to the file containing the commands to execute against the Zilla.

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -cmd ./samples/set_settings.yaml

### Send Raw Command

Send a command as detailed in the [Zilla DC Motor Controller and Hairball 2 Manual](http://www.manzanitamicro.com/downloads/category/1-zilla?download=92%3Ahb202d).

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -raw "XXX"

### Settings

Prints the current setting values applied to the controller.

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -settings

Outputs YAML to `stdout` with the following structure;

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

### Realtime

Prints the current state of the controller every 100ms.

    mmz -dongle /dev/tty.usbserial-A904RBQ7 -realtime

Outputs YAML to `stdout` with the following structure;

    Timestamp                      int64
    RxCtrlFlagByte                 int
    AverageCurrentOnMotor          int
    AvailableCurrentFromController int
    ArmDC                          int
    BatteryVoltage                 int
    MotorVoltage                   int
    ControllerTemp                 int
    SpiErrorCount                  int
    CurrentError                   string
    OperatingStatus                int
    MotorKilowatts                 float64
    StoppedState                   bool
    ShiftingInProgress             bool
    MainContactorIsOn              bool
    MotorContactorsAreOn           bool
    DirectionIsReverse             bool
    DirectionIsForward             bool
    MotorsAreInParallel            bool
    MotorsAreInSeries              bool
    MainContactorHasVoltageDrop    bool

## Setup Development Environment

Requires a [Go](https://golang.org/dl/) environment.

    git clone git@github.com:ricallinson/mmz.git $GOPATH/src/git@github.com/ricallinson/mmz
    cd $GOPATH/src/git@github.com/ricallinson/mmz
    go get ./...
    go install

## Test

    cd $GOPATH/src/git@github.com/ricallinson/mmz
    go test

## Code Coverage Report

    cd $GOPATH/src/git@github.com/ricallinson/mmz
    go test -covermode=count -coverprofile=count.out; go tool cover -html=count.out

## Help

### Connecting to Hairball Directly

On a UNIX based system you can directly connect to the hairball from a terminal session.

    screen /dev/tty.usbserial

Exit from Hairball `CTRL + A + \`.
