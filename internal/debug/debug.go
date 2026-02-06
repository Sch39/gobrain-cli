package debug

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
)

var enabled atomic.Bool

func Enabled() bool {
	return enabled.Load()
}

func Set(v bool) {
	enabled.Store(v)
	if v {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(io.Discard)
		log.SetFlags(0)
	}
}

func Printf(format string, args ...interface{}) {
	if enabled.Load() {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func Println(args ...interface{}) {
	if enabled.Load() {
		fmt.Fprintln(os.Stderr, args...)
	}
}
