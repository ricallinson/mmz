package main

type DataStream struct {
    MotorKilowatts int
    MotorVoltage int
    BatteryVoltage int
    AverageCurrentOnMotor int
    AvailableCurrentFromController int
    ControllerTemp int
    StoppedState bool
    ShiftingInProgress bool
    MainContactorIsOn bool
    MotorContactorsAreOn bool
    DirectionIsReverse bool
    DirectionIsForward bool
    MotorsAreInParallel bool
    MotorsAreInSeries bool
    MainContactorHasVoltageDrop bool
    LatestOperation string
}

func ReadLatestFromDataStream() *DataStream {
    ds := &DataStream{
        MotorVoltage: 38,
        BatteryVoltage: 48,
        AverageCurrentOnMotor: 200,
        AvailableCurrentFromController: 1000,
        ControllerTemp: 100,
        ShiftingInProgress: false,
        MainContactorIsOn: true,
        MotorContactorsAreOn: true,
        DirectionIsReverse: false,
        DirectionIsForward: true,
        MotorsAreInParallel: false,
        MotorsAreInSeries: false,
        MainContactorHasVoltageDrop: false,
        LatestOperation: "1314",
    }
    ds.MotorKilowatts = ds.MotorVoltage * ds.AverageCurrentOnMotor / 1000
    ds.LatestOperation = ds.LatestOperation + ": " + Codes[ds.LatestOperation]
    return ds
}
