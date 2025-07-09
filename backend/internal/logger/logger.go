package logger

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

// Log is the exported logger instance.
var Log *Logger

func init() {
	base := logrus.New()

	Log = &Logger{base}

	// Set output to both stdout and the file.
	Log.SetOutput(os.Stdout)

	// Enable caller reporting
	Log.SetReportCaller(true)

	// Set log format to JSON with caller information.
	Log.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})
}
