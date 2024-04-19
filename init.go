package main

import (
	"flag"
	"os"
)

var (
	OutputFilaName string
	PassedURL      string
	Verbose        bool
	TmpDir         string
	Workers        int
)

func init() {
	flag.StringVar(&OutputFilaName, "o", "", "output filename eg: ./dir/vid.mp4")
	flag.IntVar(&Workers, "w", 4, "Number of goroutines for concurrent downloading")
	flag.BoolVar(&Verbose, "v", false, "Make gsdm verbose during the operation")
	flag.Parse()
	TmpDir = os.TempDir()
	if Workers <= 0 {
		ErrLog("Number of goroutines must be a <= 1")
	} else if Workers > 20 {
		WarnLog("Setting a high number of goroutines may lead to excessive resource usage and server overload.")
	}
}
