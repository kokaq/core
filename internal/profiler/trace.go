package profiler

import (
	"log"
	"os"
	"runtime/trace"
)

func startTrace(path string) func() {
	f, err := os.Create(path)
	if err != nil {
		log.Printf("failed to create trace file: %v", err)
		return func() {}
	}

	if err := trace.Start(f); err != nil {
		log.Printf("failed to start trace: %v", err)
		return func() {}
	}

	return func() {
		trace.Stop()
		_ = f.Close()
	}
}
