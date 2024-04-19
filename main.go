package main

import (
	"flag"
	"fmt"
	"os"

	valid "github.com/asaskevich/govalidator"
)

func main() {
	if err := CleanTmpDir(); err != nil {
		ErrLog(err.Error())
	}
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [args] [link]\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	url := flag.Arg(0)
	if ok := valid.IsURL(url); !ok {
		ErrLog("Invalid URL")
	}
	if err := Run(url); err != nil {
		ErrLog(err.Error())
	}
}
