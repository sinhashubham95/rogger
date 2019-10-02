package rogger

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"
)

var (
	bufferPool *sync.Pool
)

// Errors
var (
	LoggerNotAttached = errors.New("logger not associated to the entry created")
)

// Entry is the final or intermediate logging data
type Entry struct {
	Logger *Logger

	// all the params set by the user
	Data Params

	// time at which log was created
	Time time.Time

	// level the log entry was logged at
	Level Level

	// calling method with package name
	Caller *runtime.Frame

	// log message
	Message string

	// When formatter is called in entry.log(), a Buffer may be set to entry
	Buffer *bytes.Buffer

	// err may contain a field formatting error
	err string
}

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	// start at the bottom of the stack before the package name is cached
	minCallerDepth = 1
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		Data:   make(Params),
	}
}

func (entry *Entry) HasCaller() bool {
	return entry.Logger != nil && entry.Logger.ReportCaller && entry.Caller != nil
}

func (entry *Entry) String() (string, error) {
	if entry.Logger == nil {
		return "", LoggerNotAttached
	}
	formatted, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}

// Add an error as single field to the Entry
func (entry *Entry) WithError(err error) *Entry {
	return entry.WithParam(errKey, err)
}

// Add a single param to the Entry.
func (entry *Entry) WithParam(key string, value interface{}) *Entry {
	return entry.WithParams(Params{key: value})
}

// Add a map of params to the Entry
func (entry *Entry) WithParams(params Params) *Entry {
	data := make(Params, len(entry.Data)+len(params))
	for k, v := range entry.Data {
		data[k] = v
	}
	err := entry.err
	for k, v := range params {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch t.Kind() {
			case reflect.Func:
				isErrField = true
			case reflect.Ptr:
				isErrField = t.Elem().Kind() == reflect.Func
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", k)
			if err != "" {
				err = entry.err + ", " + tmp
			} else {
				err = tmp
			}
		} else {
			data[k] = v
		}
	}
	return &Entry{
		Logger:  entry.Logger,
		Data:    data,
		Time:    entry.Time,
		Level:   entry.Level,
		Caller:  entry.Caller,
		Message: entry.Message,
		Buffer:  entry.Buffer,
		err:     err,
	}
}

// Overrides the time of the log entry.
func (entry *Entry) WithTime(t time.Time) *Entry {
	return &Entry{
		Logger:  entry.Logger,
		Data:    entry.Data,
		Time:    t,
		Level:   entry.Level,
		Caller:  entry.Caller,
		Message: entry.Message,
		Buffer:  entry.Buffer,
		err:     entry.err,
	}
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry Entry) log(l Level, msg string) {
	var buffer *bytes.Buffer

	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}

	entry.Level = l
	entry.Message = msg
	if entry.Logger.ReportCaller {
		entry.Caller = getCaller()
	}

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer

	entry.write()

	entry.Buffer = nil
}

func (entry *Entry) write() {
	entry.Logger.mu.lock()
	defer entry.Logger.mu.unlock()
	formattedLog, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
	} else {
		_, err = entry.Logger.Out.Write(formattedLog)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
	}
}

func (entry *Entry) Log(level Level, args ...interface{}) {
	if entry.Logger == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Logger not attached")
		return
	}
	if entry.Logger.IsLevelEnabled(level) {
		entry.log(level, fmt.Sprint(args...))
	}
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.Log(DebugLevel, args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.Log(InfoLevel, args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.Log(WarnLevel, args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.Log(ErrorLevel, args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.Log(FatalLevel, args...)
	if entry.Logger == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Logger not attached")
		return
	}
	entry.Logger.Exit(1)
}

func (entry *Entry) Logf(level Level, format string, args ...interface{}) {
	if entry.Logger == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Logger not attached")
		return
	}
	if entry.Logger.IsLevelEnabled(level) {
		entry.log(level, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Logf(DebugLevel, format, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Logf(InfoLevel, format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Logf(WarnLevel, format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Logf(ErrorLevel, format, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Logf(FatalLevel, format, args...)
	if entry.Logger == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Logger not attached")
		return
	}
	entry.Logger.Exit(1)
}

func (entry *Entry) Logln(level Level, args ...interface{}) {
	if entry.Logger == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Logger not attached")
		return
	}
	if entry.Logger.IsLevelEnabled(level) {
		entry.log(level, fmt.Sprintln(args...))
	}
}

func (entry *Entry) Debugln(args ...interface{}) {
	entry.Logln(DebugLevel, args...)
}

func (entry *Entry) Infoln(args ...interface{}) {
	entry.Logln(InfoLevel, args...)
}

func (entry *Entry) Warnln(args ...interface{}) {
	entry.Logln(WarnLevel, args...)
}

func (entry *Entry) Errorln(args ...interface{}) {
	entry.Logln(ErrorLevel, args...)
}

func (entry *Entry) Fatalln(args ...interface{}) {
	entry.Logln(FatalLevel, args...)
	if entry.Logger == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Logger not attached")
		return
	}
	entry.Logger.Exit(1)
}
