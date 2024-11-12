package main

import "fmt"

var debugMode = false

func logDebug(logType string, format string, args ...interface{}) {
	if debugMode {
		fmt.Printf("[DEBUG] ["+logType+"] "+format, args...)
	}
}
