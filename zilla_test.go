package main

import (
	. "github.com/ricallinson/simplebdd"
	"testing"
)

func TestZilla(t *testing.T) {

	Describe("Zilla()", func() {

		It("should read home", func() {
			z := CreateZilla(&MockPort{})
			AssertEqual(z.menuHome(), true)
		})

		It("should read settings", func() {
			z := CreateZilla(&MockPort{})
			AssertEqual(z.menuSettings(), true)
		})

		It("should read battery", func() {
			z := CreateZilla(&MockPort{})
			AssertEqual(z.menuBattery(), true)
		})

		It("should read motor", func() {
			z := CreateZilla(&MockPort{})
			AssertEqual(z.menuMotor(), true)
		})

		It("should read speed", func() {
			z := CreateZilla(&MockPort{})
			AssertEqual(z.menuSpeed(), true)
		})

		It("should read options", func() {
			z := CreateZilla(&MockPort{})
			AssertEqual(z.menuOptions(), true)
		})

		It("should read special", func() {
			z := CreateZilla(&MockPort{})
			AssertEqual(z.menuSpecial(), true)
		})

		It("should execute Refresh", func() {
			z := CreateZilla(&MockPort{})
			z.Refresh()
			AssertEqual(z.BatteryAmpLimit, 1800)
			AssertEqual(z.LowBatteryVoltageLimit, 119)
			AssertEqual(z.LowBatteryVoltageIndicator, 145)
			AssertEqual(z.NormalMotorAmpLimit, 1600)
			AssertEqual(z.SeriesMotorVoltageLimit, 429)
			AssertEqual(z.ReverseMotorAmpLimit, 700)
			AssertEqual(z.ReverseMotorVoltageLimit, 106)
			AssertEqual(z.ParallelMotorAmpLimit, 2000)
			AssertEqual(z.ParallelMotorVoltageLimit, 180)
			AssertEqual(z.ForwardRpmLimit, 7000)
			AssertEqual(z.ReverseRpmLimit, 1500)
			AssertEqual(z.MaxRpmLimit, 8000)
			AssertEqual(z.RpmSensorMotorOne, true)
			AssertEqual(z.RpmSensorMotorTwo, true)
			AssertEqual(z.AutoShiftingSeriesToParallel, true)
			AssertEqual(z.StallDetectOn, true)
			AssertEqual(z.BatteryLightPolarity, false)
			AssertEqual(z.CheckEngineLightPolarity, true)
			AssertEqual(z.ReversingContactors, true)
			AssertEqual(z.SeriesParallelContactors, true)
			AssertEqual(z.ForceParallelInReverse, false)
			AssertEqual(z.InhibitSeriesParallelShifting, false)
			AssertEqual(z.TachometerDisplayMotorAmps, false)
			AssertEqual(z.TachometerSixCylinders, false)
			AssertEqual(z.ReversesPlugInInputPolarity, false)
			AssertEqual(z.ActivateHEPI, false)
			AssertEqual(z.notUsed, false)
			AssertEqual(z.IsZ2k, true)
			AssertEqual(z.CurrentState, 1311)
			AssertEqual(z.Errors[0], "1111")
			AssertEqual(z.Errors[1], "1111")
			AssertEqual(z.Errors[2], "1111")
			AssertEqual(z.Errors[3], "1111")
			AssertEqual(z.Errors[4], "1111")
		})
	})

	Report(t)
}
