package rogger

import (
	"io"
	"os"
	"sync"
	"time"
)

// Params type is used to pass to WithParams
type Params map[string]interface{}

// Level type
type Level uint32

// convert level to a string
func (l Level) String() string {
	switch l {
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

// Logger is the type used for main logging
type Logger struct {
	// it is locked with mutex before any log is sent to this
	// default is os.Stderr.
	// better to set it to a file, which will be rotated automatically.
	Out io.Writer

	// formatter formats logs before finally sending to the writer
	Formatter Formatter

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

func (logger *Logger) IsLevelEnabled(level Level) bool {
	return level >= logger.Level
}

// Creates a new logger with default values. You can also just
// instantiate your own:
//    var log = &Logger {
//      Out: os.Stderr,
//      Formatter: new(TextFormatter),
//      Level: InfoLevel,
//    }
// It's recommended to make this a global instance called `log`.
func New() *Logger {
	return &Logger{
		Out:          os.Stderr,
		Formatter:    new(TextFormatter),
		ReportCaller: false,
		Level:        InfoLevel,
	}
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

func (logger *Logger) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{}
	logger.entryPool.Put(entry)
}

// Adds a param to the log entry, and logs when Debug, Print, Info,
// Warn, Error or Fatal is called.
func (logger *Logger) WithParam(key string, value interface{}) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithParam(key, value)
}

// Adds a list of params to the log entry, and logs when Debug, Print, Info,
// Warn, Error or Fatal is called.
func (logger *Logger) WithParams(params Params) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithParams(params)
}

// Add an error as single field to the Entry, and logs when Debug, Print, Info,
// Warn, Error or Fatal is called.
func (logger *Logger) WithError(err error) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}

// Overrides the time of the log entry.
func (logger *Logger) WithTime(t time.Time) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithTime(t)
}

func (logger *Logger) Log(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		defer logger.releaseEntry(entry)
		entry.Log(level, args...)
	}
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Log(DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Log(InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Log(WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Log(ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.Log(FatalLevel, args...)
	os.Exit(1)
}

func (logger *Logger) Logf(level Level, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		defer logger.releaseEntry(entry)
		entry.Logf(level, format, args...)
	}
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.Logf(DebugLevel, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.Logf(InfoLevel, format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.Logf(WarnLevel, format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.Logf(ErrorLevel, format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.Logf(FatalLevel, format, args...)
	os.Exit(1)
}

func (logger *Logger) Logln(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		defer logger.releaseEntry(entry)
		entry.Logln(level, args...)
	}
}

func (logger *Logger) Debugln(args ...interface{}) {
	logger.Logln(DebugLevel, args...)
}

func (logger *Logger) Infoln(args ...interface{}) {
	logger.Logln(InfoLevel, args...)
}

func (logger *Logger) Warnln(args ...interface{}) {
	logger.Logln(WarnLevel, args...)
}

func (logger *Logger) Errorln(args ...interface{}) {
	logger.Logln(ErrorLevel, args...)
}

func (logger *Logger) Fatalln(args ...interface{}) {
	logger.Logln(FatalLevel, args...)
	logger.Exit(1)
}

// exit function called to exit the application
// having a function makes us able to use the code commonly
func (*Logger) Exit(code int) {
	os.Exit(code)
}

func (logger *Logger) SetNoLock() {
	logger.mu.disable()
}

// SetLevel sets the logger level.
func (logger *Logger) SetLevel(level Level) {
	logger.mu.lock()
	defer logger.mu.unlock()
	logger.Level = level
}

// SetFormatter sets the logger formatter
// by default the formatter is set to text formatting
// you can even create a custom formatter which implements the formatter interface
func (logger *Logger) SetFormatter(formatter Formatter) {
	logger.mu.lock()
	defer logger.mu.unlock()
	logger.Formatter = formatter
}

// SetOutput sets the logger output
func (logger *Logger) SetOutput(output io.Writer) {
	logger.mu.lock()
	defer logger.mu.unlock()
	logger.Out = output
}

func (logger *Logger) SetReportCaller(reportCaller bool) {
	logger.mu.lock()
	defer logger.mu.unlock()
	logger.ReportCaller = reportCaller
}
