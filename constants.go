package rogger

import "time"

// Level data
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

// caller information
const (
	maxCallerDepth int = 25
	knownFrames    int = 4
)

// keys
const (
	timeKey  = "time"
	msgKey   = "message"
	levelKey = "level"
	errKey   = "error"
	funcKey  = "func"
	fileKey  = "file"
)

// params clash prefix
const (
	paramsPrefix = "params"
)

// formatter constants
const (
	defaultTimestampFormat = time.RFC3339
)
