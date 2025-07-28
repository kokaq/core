package profiler

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func writeHeap(path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Printf("failed to create mem profile: %v", err)
		return
	}
	defer f.Close()

	runtime.GC() // force GC
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Printf("failed to write mem profile: %v", err)
	}
}
