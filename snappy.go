package main

import (
	"errors"
	"fmt"
	"github.com/cosiner/flag"
	"github.com/golang/snappy"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Client struct {
	Create      bool     `names:"-c" usage:"Compression file to snappy format"`
	Extract     bool     `names:"-x" usage:"Decompression snappy format file"`
	Silent      bool     `names:"-s, --silent" usage:"Silent mode" default:"false"`
	SourceFiles []string `names:"-f, --file" args:"true" usage:"Files to compression/decompression"`
	OutputDir   string   `names:"-o, --output" usage:"Output directory"`
}

func (c *Client) Metadata() map[string]flag.Flag {
	const (
		usage = `
		compression:   snappy -c example.txt
		decompression: snappy -x example.snappy
		`
		version = `
			version: v1.0.0
			commit:  xxx
			date:    2022-01-01 10:00:01
		`
		desc = `
		snappy-client is a compression/decompression client for Google Snappy write by golang.
		`
	)
	return map[string]flag.Flag{
		"": {
			Usage:   usage,
			Version: version,
			Desc:    desc,
		},
	}
}

func (c *Client) ParameterCheck() (bool, error) {
	var errMsg []string
	if c.Create && c.Extract {
		errMsg = append(errMsg, "compression and decompression cannot be specified at the same time")
	}

	if !(c.Create || c.Extract) {
		errMsg = append(errMsg, "compression or decompression must be specified one")
	}

	if len(c.SourceFiles) == 0 {
		errMsg = append(errMsg, "missing files to snappy or unsnappy")
	}

	if len(errMsg) > 0 {
		return false, errors.New(strings.Join(errMsg[:], "\n"))
	} else {
		return true, nil
	}
}

func (c *Client) run() {
	_, err := c.ParameterCheck()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c.ParseOutputDir()
	if c.Create {
		c.Snappy()
	} else if c.Extract {
		c.Unsnappy()
	} else {
		fmt.Println("unsupported operation")
	}
}

func (c *Client) ParseOutputDir() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if c.OutputDir == "" {
		c.OutputDir = wd
	}

	_, err = os.Stat(c.OutputDir)
	if os.IsNotExist(err) {
		os.MkdirAll(c.OutputDir, os.ModePerm)
	}
}

func (c *Client) GetOutputFilename(path string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	if c.Extract {
		if ext == "" {
			return strings.TrimSuffix(filename, ext) + ".unsnappy"
		} else {
			return strings.TrimSuffix(filename, ext)
		}
	} else if c.Create {
		return filename + ".snappy"
	} else {
		return filename
	}
}

func (c *Client) Snappy() {
	for _, f := range c.SourceFiles {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Println(err)
			continue
		}
		encoded := snappy.Encode(nil, content)
		fullFilename := path.Join(c.OutputDir, c.GetOutputFilename(f))
		err = ioutil.WriteFile(fullFilename, encoded, 0644)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (c *Client) Unsnappy() {
	for _, f := range c.SourceFiles {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Println(err)
			continue
		}
		decoded, err := snappy.Decode(nil, content)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fullFilename := path.Join(c.OutputDir, c.GetOutputFilename(f))
		err = ioutil.WriteFile(fullFilename, decoded, 0644)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
