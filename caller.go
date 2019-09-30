package rogger

import (
	"runtime"
	"strings"
	"sync"
)

var (
	packageName    string
	minCallerDepth int
	callerInitOnce sync.Once
)

const (
	maxCallerDepth int = 25
	knownFrames    int = 4
)

// getPackageName reduces a fully qualified function name to the package name
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, 2)
		packageName = getPackageName(runtime.FuncForPC(pcs[1]).Name())
		minCallerDepth = knownFrames
	})
	// Restrict the look back frames to avoid runaway lookups
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		// If the caller isn't part of this package, we're done
		if pkg != packageName {
			return &f
		}
	}
	// if we got here, we failed to find the caller's context
	return nil
}
