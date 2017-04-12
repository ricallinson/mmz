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
			zilla.Close()
		})

		It("should return a Zilla object", func() {
			AssertEqual(reflect.TypeOf(zilla).String(), "*main.Zilla")
		})

		It("should execute a command in the Zilla.queue", func() {
			cmd := newZillaCommand()
			cmd.sendBytes([]byte{27})
			zilla.queue <- cmd
			AssertEqual(<-cmd.done, true)
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

		// Toggle

		It("should ToggleRpmSensorMotorOne", func() {
			AssertEqual(zilla.ToggleRpmSensorMotorOne(), true)
		})

		It("should ToggleRpmSensorMotorTwo", func() {
			AssertEqual(zilla.ToggleRpmSensorMotorTwo(), true)
		})

		It("should ToggleAutoShiftingSeriesToParallel", func() {
			AssertEqual(zilla.ToggleAutoShiftingSeriesToParallel(), true)
		})

		It("should ToggleStallDetectOn", func() {
			AssertEqual(zilla.ToggleStallDetectOn(), true)
		})

		It("should ToggleBatteryLightPolarity", func() {
			AssertEqual(zilla.ToggleBatteryLightPolarity(), true)
		})

		It("should ToggleCheckEngineLightPolarity", func() {
			AssertEqual(zilla.ToggleCheckEngineLightPolarity(), true)
		})

		It("should ToggleReversingContactors", func() {
			AssertEqual(zilla.ToggleReversingContactors(), true)
		})

		It("should ToggleSeriesParallelContactors", func() {
			AssertEqual(zilla.ToggleSeriesParallelContactors(), true)
		})

		It("should ToggleForceParallelInReverse", func() {
			AssertEqual(zilla.ToggleForceParallelInReverse(), true)
		})

		It("should ToggleInhibitSeriesParallelShifting", func() {
			AssertEqual(zilla.ToggleInhibitSeriesParallelShifting(), true)
		})

		It("should ToggleTachometerDisplayMotorAmps", func() {
			AssertEqual(zilla.ToggleTachometerDisplayMotorAmps(), true)
		})

		It("should ToggleTachometerSixCylinders", func() {
			AssertEqual(zilla.ToggleTachometerSixCylinders(), true)
		})

		It("should ToggleReversesPlugInInputPolarity", func() {
			AssertEqual(zilla.ToggleReversesPlugInInputPolarity(), true)
		})

		It("should ToggleActivateHEPI", func() {
			AssertEqual(zilla.ToggleActivateHEPI(), true)
		})

		It("should ToggleIsZ2k", func() {
			AssertEqual(zilla.ToggleIsZ2k(), true)
		})

		// Get

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
			time.Sleep(100 * time.Millisecond)
			zilla.GetLiveData()
			time.Sleep(100 * time.Millisecond)
			zilla.GetSettings()
			time.Sleep(100 * time.Millisecond)
			zilla.GetLiveData()
			AssertEqual(true, true)
		})
	})

	Report(t)
}
