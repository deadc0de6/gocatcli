/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package log

import (
	"fmt"
	"os"

	"github.com/TwiN/go-color"
)

var (
	eol       = "\n"
	clearLine = "\r"
	infoPre   = "[INFO] "
	errorPre  = "[ERROR] "
	warnPre   = "[WARN] "
	debugPre  = "[DEBUG] "
	// DebugMode sets debug mode flag
	DebugMode = false
)

// Info print info to stdout
func Info(text string) {
	out := color.InBlue(infoPre) + text + eol
	os.Stdout.WriteString(clearLine + out)
}

// Infof print info to stdout
func Infof(format string, a ...interface{}) {
	out := color.InBlue(fmt.Sprintf(format, a...) + eol)
	os.Stdout.WriteString(clearLine + out)
}

// Error print error to stderr
func Error(err error) {
	out := color.InRed(errorPre) + err.Error() + eol
	os.Stderr.WriteString(clearLine + out)
}

// Errorf print error to stderr
func Errorf(format string, a ...interface{}) {
	out := color.InRed(errorPre) + fmt.Sprintf(format, a...) + eol
	os.Stderr.WriteString(clearLine + out)
}

// Warn print warning to stderr
func Warn(text string) {
	out := color.InRed(warnPre) + text + eol
	os.Stderr.WriteString(clearLine + out)
}

// Warnf print warning to stderr
func Warnf(format string, a ...interface{}) {
	out := color.InRed(warnPre) + fmt.Sprintf(format, a...) + eol
	os.Stderr.WriteString(clearLine + out)
}

// Debug print debug to stderr
func Debug(text string) {
	if !DebugMode {
		return
	}
	out := color.InYellow(debugPre) + text + eol
	os.Stderr.WriteString(clearLine + out)
}

// Debugf print debug to stderr
func Debugf(format string, a ...interface{}) {
	if !DebugMode {
		return
	}
	out := color.InYellow(debugPre) + fmt.Sprintf(format, a...) + eol
	os.Stderr.WriteString(clearLine + out)
}

// Fatal error and exit
func Fatal(err error) {
	Error(err)
	os.Exit(1)
}

// Fatalf error and exit
func Fatalf(format string, a ...interface{}) {
	Errorf(format, a...)
	os.Exit(1)
}
