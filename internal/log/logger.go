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
	_, err := os.Stdout.WriteString(clearLine + out)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Infof print info to stdout
func Infof(format string, a ...interface{}) {
	out := color.InBlue(fmt.Sprintf(format, a...) + eol)
	_, err := os.Stdout.WriteString(clearLine + out)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Error print error to stderr
func Error(err error) {
	out := color.InRed(errorPre) + err.Error() + eol
	_, err = os.Stderr.WriteString(clearLine + out)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Errorf print error to stderr
func Errorf(format string, a ...interface{}) {
	out := color.InRed(errorPre) + fmt.Sprintf(format, a...) + eol
	_, err := os.Stderr.WriteString(clearLine + out)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Warn print warning to stderr
func Warn(text string) {
	out := color.InRed(warnPre) + text + eol
	_, err := os.Stderr.WriteString(clearLine + out)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Warnf print warning to stderr
func Warnf(format string, a ...interface{}) {
	out := color.InRed(warnPre) + fmt.Sprintf(format, a...) + eol
	_, err := os.Stderr.WriteString(clearLine + out)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Debug print debug to stderr
func Debug(text string) {
	if !DebugMode {
		return
	}
	out := color.InYellow(debugPre) + text + eol
	_, err := os.Stderr.WriteString(clearLine + out)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Debugf print debug to stderr
func Debugf(format string, a ...interface{}) {
	if !DebugMode {
		return
	}
	out := color.InYellow(debugPre) + fmt.Sprintf(format, a...) + eol
	_, err := os.Stderr.WriteString(clearLine + out)
	if err != nil {
		fmt.Println(err.Error())
	}
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
