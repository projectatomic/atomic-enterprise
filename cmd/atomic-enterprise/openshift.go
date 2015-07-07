package main

import (
	"os"
	"path/filepath"
	"runtime"

	ae "github.com/projectatomic/atomic-enterprise/pkg/cmd/atomic-enterprise"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/util/serviceability"
)

func main() {
	defer serviceability.BehaviorOnPanic(os.Getenv("OPENSHIFT_ON_PANIC"))()
	defer serviceability.Profile(os.Getenv("OPENSHIFT_PROFILE")).Stop()

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	basename := filepath.Base(os.Args[0])
	command := ae.CommandFor(basename)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
