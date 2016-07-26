package main

import(
    // "fmt"
)

type Zilla struct {
	BatteryAmpLimit               int  // a)BA
	LowBatteryVoltageLimit        int  // v)LBV
	LowBatteryVoltageIndicator    int  // i)LBVI
	NormalMotorAmpLimit           int  // a) Amp
	SeriesMotorVoltageLimit       int  // v) Volt
	ReverseMotorAmpLimit          int  // i) RA
	ReverseMotorVoltageLimit      int  // r) RV
	ParallelMotorAmpLimit         int  // c) PA
	ParallelMotorVoltageLimit     int  // p) PV
	ForwardRpmLimit               int  // l)Norm
	ReverseRpmLimit               int  // r)Rev
	MaxRpmLimit                   int  // x)Max
	RpmSensorMotorOne             bool // a) On
	RpmSensorMotorTwo             bool // b) On
	AutoShiftingSeriesToParallel  bool // c) On
	StallDetectOn                 bool // d) On
	BatteryLightPolarity          bool // e) Off
	CheckEngineLightPolarity      bool // f) On
	ReversingContactors           bool // g) On
	SeriesParallelContactors      bool // h) On
	ForceParallelInReverse        bool // i) Off
	InhibitSeriesParallelShifting bool // j) Off
	TachometerDisplayMotorAmps    bool // k) Off
	TachometerSixCylinders        bool // l) Off
	ReversesPlugInInputPolarity   bool // m) Off
	ActivateHEPI                  bool // n) Off
	notUsed                       bool // o) Off
	IsZ2k                         bool // p) Off
	CurrentState                  int
	Errors                        []int
    LastZillaOutput               []byte // byte array of the last Zilla output
}

func CreateZilla() *Zilla {
	z := &Zilla{}
	z.Refresh()
	return z
}

func (this *Zilla) menuHome() bool {
    return false
}

func (this *Zilla) menuSettings() bool {
    if (!this.menuHome()) {
        return false
    }
    return false
}

func (this *Zilla) menuBattery() bool {
    if (!this.menuHome()) {
        return false
    }
    return false
}

func (this *Zilla) menuMotor() bool {
    if (!this.menuHome()) {
        return false
    }
    return false
}

func (this *Zilla) menuSpeed() bool {
    if (!this.menuHome()) {
        return false
    }
    return false
}

func (this *Zilla) menuOptions() bool {
    if (!this.menuHome()) {
        return false
    }
    return false
}

func (this *Zilla) menuSpecial() bool {
    if (!this.menuHome()) {
        return false
    }
    return false
}

// Refreshes all attributes by reading them from the Zilla Controller.
func (this *Zilla) Refresh() bool {
    if (!this.menuSettings()) {
        return false
    }
	return false
}

func (this *Zilla) SetBatteryAmpLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetLowBatteryVoltageLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetNormalMotorAmpLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetSeriesMotorVoltageLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetReverseMotorAmpLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetReverseMotorVoltageLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetParallelMotorAmpLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetParallelMotorVoltageLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetForwardRpmLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetReverseRpmLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetMaxRpmLimit(val int) bool {
	return this.Refresh()
}

func (this *Zilla) SetRpmSensorMotorOne(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetRpmSensorMotorTwo(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetAutoShiftingSeriesToParallel(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetStallDetectOn(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetBatteryLightPolarity(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetCheckEngineLightPolarity(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetReversingContactors(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetSeriesParallelContactors(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetForceParallelInReverse(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetInhibitSeriesParallelShifting(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetTachometerDisplayMotorAmps(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetTachometerSixCylinders(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetReversesPlugInInputPolarity(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetActivateHEPI(val bool) bool {
	return this.Refresh()
}

func (this *Zilla) SetIsZ2k(val bool) bool {
	return this.Refresh()
}
