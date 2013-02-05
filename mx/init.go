package mx

// File: initialization of general command line flags.

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	Flag_version     = flag.Bool("v", true, "print version")
	Flag_debug       = flag.Bool("g", false, "Generate debug info")
	Flag_silent      = flag.Bool("s", false, "Don't generate any log info")
	Flag_od          = flag.String("o", "", "set output directory")
	Flag_force       = flag.Bool("f", false, "force start, clean existing output directory")
	Flag_cpuprof     = flag.Bool("cpuprof", false, "Record gopprof CPU profile")
	Flag_memprof     = flag.Bool("memprof", false, "Recored gopprof memory profile")
	Flag_maxprocs    = flag.Int("threads", 0, "maximum number of CPU threads, 0=auto")
	Flag_maxblocklen = flag.Int("maxblocklen", 1<<30, "Maximum size of concurrent blocks")
	Flag_minblocks   = flag.Int("minblocks", 1, "Minimum number of concurrent blocks")
	Flag_gpu         = flag.Int("gpu", 0, "specify GPU")
	Flag_sched       = flag.String("sched", "yield", "CUDA scheduling: auto|spin|yield|sync")
	Flag_pagelock    = flag.Bool("pagelock", true, "enable CUDA memeory page-locking")
)

func Init() {
	flag.Parse()
	if *Flag_version {
		fmt.Print("Mumax Cubed 0.0 alpha ", runtime.GOOS, "_", runtime.GOARCH, " ", runtime.Version(), "(", runtime.Compiler, ")", "\n")
	}

	initOD()
	initLog()
	initTiming()
	initGOMAXPROCS()
	initCpuProf()
	initMemProf()
}

var starttime time.Time

func initTiming() {
	starttime = time.Now()
	AtExit(func() {
		Log("run time:", time.Since(starttime))
	})
}

func initGOMAXPROCS() {
	if *Flag_maxprocs == 0 {
		*Flag_maxprocs = runtime.NumCPU()
		Log("num CPU:", *Flag_maxprocs)
	}
	procs := runtime.GOMAXPROCS(*Flag_maxprocs) // sets it
	Log("GOMAXPROCS:", procs)
}

func initCpuProf() {
	if *Flag_cpuprof {
		// start CPU profile to file
		fname := OD + "/cpu.pprof"
		f, err := os.Create(fname)
		FatalErr(err, "start CPU profile")
		err = pprof.StartCPUProfile(f)
		FatalErr(err, "start CPU profile")
		Log("writing CPU profile to", fname)

		// at exit: exec go tool pprof to generate SVG output
		AtExit(func() {
			pprof.StopCPUProfile()
			me := procselfexe()
			outfile := fname + ".svg"
			SaveCmdOutput(outfile, "go", "tool", "pprof", "-svg", me, fname)
		})
	}
}

func initMemProf() {
	if *Flag_memprof {
		AtExit(func() {
			fname := OD + "/mem.pprof"
			f, err := os.Create(fname)
			defer f.Close()
			FatalErr(err, "start memory profile")
			Log("writing memory profile to", fname)
			FatalErr(pprof.WriteHeapProfile(f), "start memory profile")
			me := procselfexe()
			outfile := fname + ".svg"
			SaveCmdOutput(outfile, "go", "tool", "pprof", "-svg", "--inuse_objects", me, fname)
		})
	}
}

// path to the executable.
func procselfexe() string {
	me, err := os.Readlink("/proc/self/exe")
	PanicErr(err)
	return me
}
