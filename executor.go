package main

type Executor struct {
	zilla    *Zilla
	Commands *ExecutorCommands
}

type ExecutorCommands struct {
	SetBatteryAmpLimit                  int
	SetLowBatteryVoltageLimit           int
	SetLowBatteryVoltageIndicator       int
	SetNormalMotorAmpLimit              int
	SetSeriesMotorVoltageLimit          int
	SetReverseMotorAmpLimit             int
	SetReverseMotorVoltageLimit         int
	SetParallelMotorAmpLimit            int
	SetParallelMotorVoltageLimit        int
	SetForwardRpmLimit                  int
	SetReverseRpmLimit                  int
	SetMaxRpmLimit                      int
	ToggleRpmSensorMotorOne             bool
	ToggleRpmSensorMotorTwo             bool
	ToggleAutoShiftingSeriesToParallel  bool
	ToggleStallDetectOn                 bool
	ToggleBatteryLightPolarity          bool
	ToggleCheckEngineLightPolarity      bool
	ToggleReversingContactors           bool
	ToggleSeriesParallelContactors      bool
	ToggleForceParallelInReverse        bool
	ToggleInhibitSeriesParallelShifting bool
	ToggleTachometerDisplayMotorAmps    bool
	ToggleTachometerSixCylinders        bool
	ToggleReversesPlugInInputPolarity   bool
	ToggleActivateHEPI                  bool
	ToggleIsZ2k                         bool
}

func NewExecutor(zilla *Zilla) *Executor {
	this := &Executor{
		zilla:    zilla,
		Commands: &ExecutorCommands{},
	}
	return this
}

func (this *Executor) Close() {
	this.zilla.Close()
}

func (this *Executor) ExecuteCommands() *ZillaSettings {
	if this.Commands.SetBatteryAmpLimit > 0 {
		this.zilla.SetBatteryAmpLimit(this.Commands.SetBatteryAmpLimit)
	}
	return this.zilla.GetSettings()
}
