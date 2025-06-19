package profiler

import (
	"log"
	"os"
	"runtime/pprof"
)

func writeGoroutines(path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Printf("failed to create goroutine dump: %v", err)
		return
	}
	defer f.Close()

	if err := pprof.Lookup("goroutine").WriteTo(f, 2); err != nil {
		log.Printf("failed to dump goroutines: %v", err)
	}
}
