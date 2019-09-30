package rogger

import (
	"bytes"
	"sync"
	"time"
)

var (
	bufferPool *sync.Pool
)

// Entry is the final or intermediate logging data
type Entry struct {
	logger *Logger
	// all the params set by the user
	data Params
	// time at which log was created
	time *time.Time
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
