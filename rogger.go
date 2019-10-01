package rogger

import (
	"io"
	"sync"
)

// Params type is used to pass to WithParams
type Params map[string]interface{}

// Level type
type Level uint32

// convert level to a string
func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	}
	return "unknown"
}

const (
	// TraceLevel level. More informational events than debug.
	TraceLevel Level = iota
	// DebugLevel level. Usually only enabled when debugging.
	DebugLevel
	// InfoLevel level. Operational information about what's going on in the application.
	InfoLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// ErrorLevel level. Used for errors that should definitely be noted.
	ErrorLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`.
	FatalLevel
)

// Logger is the type used for main logging
type Logger struct {
	// it is locked with mutex before any log is sent to this
	// default is os.Stderr.
	// better to set it to a file, which will be rotated automatically.
	Out io.Writer

	// formatter formats logs before finally sending to the writer
	Formatter

	// Flag for whether to log caller info (off by default)
	ReportCaller bool

	// The logging level the logger should log at. defaults to info.
	Level Level

	// Used to sync writing to the log. Locking is enabled by Default
	mu mutexWrap

	// Reusable empty log entries
	entryPool sync.Pool
}

type mutexWrap struct {
	m        sync.Mutex
	disabled bool
}

func (mw *mutexWrap) lock() {
	if !mw.disabled {
		mw.m.Lock()
	}
}

func (mw *mutexWrap) unlock() {
	if !mw.disabled {
		mw.m.Unlock()
	}
}

func (mw *mutexWrap) disable() {
	mw.disabled = true
}
