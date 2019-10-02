package rogger

type Formatter interface {
	Format(*Entry) ([]byte, error)
}

// this is to avoid missing fields such as time, msg, etc. which are
// added by default
func fixParamsClash(data Params, reportCaller bool) {
	// time key check
	if t, ok := data[timeKey]; ok {
		data[paramsPrefix+timeKey] = t
		delete(data, timeKey)
	}

	// msg key check
	if m, ok := data[msgKey]; ok {
		data[paramsPrefix+msgKey] = m
		delete(data, msgKey)
	}

	// level key check
	if l, ok := data[levelKey]; ok {
		data[paramsPrefix+levelKey] = l
		delete(data, levelKey)
	}

	// err key check
	if e, ok := data[errKey]; ok {
		data[paramsPrefix+errKey] = e
		delete(data, errKey)
	}

	// func and file check
	if reportCaller {
		if fu, ok := data[funcKey]; ok {
			data[paramsPrefix+funcKey] = fu
			delete(data, funcKey)
		}
		if fi, ok := data[fileKey]; ok {
			data[paramsPrefix+fileKey] = fi
			delete(data, fileKey)
		}
	}
}
