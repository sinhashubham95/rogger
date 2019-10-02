package rogger

import (
	"bytes"
	"fmt"
	"sort"
)

type TextFormatter struct {
	// Disable timestamp logging
	DisableTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output.
	DisableSorting bool
}

func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(Params)
	for k, v := range entry.Data {
		data[k] = v
	}
	fixParamsClash(data, entry.HasCaller())
	paramKeys := make([]string, 0, len(data))
	for k := range data {
		paramKeys = append(paramKeys, k)
	}
	var funcVal, fileVal string
	fixedKeys := make([]string, 0, 4+len(data))
	if !f.DisableTimestamp {
		fixedKeys = append(fixedKeys, timeKey)
	}
	if entry.Message != "" {
		fixedKeys = append(fixedKeys, msgKey)
	}
	fixedKeys = append(fixedKeys, levelKey)
	if entry.err != "" {
		fixedKeys = append(fixedKeys, errKey)
	}
	if entry.HasCaller() {
		funcVal = entry.Caller.Function
		if funcVal != "" {
			fixedKeys = append(fixedKeys, funcKey)
		}
		fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if fileVal != "" {
			fixedKeys = append(fixedKeys, fileKey)
		}
	}
	if f.DisableSorting {
		fixedKeys = append(fixedKeys, paramKeys...)
	} else {
		sort.Strings(paramKeys)
		fixedKeys = append(fixedKeys, paramKeys...)
	}
	tsFormat := f.TimestampFormat
	if tsFormat == "" {
		tsFormat = defaultTimestampFormat
	}
	buffer := entry.Buffer
	if buffer == nil {
		buffer = &bytes.Buffer{}
	}
	for _, key := range fixedKeys {
		var value interface{}
		switch key {
		case timeKey:
			value = entry.Time.Format(tsFormat)
		case msgKey:
			value = entry.Message
		case levelKey:
			value = entry.Level.String()
		case errKey:
			value = entry.err
		case funcKey:
			value = funcVal
		case fileKey:
			value = fileVal
		default:
			value = data[key]
		}
		appendData(buffer, key, value)
	}
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}

func appendData(buffer *bytes.Buffer, key string, value interface{}) {
	if buffer.Len() > 0 {
		buffer.WriteByte(' ')
	}
	buffer.WriteString(key)
	buffer.WriteByte('=')
	appendValue(buffer, value)
}

func appendValue(buffer *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}
	if needsQuoting(stringVal) {
		buffer.WriteString(fmt.Sprintf("%q", stringVal))
	} else {
		buffer.WriteString(stringVal)
	}
}

func needsQuoting(text string) bool {
	if len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}
