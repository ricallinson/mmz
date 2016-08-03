package main

var Codes map[int]string = map[int]string{
    1111: "Unknown mode, no error",
    1112: "Hairball watchdog reset",
    1113: "Hairball EEPROM CRC error",
    1114: "Controller watchdog reset",
    1121: "Controller EEPROM CRC error",
    1122: "Controller Desat error",
    1123: "Power section failed test",
    1124: "Main Contactor stuck on",
    1131: "Shorted/Loaded Controller during precharge",
    1132: "Controller did not communicate during precharge",
    1133: "Lost communitation with controller during use, either direction",
}
