package rogger

import (
	"bytes"
	"runtime"
	"sync"
	"time"
)

var (
	bufferPool *sync.Pool
)

// Entry is the final or intermediate logging data
type Entry struct {
	Logger *Logger

	// all the params set by the user
	Data Params

	// time at which log was created
	Time *time.Time

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