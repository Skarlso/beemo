package internal

import "log"

// Debug sets a flag if more verbose output is needed.
var Debug bool

// LogDebug logs a message and its arguments if debug option is enabled.
func LogDebug(v ...interface{}) {
	if Debug {
		log.Println(v...)
	}
}
