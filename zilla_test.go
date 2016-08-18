package main

import (
	. "github.com/ricallinson/simplebdd"
	"testing"
)

func TestZilla(t *testing.T) {

	Describe("Zilla()", func() {

		It("should return a Zilla object", func() {
			z, err := CreateZilla(NewMockPort())
			AssertEqual(err, nil)
			AssertNotEqual(z, nil)
		})

		It("should read home", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.menuHome(), true)
		})

		It("should read settings", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.menuSettings(), true)
		})

		It("should read battery", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.menuBattery(), true)
		})

		It("should read motor", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.menuMotor(), true)
		})

		It("should read speed", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.menuSpeed(), true)
		})

		It("should read options", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.menuOptions(), true)
		})

		It("should read special", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.menuSpecial(), true)
		})

		It("should execute Refresh", func() {
			z, _ := CreateZilla(NewMockPort())
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
			AssertEqual(z.CurrentState, "1311")
			AssertEqual(z.Errors[0], "1111")
			AssertEqual(z.Errors[1], "1111")
			AssertEqual(z.Errors[2], "1111")
			AssertEqual(z.Errors[3], "1111")
			AssertEqual(z.Errors[4], "1111")
		})

		It("should SetBatteryAmpLimit to 999", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.SetBatteryAmpLimit(999), true)
			AssertEqual(z.BatteryAmpLimit, 999)
		})

		It("should SetLowBatteryVoltageLimit to 999", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.SetLowBatteryVoltageLimit(999), true)
			AssertEqual(z.LowBatteryVoltageLimit, 999)
		})

		It("should SetLowBatteryVoltageIndicator to 999", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.SetLowBatteryVoltageIndicator(999), true)
			AssertEqual(z.LowBatteryVoltageIndicator, 999)
		})

		It("should SetNormalMotorAmpLimit to 999", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.SetNormalMotorAmpLimit(999), true)
			AssertEqual(z.NormalMotorAmpLimit, 999)
		})

		It("should SetSeriesMotorVoltageLimit to 999", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.SetSeriesMotorVoltageLimit(999), true)
			AssertEqual(z.SeriesMotorVoltageLimit, 999)
		})

		It("should ToggleRpmSensorMotorOne", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.ToggleRpmSensorMotorOne(), true)
		})

		It("should ToggleActivateHEPI", func() {
			z, _ := CreateZilla(NewMockPort())
			AssertEqual(z.ToggleActivateHEPI(), true)
		})
	})

	Report(t)
}
