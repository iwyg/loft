package loft

import (
"sync"
)

// LogLevel describes a logging level (e.g. debug, info, error, etc.)
type LogLevel int

const (
	// Debug is the log level used for debugging purpose
	Debug LogLevel = iota - 1
	// Info is the log level used general info purpose
	Info
	// Notice is the log level used for important, none critical messages
	Notice
	// Warn is the log level used for important, less critical messages
	Warn
	// Error is the log level used for logging general errors
	Error
	// Fatal is the log level used for logging fatal errors
	Fatal
	// Emergency is the log level used for logging important emergency situations
	Emergency
)

// Logger is the logging interface that wraps underlying logger services
type Logger interface {
	// Log logs arguments with a given LogLevel
	Log(LogLevel, ...interface{})
	// Logf logs formatted arguments with a given LogLevel and format string
	Logf(LogLevel, string, ...interface{})
	// Info logs arguments with a LogLevel of Info
	Info(...interface{})
	// Infof logs formatted arguments arguments with a LogLevel of Info
	Infof(string, ...interface{})
	// Notice logs arguments with a LogLevel of Notice
	Notice(...interface{})
	// Noticef logs formatted arguments arguments with a LogLevel of Notice
	Noticef(string, ...interface{})
	// Warn logs arguments with a LogLevel of Warn
	Warn(...interface{})
	// Warnf logs formatted arguments arguments with a LogLevel of Warn
	Warnf(string, ...interface{})
	// Debug logs arguments with a LogLevel of Debug
	Debug(...interface{})
	// Debugf logs formatted arguments arguments with a LogLevel of Debug
	Debugf(string, ...interface{})
	// Error logs arguments with a LogLevel of Error
	Error(...interface{})
	// Errorf logs formatted arguments arguments with a LogLevel of Error
	Errorf(string, ...interface{})
	// Emergency logs arguments with a LogLevel of Emergency
	Emergency(...interface{})
	// Emergencyf logs formatted arguments arguments with a LogLevel of Emergency
	Emergencyf(string, ...interface{})
	// Fatal logs arguments with a LogLevel of Fatal
	Fatal(...interface{})
	// Fatalf logs formatted arguments arguments with a LogLevel of Fatal
	Fatalf(string, ...interface{})
}

// DefaultLogger is the package implementation of Logger
type DefaultLogger struct {
	name     string
	handlers Handlers
	handles  map[LogLevel]Handler
	mu       sync.RWMutex
}

// PushHandler pushes a `logging.Handler` on the handlers stack
func (l *DefaultLogger) PushHandler(h Handler) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.handles = make(map[LogLevel]Handler) // reset the resolver map
	l.handlers = append(l.handlers, h)
}

// PopHandler pops a `logging.Handler` from the handlers stack
func (l *DefaultLogger) PopHandler() Handler {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.handles = make(map[LogLevel]Handler) // reset the resolver map
	handler := l.handlers[len(l.handlers)-1:][0]
	l.handlers = l.handlers[:len(l.handlers)-1]

	return handler
}

func (l *DefaultLogger) handler(lvl LogLevel) Handler {
	l.mu.Lock()
	defer l.mu.Unlock()
	if handler, ok := l.handles[lvl]; ok {
		return handler
	}

	for i := range l.handlers {
		handler := l.handlers[len(l.handlers)-1-i]
		if handler.Handles(lvl) {
			l.handles[lvl] = handler
			return handler
		}
	}

	return nil
}

// Log logs arguments with a given LogLevel
func (l *DefaultLogger) Log(lvl LogLevel, args ...interface{}) {
	if logger := l.handler(lvl); logger != nil {
		logger.Handle(lvl, l.name, args...)
	}
}

// Logf logs formatted arguments with a given LogLevel and format string
func (l *DefaultLogger) Logf(lvl LogLevel, format string, args ...interface{}) {
	if logger := l.handler(lvl); logger != nil {
		logger.Handlef(lvl, l.name, format, args...)
	}
}

// Info logs arguments with a LogLevel of Info
func (l *DefaultLogger) Info(args ...interface{}) {
	l.Log(Info, args...)
}

// Infof logs formated arguments arguments with a LogLevel of Info
func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.Logf(Info, format, args...)
}

// Notice logs arguments with a LogLevel of Notice
func (l *DefaultLogger) Notice(args ...interface{}) {
	l.Log(Notice, args...)
}

// Noticef logs formatted arguments arguments with a LogLevel of Notice
func (l *DefaultLogger) Noticef(format string, args ...interface{}) {
	l.Logf(Notice, format, args...)
}

// Debug logs arguments with a LogLevel of Debug
func (l *DefaultLogger) Debug(args ...interface{}) {
	l.Log(Debug, args...)
}

// Debugf logs formatted arguments arguments with a LogLevel of Debug
func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	l.Logf(Debug, format, args...)
}

// Warn logs arguments with a LogLevel of Warn
func (l *DefaultLogger) Warn(args ...interface{}) {
	l.Log(Warn, args...)
}

// Warnf logs formatted arguments arguments with a LogLevel of Warn
func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	l.Logf(Warn, format, args...)
}

// Error logs arguments with a LogLevel of Error
func (l *DefaultLogger) Error(args ...interface{}) {
	l.Log(Error, args...)
}

// Errorf logs formatted arguments arguments with a LogLevel of Error
func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.Logf(Error, format, args...)
}

// Emergency logs arguments with a LogLevel of Emergency
func (l *DefaultLogger) Emergency(args ...interface{}) {
	l.Log(Emergency, args...)
}

// Emergencyf logs formatted arguments arguments with a LogLevel of Emergency
func (l *DefaultLogger) Emergencyf(format string, args ...interface{}) {
	l.Logf(Emergency, format, args...)
}

// Fatal logs arguments with a LogLevel of Fatal
func (l *DefaultLogger) Fatal(args ...interface{}) {
	l.Log(Fatal, args...)
}

// Fatalf logs formatted arguments arguments with a LogLevel of Fatal
func (l *DefaultLogger) Fatalf(format string, args ...interface{}) {
	l.Logf(Fatal, format, args...)
}

// New creates a new *DefaultLogger instance
func New(name string, handlers Handlers) *DefaultLogger {
	if handlers == nil {
		handlers = []Handler{}
	}
	if name == "" {
		name = "default"
	}
	return &DefaultLogger{
		name:     name,
		handlers: handlers,
		handles:  make(map[LogLevel]Handler),
	}
}

