package main

import (
	"fmt"
	"github.com/cosiner/flag"
	"os"
)

func main() {
	var c Client
	err := flag.NewFlagSet(flag.Flag{}).ParseStruct(&c, os.Args...)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	c.run()
}
