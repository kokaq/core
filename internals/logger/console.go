package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Map log levels to ANSI color codes
var levelColors = map[string]string{
	"DEBUG": "\033[36m", // Cyan
	"INFO":  "\033[32m", // Green
	"WARN":  "\033[33m", // Yellow
	"ERROR": "\033[31m", // Red
}

// ANSI reset
const colorReset = "\033[0m"

func ConsoleLog(level string, format string, args ...any) {
	// Get caller info
	pc, _, _, ok := runtime.Caller(1)

	var (
		methodName = "unknown"
		//location   = "unknown"
	)

	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			// Get just the method name (e.g., strip package path)
			fullName := fn.Name()
			methodName = filepath.Ext(fullName)[1:] // remove leading dot
		}
		//location = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	// Format timestamp and message
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// Safe message formatting
	var message string
	if strings.Contains(format, "%") {
		// Use Sprintf only if formatting directives are present
		message = fmt.Sprintf(format, args...)
	} else {
		// Otherwise concatenate all args
		message = format
		if len(args) > 0 {
			message += " " + fmt.Sprint(args...)
		}
	}

	// Color based on level
	levelUpper := strings.ToUpper(level)
	color := levelColors[levelUpper]
	if color == "" {
		color = "\033[37m" // Default: white
	}
	// Final log output
	fmt.Printf("[%s] [%s] [%s] [%s] %s\n",
		timestamp[:10], // [DATE]
		timestamp[11:], // [TIME]
		levelUpper,     // [LEVEL]
		methodName,     // [METHOD]
		//location,       // [FILE:LINE]
		message, // MESSAGE
	)
}
