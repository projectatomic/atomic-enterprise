package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/projectatomic/atomic-enterprise/pkg/cmd/admin"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/cli"
	"github.com/projectatomic/atomic-enterprise/pkg/cmd/atomic-enterprise"
)

func OutDir(path string) (string, error) {
	outDir, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	stat, err := os.Stat(outDir)
	if err != nil {
		return "", err
	}

	if !stat.IsDir() {
		return "", fmt.Errorf("output directory %s is not a directory\n", outDir)
	}
	outDir = outDir + "/"
	return outDir, nil
}

func main() {
	// use os.Args instead of "flags" because "flags" will mess up the man pages!
	path := "rel-eng/completions/bash/"
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [output directory]\n", os.Args[0])
		os.Exit(1)
	}

	outDir, err := OutDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}
	outFile_atomic-enterprise := outDir + "atomic-enterprise"
	atomic-enterprise := atomic-enterprise.NewCommandAtomicEnterprise("atomic-enterprise")
	atomic-enterprise.GenBashCompletionFile(outFile_atomic-enterprise)

	outFile_osc := outDir + "oc"
	oc := cli.NewCommandCLI("oc", "atomic-enterprise cli")
	oc.GenBashCompletionFile(outFile_osc)

	outFile_osadm := outDir + "oadm"
	oadm := admin.NewCommandAdmin("oadm", "atomic-enterprise admin", ioutil.Discard)
	oadm.GenBashCompletionFile(outFile_osadm)
}
