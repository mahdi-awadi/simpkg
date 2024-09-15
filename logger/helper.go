package logger

// Debug func
func Debug(message string, data any, args ...any) {
	Instance.Debug(message, data, args...)
}

// Info func
func Info(message string, data any, args ...any) {
	Instance.Info(message, data, args...)
}

// Success func
func Success(message string, data any, args ...any) {
	Instance.Success(message, data, args...)
}

// Error func
func Error(message string, data any, args ...any) {
	Instance.Error(message, data, args...)
}
