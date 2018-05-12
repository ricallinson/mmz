package main

import (
	. "github.com/ricallinson/simplebdd"
	"reflect"
	"testing"
	"time"
)

func TestZilla(t *testing.T) {

	Describe("Zilla()", func() {

		var zilla *Zilla

		BeforeEach(func() {
			zilla, _ = NewZilla(NewMockPort())
		})

		AfterEach(func() {

		})

		It("should return a Zilla object", func() {
			AssertEqual(reflect.TypeOf(zilla).String(), "*main.Zilla")
		})

		It("should execute a command in the Zilla.queue", func() {
			cmd := newZillaCommand()
			cmd.sendBytes([]byte{27})
			zilla.sendCommand(cmd)
			AssertEqual(true, true)
		})

		// Set

		It("should SetBatteryAmpLimit to 999", func() {
			AssertEqual(zilla.SetBatteryAmpLimit(999), true)
			AssertEqual(zilla.GetSettings().BatteryAmpLimit, 999)
		})

		It("should SetLowBatteryVoltageLimit to 999", func() {
			AssertEqual(zilla.SetLowBatteryVoltageLimit(999), true)
			AssertEqual(zilla.GetSettings().LowBatteryVoltageLimit, 999)
		})

		It("should SetLowBatteryVoltageIndicator to 999", func() {
			AssertEqual(zilla.SetLowBatteryVoltageIndicator(999), true)
			AssertEqual(zilla.GetSettings().LowBatteryVoltageIndicator, 999)
		})

		It("should SetNormalMotorAmpLimit to 999", func() {
			AssertEqual(zilla.SetNormalMotorAmpLimit(999), true)
			AssertEqual(zilla.GetSettings().NormalMotorAmpLimit, 999)
		})

		It("should SetSeriesMotorVoltageLimit to 999", func() {
			AssertEqual(zilla.SetSeriesMotorVoltageLimit(999), true)
			AssertEqual(zilla.GetSettings().SeriesMotorVoltageLimit, 999)
		})

		It("should SetReverseMotorAmpLimit to 999", func() {
			AssertEqual(zilla.SetReverseMotorAmpLimit(999), true)
			AssertEqual(zilla.GetSettings().ReverseMotorAmpLimit, 999)
		})

		It("should SetReverseMotorVoltageLimit to 999", func() {
			AssertEqual(zilla.SetReverseMotorVoltageLimit(999), true)
			AssertEqual(zilla.GetSettings().ReverseMotorVoltageLimit, 999)
		})

		It("should SetParallelMotorAmpLimit to 999", func() {
			AssertEqual(zilla.SetParallelMotorAmpLimit(999), true)
			AssertEqual(zilla.GetSettings().ParallelMotorAmpLimit, 999)
		})

		It("should SetParallelMotorVoltageLimit to 999", func() {
			AssertEqual(zilla.SetParallelMotorVoltageLimit(999), true)
			AssertEqual(zilla.GetSettings().ParallelMotorVoltageLimit, 999)
		})

		It("should SetForwardRpmLimit to 999", func() {
			AssertEqual(zilla.SetForwardRpmLimit(999), true)
			AssertEqual(zilla.GetSettings().ForwardRpmLimit, 999)
		})

		It("should SetReverseRpmLimit to 999", func() {
			AssertEqual(zilla.SetReverseRpmLimit(999), true)
			AssertEqual(zilla.GetSettings().ReverseRpmLimit, 999)
		})

		It("should SetMaxRpmLimit to 999", func() {
			AssertEqual(zilla.SetMaxRpmLimit(999), true)
			AssertEqual(zilla.GetSettings().MaxRpmLimit, 999)
		})

		// Set

		It("should SetRpmSensorMotorOne", func() {
			AssertEqual(zilla.SetRpmSensorMotorOne(false), true)
			AssertEqual(zilla.SetRpmSensorMotorOne(true), true)
			AssertEqual(zilla.SetRpmSensorMotorOne(false), true)
			AssertEqual(zilla.SetRpmSensorMotorOne(true), true)
		})

		It("should SetRpmSensorMotorTwo", func() {
			AssertEqual(zilla.SetRpmSensorMotorTwo(false), true)
			AssertEqual(zilla.SetRpmSensorMotorTwo(true), true)
			AssertEqual(zilla.SetRpmSensorMotorTwo(false), true)
			AssertEqual(zilla.SetRpmSensorMotorTwo(true), true)
		})

		It("should SetAutoShiftingSeriesToParallel", func() {
			AssertEqual(zilla.SetAutoShiftingSeriesToParallel(false), true)
			AssertEqual(zilla.SetAutoShiftingSeriesToParallel(true), true)
			AssertEqual(zilla.SetAutoShiftingSeriesToParallel(false), true)
			AssertEqual(zilla.SetAutoShiftingSeriesToParallel(true), true)
		})

		It("should SetStallDetectOn", func() {
			AssertEqual(zilla.SetStallDetectOn(false), true)
			AssertEqual(zilla.SetStallDetectOn(true), true)
			AssertEqual(zilla.SetStallDetectOn(false), true)
			AssertEqual(zilla.SetStallDetectOn(true), true)
		})

		It("should SetBatteryLightPolarity", func() {
			AssertEqual(zilla.SetBatteryLightPolarity(false), true)
			AssertEqual(zilla.SetBatteryLightPolarity(true), true)
			AssertEqual(zilla.SetBatteryLightPolarity(false), true)
			AssertEqual(zilla.SetBatteryLightPolarity(true), true)
		})

		It("should SetCheckEngineLightPolarity", func() {
			AssertEqual(zilla.SetCheckEngineLightPolarity(false), true)
			AssertEqual(zilla.SetCheckEngineLightPolarity(true), true)
			AssertEqual(zilla.SetCheckEngineLightPolarity(false), true)
			AssertEqual(zilla.SetCheckEngineLightPolarity(true), true)
		})

		It("should SetReversingContactors", func() {
			AssertEqual(zilla.SetReversingContactors(false), true)
			AssertEqual(zilla.SetReversingContactors(true), true)
			AssertEqual(zilla.SetReversingContactors(false), true)
			AssertEqual(zilla.SetReversingContactors(true), true)
		})

		It("should SetSeriesParallelContactors", func() {
			AssertEqual(zilla.SetSeriesParallelContactors(false), true)
			AssertEqual(zilla.SetSeriesParallelContactors(true), true)
			AssertEqual(zilla.SetSeriesParallelContactors(false), true)
			AssertEqual(zilla.SetSeriesParallelContactors(true), true)
		})

		It("should SetForceParallelInReverse", func() {
			AssertEqual(zilla.SetForceParallelInReverse(false), true)
			AssertEqual(zilla.SetForceParallelInReverse(true), true)
			AssertEqual(zilla.SetForceParallelInReverse(false), true)
			AssertEqual(zilla.SetForceParallelInReverse(true), true)
		})

		It("should SetInhibitSeriesParallelShifting", func() {
			AssertEqual(zilla.SetInhibitSeriesParallelShifting(false), true)
			AssertEqual(zilla.SetInhibitSeriesParallelShifting(true), true)
			AssertEqual(zilla.SetInhibitSeriesParallelShifting(false), true)
			AssertEqual(zilla.SetInhibitSeriesParallelShifting(true), true)
		})

		It("should SetTachometerDisplayMotorAmps", func() {
			AssertEqual(zilla.SetTachometerDisplayMotorAmps(false), true)
			AssertEqual(zilla.SetTachometerDisplayMotorAmps(true), true)
			AssertEqual(zilla.SetTachometerDisplayMotorAmps(false), true)
			AssertEqual(zilla.SetTachometerDisplayMotorAmps(true), true)
		})

		It("should SetTachometerSixCylinders", func() {
			AssertEqual(zilla.SetTachometerSixCylinders(false), true)
			AssertEqual(zilla.SetTachometerSixCylinders(true), true)
			AssertEqual(zilla.SetTachometerSixCylinders(false), true)
			AssertEqual(zilla.SetTachometerSixCylinders(true), true)
		})

		It("should SetReversesPlugInInputPolarity", func() {
			AssertEqual(zilla.SetReversesPlugInInputPolarity(false), true)
			AssertEqual(zilla.SetReversesPlugInInputPolarity(true), true)
			AssertEqual(zilla.SetReversesPlugInInputPolarity(false), true)
			AssertEqual(zilla.SetReversesPlugInInputPolarity(true), true)
		})

		It("should SetActivateHEPI", func() {
			AssertEqual(zilla.SetActivateHEPI(false), true)
			AssertEqual(zilla.SetActivateHEPI(true), true)
			AssertEqual(zilla.SetActivateHEPI(false), true)
			AssertEqual(zilla.SetActivateHEPI(true), true)
		})

		It("should SetIsZ2k", func() {
			AssertEqual(zilla.SetIsZ2k(false), true)
			AssertEqual(zilla.SetIsZ2k(true), true)
			AssertEqual(zilla.SetIsZ2k(false), true)
			AssertEqual(zilla.SetIsZ2k(true), true)
		})

		// // Get

		It("should execute GetSettings", func() {
			settings := zilla.GetSettings()
			AssertEqual(settings.BatteryAmpLimit, 1800)
			AssertEqual(settings.LowBatteryVoltageLimit, 119)
			AssertEqual(settings.LowBatteryVoltageIndicator, 145)
			AssertEqual(settings.NormalMotorAmpLimit, 1600)
			AssertEqual(settings.SeriesMotorVoltageLimit, 429)
			AssertEqual(settings.ReverseMotorAmpLimit, 700)
			AssertEqual(settings.ReverseMotorVoltageLimit, 106)
			AssertEqual(settings.ParallelMotorAmpLimit, 2000)
			AssertEqual(settings.ParallelMotorVoltageLimit, 180)
			AssertEqual(settings.ForwardRpmLimit, 7000)
			AssertEqual(settings.ReverseRpmLimit, 1500)
			AssertEqual(settings.MaxRpmLimit, 8000)
			AssertEqual(settings.RpmSensorMotorOne, true)
			AssertEqual(settings.RpmSensorMotorTwo, true)
			AssertEqual(settings.AutoShiftingSeriesToParallel, true)
			AssertEqual(settings.StallDetectOn, true)
			AssertEqual(settings.BatteryLightPolarity, false)
			AssertEqual(settings.CheckEngineLightPolarity, true)
			AssertEqual(settings.ReversingContactors, true)
			AssertEqual(settings.SeriesParallelContactors, true)
			AssertEqual(settings.ForceParallelInReverse, false)
			AssertEqual(settings.InhibitSeriesParallelShifting, false)
			AssertEqual(settings.TachometerDisplayMotorAmps, false)
			AssertEqual(settings.TachometerSixCylinders, false)
			AssertEqual(settings.ReversesPlugInInputPolarity, false)
			AssertEqual(settings.ActivateHEPI, false)
			AssertEqual(settings.notUsed, false)
			AssertEqual(settings.IsZ2k, true)
			AssertEqual(settings.CurrentState, "1311")
			AssertEqual(settings.Errors[0], "1111")
			AssertEqual(settings.Errors[1], "1111")
			AssertEqual(settings.Errors[2], "1111")
			AssertEqual(settings.Errors[3], "1111")
			AssertEqual(settings.Errors[4], "1111")
		})

		It("should start logging and then show settings", func() {
			go zilla.StartLogging("./logs/testing.dat")
			time.Sleep(1 * time.Second)
			zilla.StopLogging()
			zilla.GetLiveData()
			AssertEqual(true, true)
		})
	})

	Report(t)
}
