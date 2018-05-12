package main

type Executor struct {
	zilla    *Zilla
	Commands *ExecutorCommands
}

type ExecutorCommands struct {
	SetBatteryAmpLimit               int  `yaml:"SetBatteryAmpLimit"`
	SetLowBatteryVoltageLimit        int  `yaml:"SetLowBatteryVoltageLimit"`
	SetLowBatteryVoltageIndicator    int  `yaml:"SetLowBatteryVoltageIndicator"`
	SetNormalMotorAmpLimit           int  `yaml:"SetNormalMotorAmpLimit"`
	SetSeriesMotorVoltageLimit       int  `yaml:"SetSeriesMotorVoltageLimit"`
	SetReverseMotorAmpLimit          int  `yaml:"SetReverseMotorAmpLimit"`
	SetReverseMotorVoltageLimit      int  `yaml:"SetReverseMotorVoltageLimit"`
	SetParallelMotorAmpLimit         int  `yaml:"SetParallelMotorAmpLimit"`
	SetParallelMotorVoltageLimit     int  `yaml:"SetParallelMotorVoltageLimit"`
	SetForwardRpmLimit               int  `yaml:"SetForwardRpmLimit"`
	SetReverseRpmLimit               int  `yaml:"SetReverseRpmLimit"`
	SetMaxRpmLimit                   int  `yaml:"SetMaxRpmLimit"`
	SetRpmSensorMotorOne             bool `yaml:"SetRpmSensorMotorOne"`
	SetRpmSensorMotorTwo             bool `yaml:"SetRpmSensorMotorTwo"`
	SetAutoShiftingSeriesToParallel  bool `yaml:"SetAutoShiftingSeriesToParallel"`
	SetStallDetectOn                 bool `yaml:"SetStallDetectOn"`
	SetBatteryLightPolarity          bool `yaml:"SetBatteryLightPolarity"`
	SetCheckEngineLightPolarity      bool `yaml:"SetCheckEngineLightPolarity"`
	SetReversingContactors           bool `yaml:"SetReversingContactors"`
	SetSeriesParallelContactors      bool `yaml:"SetSeriesParallelContactors"`
	SetForceParallelInReverse        bool `yaml:"SetForceParallelInReverse"`
	SetInhibitSeriesParallelShifting bool `yaml:"SetInhibitSeriesParallelShifting"`
	SetTachometerDisplayMotorAmps    bool `yaml:"SetTachometerDisplayMotorAmps"`
	SetTachometerSixCylinders        bool `yaml:"SetTachometerSixCylinders"`
	SetReversesPlugInInputPolarity   bool `yaml:"SetReversesPlugInInputPolarity"`
	SetActivateHEPI                  bool `yaml:"SetActivateHEPI"`
	SetIsZ2k                         bool `yaml:"SetIsZ2k"`
}

func NewExecutor(zilla *Zilla) *Executor {
	this := &Executor{
		zilla:    zilla,
		Commands: &ExecutorCommands{},
	}
	return this
}

func (this *Executor) ExecuteCommands() *ZillaSettings {
	if this.Commands.SetBatteryAmpLimit > 0 {
		this.zilla.SetBatteryAmpLimit(this.Commands.SetBatteryAmpLimit)
	}
	if this.Commands.SetLowBatteryVoltageLimit > 0 {
		this.zilla.SetLowBatteryVoltageLimit(this.Commands.SetLowBatteryVoltageLimit)
	}
	if this.Commands.SetLowBatteryVoltageIndicator > 0 {
		this.zilla.SetLowBatteryVoltageIndicator(this.Commands.SetLowBatteryVoltageIndicator)
	}
	if this.Commands.SetNormalMotorAmpLimit > 0 {
		this.zilla.SetNormalMotorAmpLimit(this.Commands.SetNormalMotorAmpLimit)
	}
	if this.Commands.SetSeriesMotorVoltageLimit > 0 {
		this.zilla.SetSeriesMotorVoltageLimit(this.Commands.SetSeriesMotorVoltageLimit)
	}
	if this.Commands.SetReverseMotorAmpLimit > 0 {
		this.zilla.SetReverseMotorAmpLimit(this.Commands.SetReverseMotorAmpLimit)
	}
	if this.Commands.SetReverseMotorVoltageLimit > 0 {
		this.zilla.SetReverseMotorVoltageLimit(this.Commands.SetReverseMotorVoltageLimit)
	}
	if this.Commands.SetParallelMotorAmpLimit > 0 {
		this.zilla.SetParallelMotorAmpLimit(this.Commands.SetParallelMotorAmpLimit)
	}
	if this.Commands.SetParallelMotorVoltageLimit > 0 {
		this.zilla.SetParallelMotorVoltageLimit(this.Commands.SetParallelMotorVoltageLimit)
	}
	if this.Commands.SetForwardRpmLimit > 0 {
		this.zilla.SetForwardRpmLimit(this.Commands.SetForwardRpmLimit)
	}
	if this.Commands.SetReverseRpmLimit > 0 {
		this.zilla.SetReverseRpmLimit(this.Commands.SetReverseRpmLimit)
	}
	if this.Commands.SetMaxRpmLimit > 0 {
		this.zilla.SetMaxRpmLimit(this.Commands.SetMaxRpmLimit)
	}
	this.zilla.SetRpmSensorMotorOne(this.Commands.SetRpmSensorMotorOne)
	this.zilla.SetRpmSensorMotorTwo(this.Commands.SetRpmSensorMotorTwo)
	this.zilla.SetAutoShiftingSeriesToParallel(this.Commands.SetAutoShiftingSeriesToParallel)
	this.zilla.SetStallDetectOn(this.Commands.SetStallDetectOn)
	this.zilla.SetBatteryLightPolarity(this.Commands.SetBatteryLightPolarity)
	this.zilla.SetCheckEngineLightPolarity(this.Commands.SetCheckEngineLightPolarity)
	this.zilla.SetReversingContactors(this.Commands.SetReversingContactors)
	this.zilla.SetSeriesParallelContactors(this.Commands.SetSeriesParallelContactors)
	this.zilla.SetForceParallelInReverse(this.Commands.SetForceParallelInReverse)
	this.zilla.SetInhibitSeriesParallelShifting(this.Commands.SetInhibitSeriesParallelShifting)
	this.zilla.SetTachometerDisplayMotorAmps(this.Commands.SetTachometerDisplayMotorAmps)
	this.zilla.SetTachometerSixCylinders(this.Commands.SetTachometerSixCylinders)
	this.zilla.SetReversesPlugInInputPolarity(this.Commands.SetReversesPlugInInputPolarity)
	this.zilla.SetActivateHEPI(this.Commands.SetActivateHEPI)
	this.zilla.SetIsZ2k(this.Commands.SetIsZ2k)
	return this.zilla.GetSettings()
}
