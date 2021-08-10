package main

import (
	"github.com/cosiner/flag"
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func main() {
	var c Client
	err := flag.NewFlagSet(flag.Flag{}).ParseStruct(&c, os.Args...)
	if err != nil {
		log.Fatal(err)
	}

	if c.Verbose {
		log.SetLevel(logrus.TraceLevel)
	}

	c.run()
}
