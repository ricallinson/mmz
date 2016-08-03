package main

var Codes map[int]string = map[int]string{
    1111: "Unknown mode, no error.",
    1112: "Hairball watchdog reset.",
    1113: "Hairball EEPROM CRC error.",
    1114: "Controller watchdog reset.",
    1121: "Controller EEPROM CRC error.",
    1122: "Controller Desat error.",
    1123: "Power section failed test.",
    1124: "Main contactor stuck on.",
    1131: "Shorted/Loaded Controller during precharge.",
    1132: "Controller did not communicate during precharge.",
    1133: "Lost communication with controller during use, either direction.",
    1134: "Lost communication to controller, still receiving from controller.",
    1141: "Main contactor high resistance.",
    1142: "Controller still not off. Main contactor trying to turn off.",
    1143: "Motor contactor state machine got an illegal value, software error.",
    1144: "Motor contactors no long match requested state.",
    1211: "Controller still not off. Motor contactors trying to turn off.",
    1212: "Motor contactors did not turn off.",
    1213: "Motor contactors did not turn on.",
    1214: "OPen pot wire.",
    1221: "Major over-speed, either motor over red line.",
    1222: "Unused.",
    1223: "SLI battery below warning threshold.",
    1224: "SLI battery too low and caused shutdown of controller."
    1231: "Propulsion pack open, no contactor drop and controller is not responding.",
    1232: "This should never happen. Contact the factory.",
    1233: "Hall effect pedal input invalid.",
    1234: "Motor voltage is high on startup.",
    1241: "Key input not on while start is asserted.",
    
}
