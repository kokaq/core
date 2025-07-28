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
	pc, file, line, ok := runtime.Caller(1)

	var (
		methodName = "unknown"
		location   = "unknown"
	)

	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			// Get just the method name (e.g., strip package path)
			fullName := fn.Name()
			methodName = filepath.Ext(fullName)[1:] // remove leading dot
		}
		location = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	// Format timestamp and message
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)

	// Color based on level
	levelUpper := strings.ToUpper(level)
	color := levelColors[levelUpper]
	if color == "" {
		color = "\033[37m" // Default: white
	}
	// Final log output
	fmt.Printf("[%s] [%s] [%s] [%s] [%s] %s\n",
		timestamp[:10], // [DATE]
		timestamp[11:], // [TIME]
		methodName,     // [METHOD]
		levelUpper,     // [LEVEL]
		location,       // [FILE:LINE]
		message,        // MESSAGE
	)
}
