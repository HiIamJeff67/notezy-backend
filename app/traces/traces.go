package traces

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Trace struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
	Package  string `json:"package"`
}

func (t Trace) String() string {
	return fmt.Sprintf("|%s|%s|%s|(%d)", t.Package, t.File, t.Function, t.Line)
}

func (t Trace) FileLineString() string {
	return fmt.Sprintf("|%s|(%d)", t.File, t.Line)
}

func GetTrace(skip int) Trace { // 0th caller
	skip += 1

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return Trace{} // return zero value of Trace when it is not available to get previous trace information
	}

	functionName := "unknown"
	packageName := "unknown"

	if fn := runtime.FuncForPC(pc); fn != nil {
		fullFunctionName := fn.Name()
		parts := strings.Split(fullFunctionName, ".")
		if len(parts) >= 2 {
			packageName = strings.Join(parts[:len(parts)-1], ".")
			functionName = parts[len(parts)-1]
		} else {
			functionName = fullFunctionName
		}
	}

	return Trace{
		File:     filepath.Base(file),
		Line:     line,
		Function: functionName,
		Package:  packageName,
	}
}

func GetTraces(skip int, maxTraceDepth int) []Trace { // 0th caller
	skip += 1
	var traces []Trace

	for i := skip; i < skip+maxTraceDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break // the end of the trace stack
		}

		functionName := "unknown"
		packageName := "unknown"

		if fn := runtime.FuncForPC(pc); fn != nil {
			fullFunctionName := fn.Name()
			parts := strings.Split(fullFunctionName, ".")
			if len(parts) >= 2 {
				packageName = strings.Join(parts[:len(parts)-1], ".")
				functionName = parts[len(parts)-1]
			} else {
				functionName = fullFunctionName
			}
		}

		traces = append(traces, Trace{
			File:     filepath.Base(file),
			Line:     line,
			Function: functionName,
			Package:  packageName,
		})
	}

	return traces
}
