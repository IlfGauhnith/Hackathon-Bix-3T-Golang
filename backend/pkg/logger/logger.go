package logger

import (
	"fmt"
	"io"
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

	// Enable caller reporting
	Log.SetReportCaller(true)

	if err := os.MkdirAll("logs", 0755); err != nil {
		Log.Fatalf("Failed to create logs directory: %v", err)
	}

	logFile, err := os.OpenFile("logs/backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Fatalf("Failed to open log file: %v", err)
	}

	// Set output to both stdout and the file.
	Log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// Set log format to JSON with caller information.
	Log.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})

	Log.SetLevel(logrus.DebugLevel)
}
