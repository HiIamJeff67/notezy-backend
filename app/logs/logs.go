package logs

import (
	"fmt"
	"log"

	"github.com/comail/colog"
)

func init() {
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ltime, // enable to record the time
	})
	colog.Register() // register colog to handle the print output and take over the stdout io writer
	colog.SetMinLevel(colog.LInfo)
	colog.SetDefaultLevel(colog.LInfo)
	colog.ParseFields(false)
}

// Trace is used to print trace level message
func Trace(fileLine string, v ...interface{}) {
	log.Println("t:", fileLine, fmt.Sprintln(v...))
}

// Debug is used to print debug level message
func Debug(fileLine string, v ...interface{}) {
	log.Println("d:", fileLine, fmt.Sprintln(v...))
}

// Info is used to print info level message
func Info(fileLine string, v ...interface{}) {
	log.Println("i:", fileLine, fmt.Sprintln(v...))
}

// Warn is used to print warning level message
func Warn(fileLine string, v ...interface{}) {
	log.Println("w:", fileLine, fmt.Sprintln(v...))
}

// Error is used to print error level message
func Error(fileLine string, v ...interface{}) {
	log.Println("e:", fileLine, fmt.Sprintln(v...))
}

// Alert is used to print alert level message
func Alert(fileLine string, v ...interface{}) {
	log.Println("alert:", fileLine, fmt.Sprintln(v...))
}

// FTrace is used to print trace level message with format
func FTrace(fileLine string, format string, v ...interface{}) {
	log.Println("t:", fileLine, fmt.Sprintf(format, v...))
}

// FDebug is used to print debug level message with format
func FDebug(fileLine string, format string, v ...interface{}) {
	log.Println("d:", fileLine, fmt.Sprintf(format, v...))
}

// FInfo is used to print info level message with format
func FInfo(fileLine string, format string, v ...interface{}) {
	log.Println("i:", fileLine, fmt.Sprintf(format, v...))
}

// FWarn is used to print warning level message with format
func FWarn(fileLine string, format string, v ...interface{}) {
	log.Println("w:", fileLine, fmt.Sprintf(format, v...))
}

// FError is used to print error level message with format
func FError(fileLine string, format string, v ...interface{}) {
	log.Println("e:", fileLine, fmt.Sprintf(format, v...))
}

// FAlert is used to print alert level message with format
func FAlert(fileLine string, format string, v ...interface{}) {
	log.Println("alert:", fileLine, fmt.Sprintf(format, v...))
}
