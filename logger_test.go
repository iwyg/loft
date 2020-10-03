package loft

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"
)

type levelRecorder struct {
	handles LogLevel
	lvl     []LogLevel
	mu      sync.Mutex
}

func (r *levelRecorder) Records() []LogLevel {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lvl
}

func (r *levelRecorder) Handle(lvl LogLevel, name string, v ...interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lvl = append(r.lvl, lvl)
}

func (r *levelRecorder) Handlef(lvl LogLevel, name string, format string, v ...interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lvl = append(r.lvl, lvl)
}

func (r *levelRecorder) Handles(lvl LogLevel) bool {
	return lvl >= r.handles
}

type testWriter struct {
	t   testing.TB
	c   string
	buf []byte
	mu  sync.Mutex
}

func (w *testWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return string(w.buf)
}

func (w *testWriter) Write(n []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf = n
	return len(n), nil
}

var logstr = "%s testing.%s: %s\n"

func fmLogStr(t time.Time, cc, msg string) string {
	return fmt.Sprintf(
		logstr,
		t.Format("2006/01/02 15:04:05"),
		cc,
		msg,
	)
}

func TestLoggerLogCascadingHandlers(t *testing.T) {
	dbgW := &testWriter{t: t, c: "dbg"}
	stdW := &testWriter{t: t, c: "std"}
	errW := &testWriter{t: t, c: "err"}
	logger := New("testing", []Handler{
		NewStdHandler(Debug, dbgW, log.LstdFlags),
		NewStdHandler(Info, stdW, log.LstdFlags),
		NewStdHandler(Error, errW, log.LstdFlags),
	})

	d := time.Now()

	logger.Info("ok!")
	logger.Debug("debug!")
	logger.Error("alert!")

	// this should produce an output like:
	// 2019/11/14 12:31:13 testing.INFO: ok!
	// 2019/11/14 12:31:13 testing.DEBUG: debug!
	// 2019/11/14 12:31:13 testing.ERROR: alert!

	wDbg := fmLogStr(d, "DEBUG", "debug!")
	wStd := fmLogStr(d, "INFO", "ok!")
	wErr := fmLogStr(d, "ERROR", "alert!")

	if dbgW.String() != wDbg {
		t.Errorf("want %q, got %q", wDbg, dbgW.String())
	}
	if stdW.String() != wStd {
		t.Errorf("want %q, got %q", wStd, stdW.String())
	}
	if errW.String() != wErr {
		t.Errorf("want %q, got %q", wStd, errW.String())
	}
}

func TestLoggerDebug(t *testing.T) {

	for _, lvl := range []LogLevel{Debug, Info, Notice, Warn, Error, Emergency, Fatal} {
		h := &levelRecorder{handles: Debug}
		logger := New("testing", []Handler{h})
		exp := []LogLevel{lvl, lvl}
		switch lvl {
		case Debug:
			logger.Debug("f")
			logger.Debugf("f")
		case Info:
			logger.Info("f")
			logger.Infof("f")
		case Notice:
			logger.Notice("f")
			logger.Noticef("f")
		case Warn:
			logger.Warn("f")
			logger.Warnf("f")
		case Error:
			logger.Error("f")
			logger.Errorf("f")
		case Emergency:
			logger.Emergency("f")
			logger.Emergencyf("f")
		case Fatal:
			logger.Fatal("f")
			logger.Fatalf("f")
		}

		rec := h.Records()
		if !reflect.DeepEqual(rec, exp) {
			t.Errorf("want %#v, got %#v", exp, rec)
		}
	}
}

