package main

import "github.com/kokaq/core/profiler"

func main() {
	stop := profiler.Start(profiler.Config{
		CPUProfilePath:    "cpu.prof",
		MemProfilePath:    "mem.prof",
		BlockProfilePath:  "block.prof",
		GoroutineDumpPath: "goroutines.prof",
		TracePath:         "trace.out",
	})
	defer stop()
}
