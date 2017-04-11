package main

var CodeList []uint64 = []uint64{
	1111,
	1112,
	1113,
	1114,
	1121,
	1122,
	1123,
	1124,
	1131,
	1132,
	1133,
	1134,
	1141,
	1142,
	1143,
	1144,
	1211,
	1212,
	1213,
	1214,
	1221,
	1222,
	1223,
	1224,
	1231,
	1232,
	1233,
	1234,
	1241,
	1311,
	1312,
	1313,
	1314,
	1321,
	1322,
	1323,
	1324,
	1331,
	1332,
	1333,
	1334,
	1411,
	1414,
}

var Codes map[string]string = map[string]string{
	// Error codes.
	"1111": "Unknown mode, no error.",
	"1112": "Hairball watchdog reset.",
	"1113": "Hairball EEPROM CRC error.",
	"1114": "Controller watchdog reset.",
	"1121": "Controller EEPROM CRC error.",
	"1122": "Controller Desat error.",
	"1123": "Power section failed test.",
	"1124": "Main contactor stuck on.",
	"1131": "Shorted/Loaded Controller during precharge.",
	"1132": "Controller did not communicate during precharge.",
	"1133": "Lost communication with controller during use, either direction.",
	"1134": "Lost communication to controller, still receiving from controller.",
	"1141": "Main contactor high resistance.",
	"1142": "Controller still not off. Main contactor trying to turn off.",
	"1143": "Motor contactor state machine got an illegal value, software error.",
	"1144": "Motor contactors no long match requested state.",
	"1211": "Controller still not off. Motor contactors trying to turn off.",
	"1212": "Motor contactors did not turn off.",
	"1213": "Motor contactors did not turn on.",
	"1214": "Open pot wire.",
	"1221": "Major over-speed, either motor over red line.",
	"1222": "Unused.",
	"1223": "SLI battery below warning threshold.",
	"1224": "SLI battery too low and caused shutdown of controller.",
	"1231": "Propulsion pack open, no contactor drop and controller is not responding.",
	"1232": "This should never happen. Contact the factory.",
	"1233": "Hall effect pedal input invalid.",
	"1234": "Motor voltage is high on startup.",
	"1241": "Key input not on while start is asserted.",
	// Operating codes.
	"1311": "Waiting for key.",
	"1312": "Waiting for start signal.",
	"1313": "Waiting for zero pot.",
	"1314": "Waiting for throttle input.",
	"1321": "Waiting for go button.", // Drag race mode only.
	"1322": "Direction selected not allowed (rolling too fast or inactive state).",
	"1323": "Battery voltage limit active.",
	"1324": "Motor current limit active.",
	"1331": "Battery current limit active.",
	"1332": "Temperature current limit active",
	"1333": "SPI packet error in controller.",
	"1334": "Controller waiting for enable signal.",
	"1411": "Normal driving.",
	"1414": "Waiting for the vehicle to be unplugged.",
	// Operating states.
	"G": "Shifting in progress.",
	"O": "Main contactor is on OK.",
	"M": "Motor contactors are on OK.",
	"R": "Direction is reverse.",
	"F": "Direction is forward.",
	"P": "Motors are in parallel.",
	"S": "Stopped state or motors are in series.",
	"V": "Main contactor has voltage drop (>5V) across power terminals.",
}
