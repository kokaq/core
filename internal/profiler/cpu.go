package profiler

import (
	"log"
	"os"
	"runtime/pprof"
)

func startCPU(path string) func() {
	f, err := os.Create(path)
	if err != nil {
		log.Printf("failed to create CPU profile: %v", err)
		return func() {}
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Printf("failed to start CPU profiling: %v", err)
		return func() {}
	}

	return func() {
		pprof.StopCPUProfile()
		_ = f.Close()
	}
}
