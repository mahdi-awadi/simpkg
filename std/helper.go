package std

import (
	"fmt"
	"os"
	"time"
)

type OutType int

// keep debug mode
var isDebugMode = true

const (
	OutInfo OutType = iota
	OutDebug
	OutError
	OutSuccess
)

// SetIsDebug set debug mode
func SetIsDebug(v bool) {
	isDebugMode = v
}

func Out(stdType OutType, format string, args ...any) {
	if !isDebugMode && stdType != OutError {
		return
	}

	var logType string
	switch stdType {
	case OutInfo:
		logType = "INFO"
	case OutDebug:
		logType = "DEBU"
	case OutError:
		logType = "ERRO"
	case OutSuccess:
		logType = "SUCC"
	}

	format = "%s - %s | " + format + " \n"
	args = append([]any{
		time.Now().Format("15:04:05.000"),
		logType,
	}, args...)
	fmt.Printf(format, args...)
}

func Debug(format string, args ...any) {
	Out(OutDebug, format, args...)
}

func Info(format string, args ...any) {
	Out(OutInfo, format, args...)
}

func Success(format string, args ...any) {
	Out(OutSuccess, format, args...)
}

func Error(format string, args ...any) {
	Out(OutError, format, args...)
}

// PressAnyKeyToExit wait for key press to exit
func PressAnyKeyToExit() {
	Info("Press `Enter` key to exit...")
	_, _ = fmt.Scanln()
	os.Exit(1)
}
