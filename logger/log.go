package logger

import (
	"fmt"
	"github.com/go-per/simpkg/types"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-per/simpkg/helpers"
	"github.com/go-per/simpkg/i18n"
	"github.com/go-per/simpkg/parse"
	"github.com/go-per/simpkg/std"
)

// Instance is Logger instance
var Instance ILogger

type MessageType string

const (
	MessageDebug   MessageType = "debug"
	MessageInfo    MessageType = "info"
	MessageSuccess MessageType = "success"
	MessageError   MessageType = "error"
)

// ILogger interface
type ILogger interface {
	SetSplitByDays(bool) ILogger
	SetFilename(string) ILogger
	SetPageSize(int) ILogger
	SetRootPath(string) ILogger
	Log(LogItem, ...any)
	Debug(string, any, ...any)
	Info(string, any, ...any)
	Success(string, any, ...any)
	Error(string, any, ...any)
	Get(int, int, ...string) logsResponse
	LogFiles() []string
}

// Logger struct
type Logger struct {
	fileName    string
	splitByDays bool
	pageSize    int
	rootPath    string
}

// LogItem struct
type LogItem struct {
	Time        types.JSONTime `json:"time"`
	Message     string         `json:"message"`
	Data        any            `json:"data,omitempty"`
	MessageType MessageType    `json:"type"`
}

// logsResponse struct
type logsResponse struct {
	Offset int       `json:"offset"`
	Limit  int       `json:"limit"`
	Size   int       `json:"size"`
	Count  int       `json:"count"`
	Logs   []LogItem `json:"logs"`
	Error  string    `json:"error,omitempty"`
}

// init func
func init() {
	Instance = New()
}

// New create new logger
func New() *Logger {
	return &Logger{
		fileName:    "log.txt",
		pageSize:    20,
		rootPath:    "logs",
		splitByDays: true,
	}
}

// SetSplitByDays set split into days
func (logger *Logger) SetSplitByDays(split bool) ILogger {
	logger.splitByDays = split
	return logger
}

// SetFilename set log file name
func (logger *Logger) SetFilename(fileName string) ILogger {
	logger.fileName = fileName
	return logger
}

// SetPageSize set log page size
func (logger *Logger) SetPageSize(pageSize int) ILogger {
	logger.pageSize = pageSize
	return logger
}

// SetRootPath set log root path
func (logger *Logger) SetRootPath(rootPath string) ILogger {
	logger.rootPath = rootPath
	return logger
}

// Debug log
func (logger *Logger) Debug(message string, data any, args ...any) {
	logger.log(message, MessageDebug, data, args...)
}

// Info log
func (logger *Logger) Info(message string, data any, args ...any) {
	logger.log(message, MessageInfo, data, args...)
}

// Success log
func (logger *Logger) Success(message string, data any, args ...any) {
	logger.log(message, MessageSuccess, data, args...)
}

// Error log
func (logger *Logger) Error(message string, data any, args ...any) {
	logger.log(message, MessageError, data, args...)
}

// Log logs a message
func (logger *Logger) Log(l LogItem, args ...any) {
	logger.log(l.Message, l.MessageType, l.Data, args...)
}

// Get returns logs
func (logger *Logger) Get(offset int, size int, fileName ...string) (response logsResponse) {
	logFile, logPath := logger.getLogFile(fileName...)
	content, err := helpers.ReadFile(logPath)

	response = logsResponse{
		Offset: offset,
		Limit:  logger.pageSize,
		Size:   0,
		Count:  0,
		Logs:   []LogItem{},
	}

	if err != nil {
		std.Error("Could not read log file: %v", err.Error())
		response.Error = i18n.Translate("file_not_exists", logFile)
		return
	}

	var logs []LogItem
	strContent, _ := strings.CutSuffix(string(content), ",")
	jsonContent := []byte("[" + strContent + "]")
	err = parse.Decode(jsonContent, &logs)
	if err != nil {
		std.Error("Could not decode log file: %v", err.Error())
		response.Error = i18n.Translate("file_not_decoded", logFile)
		return
	}

	response.Size = len(logs)
	if response.Size == 0 || offset >= response.Size {
		response.Error = i18n.Translate("out_of_range", offset)
		return
	}

	// set limit
	limit := logger.pageSize - 1
	if size > 0 {
		limit = offset + size
	}

	if limit > response.Size {
		limit = response.Size
	}

	response.Count = limit - offset
	response.Logs = logs[offset:limit]
	return
}

// LogFiles returns log files
func (logger *Logger) LogFiles() []string {
	extension := filepath.Ext(logger.fileName)
	files, err := filepath.Glob(filepath.Join(logger.rootPath, "*"+extension))
	if err != nil {
		std.Error("Could not list log files: %v", err.Error())
		return []string{}
	}

	var logFiles []string
	for _, file := range files {
		logFiles = append(logFiles, strings.Replace(filepath.Base(file), extension, "", -1))
	}

	return logFiles
}

// log message
func (logger *Logger) log(message string, messageTye MessageType, data any, args ...any) {
	msg := message
	if len(args) > 0 {
		msg = fmt.Sprintf(message, args...)
	}

	logItem := LogItem{
		Time:        types.JSONTime(time.Now()),
		Message:     msg,
		MessageType: messageTye,
		Data:        data,
	}

	logger.dump(logItem)
}

// dump log to file
func (logger *Logger) dump(item LogItem) {
	_, p := logger.getLogFile()
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()

	if err != nil {
		std.Error("Could not open log file: %v", err.Error())
		return
	}

	logItem, err := parse.Encode(item)
	if err != nil {
		std.Error("Could not encode log item: %v", err.Error())
		return
	}

	_, err = f.WriteString(string(logItem) + ",\n")
}

// getLogFile returns log file path
func (logger *Logger) getLogFile(name ...string) (n string, p string) {
	fileName := logger.fileName
	extension := filepath.Ext(logger.fileName)
	inputFileName := false
	if name != nil && len(name) > 0 && name[0] != "" {
		fileName = name[0]
		inputFileName = true
	}

	fileName = strings.Replace(fileName, extension, "", 1)
	if !inputFileName && logger.splitByDays {
		fileName = fmt.Sprintf("%s-%s%s", fileName, time.Now().Format("2006-01-02"), extension)
	}
	if inputFileName {
		fileName = fmt.Sprintf("%s%s", fileName, extension)
	}

	_ = helpers.EnsureDir(logger.rootPath)
	return fileName, filepath.Join(logger.rootPath, fileName)
}
