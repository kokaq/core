package profiler

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func writeBlock(path string) {
	runtime.SetBlockProfileRate(1)
	f, err := os.Create(path)
	if err != nil {
		log.Printf("failed to create block profile: %v", err)
		return
	}
	defer f.Close()

	if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
		log.Printf("failed to write block profile: %v", err)
	}
}
