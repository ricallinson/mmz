package main

import (
	. "github.com/ricallinson/simplebdd"
	"reflect"
	"testing"
)

func TestExecutor(t *testing.T) {

	var e *Executor

	BeforeEach(func() {
		zilla, _ := NewZilla(NewMockPort())
		e = NewExecutor(zilla)
	})

	AfterEach(func() {

	})

	Describe("Executor()", func() {

		It("should return an Executor object", func() {
			AssertEqual(reflect.TypeOf(e).String(), "*main.Executor")
		})

		It("should populate the ExecutorCommands object from YAML", func() {
			e.Commands = readYamlFileToExecutorCommands("./fixtures/all_commands.yaml")
			AssertEqual(e.Commands.SetBatteryAmpLimit, 1800)
			AssertEqual(e.Commands.SetLowBatteryVoltageLimit, 119)
			AssertEqual(e.Commands.SetLowBatteryVoltageIndicator, 145)
			AssertEqual(e.Commands.SetNormalMotorAmpLimit, 1600)
			AssertEqual(e.Commands.SetSeriesMotorVoltageLimit, 429)
			AssertEqual(e.Commands.SetReverseMotorAmpLimit, 700)
			AssertEqual(e.Commands.SetReverseMotorVoltageLimit, 106)
			AssertEqual(e.Commands.SetParallelMotorAmpLimit, 2000)
			AssertEqual(e.Commands.SetParallelMotorVoltageLimit, 180)
			AssertEqual(e.Commands.SetForwardRpmLimit, 7000)
			AssertEqual(e.Commands.SetReverseRpmLimit, 1500)
			AssertEqual(e.Commands.SetMaxRpmLimit, 8000)
			AssertEqual(e.Commands.SetRpmSensorMotorOne, true)
			AssertEqual(e.Commands.SetRpmSensorMotorTwo, true)
			AssertEqual(e.Commands.SetAutoShiftingSeriesToParallel, true)
			AssertEqual(e.Commands.SetStallDetectOn, true)
			AssertEqual(e.Commands.SetBatteryLightPolarity, false)
			AssertEqual(e.Commands.SetCheckEngineLightPolarity, true)
			AssertEqual(e.Commands.SetReversingContactors, true)
			AssertEqual(e.Commands.SetSeriesParallelContactors, true)
			AssertEqual(e.Commands.SetForceParallelInReverse, false)
			AssertEqual(e.Commands.SetInhibitSeriesParallelShifting, false)
			AssertEqual(e.Commands.SetTachometerDisplayMotorAmps, false)
			AssertEqual(e.Commands.SetTachometerSixCylinders, false)
			AssertEqual(e.Commands.SetReversesPlugInInputPolarity, false)
			AssertEqual(e.Commands.SetActivateHEPI, false)
			AssertEqual(e.Commands.SetIsZ1k, true)
		})

		It("should return given value from SetBatteryAmpLimit", func() {
			e.Commands.SetBatteryAmpLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.BatteryAmpLimit, 100)
		})

		It("should return given value from SetLowBatteryVoltageLimit", func() {
			e.Commands.SetLowBatteryVoltageLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.LowBatteryVoltageLimit, 100)
		})

		It("should return given value from SetLowBatteryVoltageIndicator", func() {
			e.Commands.SetLowBatteryVoltageIndicator = 100
			r := e.ExecuteCommands()
			AssertEqual(r.LowBatteryVoltageIndicator, 100)
		})

		It("should return given value from SetNormalMotorAmpLimit", func() {
			e.Commands.SetNormalMotorAmpLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.NormalMotorAmpLimit, 100)
		})

		It("should return given value from SetSeriesMotorVoltageLimit", func() {
			e.Commands.SetSeriesMotorVoltageLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.SeriesMotorVoltageLimit, 100)
		})

		It("should return given value from SetReverseMotorAmpLimit", func() {
			e.Commands.SetReverseMotorAmpLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.ReverseMotorAmpLimit, 100)
		})

		It("should return given value from SetReverseMotorVoltageLimit", func() {
			e.Commands.SetReverseMotorVoltageLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.ReverseMotorVoltageLimit, 100)
		})

		It("should return given value from SetParallelMotorAmpLimit", func() {
			e.Commands.SetParallelMotorAmpLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.ParallelMotorAmpLimit, 100)
		})

		It("should return given value from SetParallelMotorVoltageLimit", func() {
			e.Commands.SetParallelMotorVoltageLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.ParallelMotorVoltageLimit, 100)
		})

		It("should return given value from SetForwardRpmLimit", func() {
			e.Commands.SetForwardRpmLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.ForwardRpmLimit, 100)
		})

		It("should return given value from SetReverseRpmLimit", func() {
			e.Commands.SetReverseRpmLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.ReverseRpmLimit, 100)
		})

		It("should return given value from SetMaxRpmLimit", func() {
			e.Commands.SetMaxRpmLimit = 100
			r := e.ExecuteCommands()
			AssertEqual(r.MaxRpmLimit, 100)
		})

		It("should return given value from SetRpmSensorMotorOne", func() {
			e.Commands.SetRpmSensorMotorOne = false
			r := e.ExecuteCommands()
			AssertEqual(r.RpmSensorMotorOne, false)
		})

		It("should return given value from SetRpmSensorMotorTwo", func() {
			e.Commands.SetRpmSensorMotorTwo = false
			r := e.ExecuteCommands()
			AssertEqual(r.RpmSensorMotorTwo, false)
		})

		It("should return given value from SetAutoShiftingSeriesToParallel", func() {
			e.Commands.SetAutoShiftingSeriesToParallel = false
			r := e.ExecuteCommands()
			AssertEqual(r.AutoShiftingSeriesToParallel, false)
		})

		It("should return given value from SetStallDetectOn", func() {
			e.Commands.SetStallDetectOn = false
			r := e.ExecuteCommands()
			AssertEqual(r.StallDetectOn, false)
		})

		It("should return given value from SetBatteryLightPolarity", func() {
			e.Commands.SetBatteryLightPolarity = false
			r := e.ExecuteCommands()
			AssertEqual(r.BatteryLightPolarity, false)
		})

		It("should return given value from SetCheckEngineLightPolarity", func() {
			e.Commands.SetCheckEngineLightPolarity = false
			r := e.ExecuteCommands()
			AssertEqual(r.CheckEngineLightPolarity, false)
		})

		It("should return given value from SetReversingContactors", func() {
			e.Commands.SetReversingContactors = false
			r := e.ExecuteCommands()
			AssertEqual(r.ReversingContactors, false)
		})

		It("should return given value from SetSeriesParallelContactors", func() {
			e.Commands.SetSeriesParallelContactors = false
			r := e.ExecuteCommands()
			AssertEqual(r.SeriesParallelContactors, false)
		})

		It("should return given value from SetForceParallelInReverse", func() {
			e.Commands.SetForceParallelInReverse = false
			r := e.ExecuteCommands()
			AssertEqual(r.ForceParallelInReverse, false)
		})

		It("should return given value from SetInhibitSeriesParallelShifting", func() {
			e.Commands.SetInhibitSeriesParallelShifting = false
			r := e.ExecuteCommands()
			AssertEqual(r.InhibitSeriesParallelShifting, false)
		})

		It("should return given value from SetTachometerDisplayMotorAmps", func() {
			e.Commands.SetTachometerDisplayMotorAmps = false
			r := e.ExecuteCommands()
			AssertEqual(r.TachometerDisplayMotorAmps, false)
		})

		It("should return given value from SetTachometerSixCylinders", func() {
			e.Commands.SetTachometerSixCylinders = false
			r := e.ExecuteCommands()
			AssertEqual(r.TachometerSixCylinders, false)
		})

		It("should return given value from SetReversesPlugInInputPolarity", func() {
			e.Commands.SetReversesPlugInInputPolarity = false
			r := e.ExecuteCommands()
			AssertEqual(r.ReversesPlugInInputPolarity, false)
		})

		It("should return given value from SetActivateHEPI", func() {
			e.Commands.SetActivateHEPI = false
			r := e.ExecuteCommands()
			AssertEqual(r.ActivateHEPI, false)
		})

		It("should return given value from SetIsZ1k", func() {
			e.Commands.SetIsZ1k = false
			r := e.ExecuteCommands()
			AssertEqual(r.IsZ1k, false)
		})

	})

	Report(t)
}
