package client

import "fmt"

type RequestLogger struct{}

func (l *RequestLogger) Errorf(format string, v ...any) {
	fmt.Println(fmt.Sprintf("[ERROR] "+format, v...))
}

func (l *RequestLogger) Warnf(format string, v ...any) {
	fmt.Println(fmt.Sprintf("[WARNING] "+format, v...))
}

func (l *RequestLogger) Debugf(format string, v ...any) {
	fmt.Println(fmt.Sprintf("[DEBUG] "+format, v...))
}
