package profiler

type Config struct {
	CPUProfilePath    string
	MemProfilePath    string
	BlockProfilePath  string
	TracePath         string
	GoroutineDumpPath string
}

func Start(cfg Config) func() {
	stoppers := []func(){}

	if cfg.CPUProfilePath != "" {
		stop := startCPU(cfg.CPUProfilePath)
		stoppers = append(stoppers, stop)
	}

	if cfg.TracePath != "" {
		stop := startTrace(cfg.TracePath)
		stoppers = append(stoppers, stop)
	}

	return func() {
		for _, stop := range stoppers {
			stop()
		}

		if cfg.MemProfilePath != "" {
			writeHeap(cfg.MemProfilePath)
		}

		if cfg.BlockProfilePath != "" {
			writeBlock(cfg.BlockProfilePath)
		}

		if cfg.GoroutineDumpPath != "" {
			writeGoroutines(cfg.GoroutineDumpPath)
		}
	}
}
