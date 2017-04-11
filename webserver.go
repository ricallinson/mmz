package main

import (
	"github.com/goforgery/forgery2"
	"github.com/goforgery/mustache"
	"github.com/goforgery/static"
	"log"
	"strconv"
)

// Starts the web server.
func StartWebServer(port int, zilla *Zilla) {

	app := f.CreateApp()
	app.Use(static.Create())
	app.Engine(".html", mustache.Create())

	app.Get("/", func(req *f.Request, res *f.Response, next func()) {
		res.Render("index.html", zilla.GetLiveData())
	})

	app.Get("/livedata", func(req *f.Request, res *f.Response, next func()) {
		res.Send(zilla.GetLiveData())
	})

	app.Get("/update/:attribute", func(req *f.Request, res *f.Response, next func()) {
		var status bool
		attribute := req.Param("attribute")
		value, _ := strconv.Atoi(req.Query("value"))
		switch attribute {
		case "BatteryAmpLimit":
			status = zilla.SetBatteryAmpLimit(value)
		case "LowBatteryVoltageLimit":
			status = zilla.SetLowBatteryVoltageLimit(value)
		case "LowBatteryVoltageIndicator":
			status = zilla.SetLowBatteryVoltageIndicator(value)
		case "NormalMotorAmpLimit":
			status = zilla.SetNormalMotorAmpLimit(value)
		case "SeriesMotorVoltageLimit":
			status = zilla.SetSeriesMotorVoltageLimit(value)
		case "ReverseMotorAmpLimit":
			status = zilla.SetReverseMotorAmpLimit(value)
		case "ReverseMotorVoltageLimit":
			status = zilla.SetReverseMotorVoltageLimit(value)
		case "ParallelMotorAmpLimit":
			status = zilla.SetParallelMotorAmpLimit(value)
		case "ParallelMotorVoltageLimit":
			status = zilla.SetParallelMotorVoltageLimit(value)
		case "ForwardRpmLimit":
			status = zilla.SetForwardRpmLimit(value)
		case "ReverseRpmLimit":
			status = zilla.SetReverseRpmLimit(value)
		case "MaxRpmLimit":
			status = zilla.SetMaxRpmLimit(value)
		case "RpmSensorMotorOne":
			status = zilla.ToggleRpmSensorMotorOne()
			value = Btoi(zilla.GetSettings().RpmSensorMotorOne)
		case "RpmSensorMotorTwo":
			status = zilla.ToggleRpmSensorMotorTwo()
			value = Btoi(zilla.GetSettings().RpmSensorMotorTwo)
		case "AutoShiftingSeriesToParallel":
			status = zilla.ToggleAutoShiftingSeriesToParallel()
			value = Btoi(zilla.GetSettings().AutoShiftingSeriesToParallel)
		case "StallDetectOn":
			status = zilla.ToggleStallDetectOn()
			value = Btoi(zilla.GetSettings().StallDetectOn)
		case "BatteryLightPolarity":
			status = zilla.ToggleBatteryLightPolarity()
			value = Btoi(zilla.GetSettings().BatteryLightPolarity)
		case "CheckEngineLightPolarity":
			status = zilla.ToggleCheckEngineLightPolarity()
			value = Btoi(zilla.GetSettings().CheckEngineLightPolarity)
		case "ReversingContactors":
			status = zilla.ToggleReversingContactors()
			value = Btoi(zilla.GetSettings().ReversingContactors)
		case "SeriesParallelContactors":
			status = zilla.ToggleSeriesParallelContactors()
			value = Btoi(zilla.GetSettings().SeriesParallelContactors)
		case "ForceParallelInReverse":
			status = zilla.ToggleForceParallelInReverse()
			value = Btoi(zilla.GetSettings().ForceParallelInReverse)
		case "InhibitSeriesParallelShifting":
			status = zilla.ToggleInhibitSeriesParallelShifting()
			value = Btoi(zilla.GetSettings().InhibitSeriesParallelShifting)
		case "TachometerDisplayMotorAmps":
			status = zilla.ToggleTachometerDisplayMotorAmps()
			value = Btoi(zilla.GetSettings().TachometerDisplayMotorAmps)
		case "TachometerSixCylinders":
			status = zilla.ToggleTachometerSixCylinders()
			value = Btoi(zilla.GetSettings().TachometerSixCylinders)
		case "ReversesPlugInInputPolarity":
			status = zilla.ToggleReversesPlugInInputPolarity()
			value = Btoi(zilla.GetSettings().ReversesPlugInInputPolarity)
		case "ActivateHEPI":
			status = zilla.ToggleActivateHEPI()
			value = Btoi(zilla.GetSettings().ActivateHEPI)
		case "IsZ2k":
			status = zilla.ToggleIsZ2k()
			value = Btoi(zilla.GetSettings().IsZ2k)
		}
		res.Send(map[string]interface{}{"status": status, "attribute": attribute, "value": value})
	})

	app.Get("/settings", func(req *f.Request, res *f.Response, next func()) {
		res.Render("settings.html", zilla.GetSettings())
	})

	log.Printf("The Manzanita Micro Zilla interface is now running on port '%d'.\n", port)

	app.Listen(port)
}
