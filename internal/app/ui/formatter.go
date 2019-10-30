package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// Title prints a title.
func Title(format string, args ...interface{}) {
	color.Blue(format+"\n", args...)
}

// Phase prints a phase.
func Phase(format string, args ...interface{}) {
	fmt.Printf("- "+format+"\n", args...)
}

// Step prints a step.
func Step(format string, args ...interface{}) {
	color.Green("  "+format+"\n", args...)
}

// Trace prints a trace statement.
func Trace(format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}

// Debug prints a debug statement.
func Debug(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Info prints an info statement.
func Info(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Warn prints an warn statement.
func Warn(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Error prints an error statement.
func Error(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, color.RedString(format, args...))
}

// Panic prints an panic statement.
func Panic(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

// Fatal prints an fatal statement.
func Fatal(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
