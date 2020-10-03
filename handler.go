package loft
import (
	"io"
	"log"
	"strings"
)

// Handler is a interface describing a loggin handler
// Implementation should wrap underlying loggin services
type Handler interface {
	// Handles determines if this handler can handle a given `logging.LogLevel` `lvl`
	Handles(LogLevel) bool
	// Handle passes arguments to the underlying logging service
	Handle(lvl LogLevel, name string, args ...interface{})
	// Handlef passes a formatted message to the underlying logging service
	Handlef(lvl LogLevel, name, format string, args ...interface{})
}

// Handlers is a list of logging.Handler
type Handlers []Handler

var labelMap = map[LogLevel]string{
	Debug:     "DEBUG",
	Info:      "INFO",
	Notice:    "NOTICE",
	Warn:      "WARN",
	Error:     "ERROR",
	Fatal:     "FATAL",
	Emergency: "EMERGENCY",
}

func prefix(lvl LogLevel, name string) string {
	return strings.Join([]string{name, ".", labelMap[lvl], ": "}, "")
}

var StdLogHandlerPrefixFunc = prefix

// StdLogHandler is a Handler that uses `*log.Logger` as the underlying printer
type StdLogHandler struct {
	handles LogLevel
	logger  *log.Logger
}

// Handles determines if this handler can handle a given `logging.LogLevel` `lvl`
func (l *StdLogHandler) Handles(lvl LogLevel) bool {
	return lvl >= l.handles
}

// Handle passes arguments to the underlying logging service
func (l *StdLogHandler) Handle(lvl LogLevel, name string, args ...interface{}) {
	v := append([]interface{}{prefix(lvl, name)}, args...)
	l.logger.Print(v...)
}

// Handlef passes a formatted message to the underlying logging service
func (l *StdLogHandler) Handlef(lvl LogLevel, name, format string, args ...interface{}) {
	l.logger.Printf(strings.Join([]string{prefix(lvl, name), format}, ""), args...)
}

// NewStdHandler creates a default implementation of `Handler`
func NewStdHandler(lvl LogLevel, w io.Writer, flag int) *StdLogHandler {
	return &StdLogHandler{handles: lvl, logger: log.New(w, "", flag)}
}
