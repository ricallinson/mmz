package main

import (
	. "github.com/ricallinson/simplebdd"
	"testing"
)

func TestZilla(t *testing.T) {

	Describe("Zilla()", func() {

		var zilla *Zilla

		BeforeEach(func() {
			var err error
			zilla, err = NewZilla(NewMockPort())
			if err != nil {
				panic(err)
			}
		})

		AfterEach(func() {
			zilla.CloseLog()
		})

		It("should return a Zilla object", func() {
			zilla, err := NewZilla(NewMockPort())
			zilla.CloseLog()
			AssertEqual(err, nil)
			AssertNotEqual(zilla, nil)
		})

		It("should read home", func() {
			AssertEqual(zilla.menuHome(), true)
		})

		It("should read settings", func() {
			AssertEqual(zilla.menuSettings(), true)
		})

		It("should read battery", func() {
			AssertEqual(zilla.menuBattery(), true)
		})

		It("should read motor", func() {
			AssertEqual(zilla.menuMotor(), true)
		})

		It("should read speed", func() {
			AssertEqual(zilla.menuSpeed(), true)
		})

		It("should read options", func() {
			AssertEqual(zilla.menuOptions(), true)
		})

		It("should read special", func() {
			AssertEqual(zilla.menuSpecial(), true)
		})

		It("should execute Refresh", func() {
			zilla.Refresh()
			AssertEqual(zilla.BatteryAmpLimit, 1800)
			AssertEqual(zilla.LowBatteryVoltageLimit, 119)
			AssertEqual(zilla.LowBatteryVoltageIndicator, 145)
			AssertEqual(zilla.NormalMotorAmpLimit, 1600)
			AssertEqual(zilla.SeriesMotorVoltageLimit, 429)
			AssertEqual(zilla.ReverseMotorAmpLimit, 700)
			AssertEqual(zilla.ReverseMotorVoltageLimit, 106)
			AssertEqual(zilla.ParallelMotorAmpLimit, 2000)
			AssertEqual(zilla.ParallelMotorVoltageLimit, 180)
			AssertEqual(zilla.ForwardRpmLimit, 7000)
			AssertEqual(zilla.ReverseRpmLimit, 1500)
			AssertEqual(zilla.MaxRpmLimit, 8000)
			AssertEqual(zilla.RpmSensorMotorOne, true)
			AssertEqual(zilla.RpmSensorMotorTwo, true)
			AssertEqual(zilla.AutoShiftingSeriesToParallel, true)
			AssertEqual(zilla.StallDetectOn, true)
			AssertEqual(zilla.BatteryLightPolarity, false)
			AssertEqual(zilla.CheckEngineLightPolarity, true)
			AssertEqual(zilla.ReversingContactors, true)
			AssertEqual(zilla.SeriesParallelContactors, true)
			AssertEqual(zilla.ForceParallelInReverse, false)
			AssertEqual(zilla.InhibitSeriesParallelShifting, false)
			AssertEqual(zilla.TachometerDisplayMotorAmps, false)
			AssertEqual(zilla.TachometerSixCylinders, false)
			AssertEqual(zilla.ReversesPlugInInputPolarity, false)
			AssertEqual(zilla.ActivateHEPI, false)
			AssertEqual(zilla.notUsed, false)
			AssertEqual(zilla.IsZ2k, true)
			AssertEqual(zilla.CurrentState, "1311")
			AssertEqual(zilla.Errors[0], "1111")
			AssertEqual(zilla.Errors[1], "1111")
			AssertEqual(zilla.Errors[2], "1111")
			AssertEqual(zilla.Errors[3], "1111")
			AssertEqual(zilla.Errors[4], "1111")
		})

		// Set

		It("should SetBatteryAmpLimit to 999", func() {
			AssertEqual(zilla.SetBatteryAmpLimit(999), true)
			AssertEqual(zilla.BatteryAmpLimit, 999)
		})

		It("should SetLowBatteryVoltageLimit to 999", func() {
			AssertEqual(zilla.SetLowBatteryVoltageLimit(999), true)
			AssertEqual(zilla.LowBatteryVoltageLimit, 999)
		})

		It("should SetLowBatteryVoltageIndicator to 999", func() {
			AssertEqual(zilla.SetLowBatteryVoltageIndicator(999), true)
			AssertEqual(zilla.LowBatteryVoltageIndicator, 999)
		})

		It("should SetNormalMotorAmpLimit to 999", func() {
			AssertEqual(zilla.SetNormalMotorAmpLimit(999), true)
			AssertEqual(zilla.NormalMotorAmpLimit, 999)
		})

		It("should SetSeriesMotorVoltageLimit to 999", func() {
			AssertEqual(zilla.SetSeriesMotorVoltageLimit(999), true)
			AssertEqual(zilla.SeriesMotorVoltageLimit, 999)
		})

		It("should SetReverseMotorAmpLimit to 999", func() {
			AssertEqual(zilla.SetReverseMotorAmpLimit(999), true)
			AssertEqual(zilla.ReverseMotorAmpLimit, 999)
		})

		It("should SetReverseMotorVoltageLimit to 999", func() {
			AssertEqual(zilla.SetReverseMotorVoltageLimit(999), true)
			AssertEqual(zilla.ReverseMotorVoltageLimit, 999)
		})

		It("should SetParallelMotorAmpLimit to 999", func() {
			AssertEqual(zilla.SetParallelMotorAmpLimit(999), true)
			AssertEqual(zilla.ParallelMotorAmpLimit, 999)
		})

		It("should SetParallelMotorVoltageLimit to 999", func() {
			AssertEqual(zilla.SetParallelMotorVoltageLimit(999), true)
			AssertEqual(zilla.ParallelMotorVoltageLimit, 999)
		})

		It("should SetForwardRpmLimit to 999", func() {
			AssertEqual(zilla.SetForwardRpmLimit(999), true)
			AssertEqual(zilla.ForwardRpmLimit, 999)
		})

		It("should SetReverseRpmLimit to 999", func() {
			AssertEqual(zilla.SetReverseRpmLimit(999), true)
			AssertEqual(zilla.ReverseRpmLimit, 999)
		})

		It("should SetMaxRpmLimit to 999", func() {
			AssertEqual(zilla.SetMaxRpmLimit(999), true)
			AssertEqual(zilla.MaxRpmLimit, 999)
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
	})

	Report(t)
}
