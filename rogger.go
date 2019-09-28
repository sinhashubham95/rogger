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

// Logger is the type used for
type Logger struct {
	// it is locked with mutex before any log is sent to this
	// default is os.Stderr.
	// better to set it to a file, which will be rotated automatically.
	out io.Writer

	// Flag for whether to log caller info (off by default)
	reportCaller bool

	// The logging level the logger should log at. defaults to info.
	level Level

	// Used to sync writing to the log. Locking is enabled by Default
	mu MutexWrap
}

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}
