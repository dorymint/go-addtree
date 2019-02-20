// include text to after matched string
//
//	$ goins -src /path/file -match "Tag" -ins "srouce strings"
//
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	Name    = "goinc"
	Version = "0.0.0"
)

func include(srcR io.Reader, match string, incR io.Reader) (*bytes.Buffer, error) {
	b, err := ioutil.ReadAll(incR)
	if err != nil {
		return nil, err
	}

	var s []string
	matched := false
	sc := bufio.NewScanner(srcR)
	for sc.Scan() {
		if sc.Err() != nil {
			return nil, sc.Err()
		}
		if !matched && match == sc.Text() {
			s = append(s, sc.Text(), string(b))
			matched = true
		} else {
			s = append(s, sc.Text())
		}
	}
	if !matched {
		return nil, errors.New("not matched")
	}

	return bytes.NewBufferString(strings.Join(s, "\n")), nil
}

func isRegular(f *os.File) error {
	if fi, err := f.Stat(); err != nil {
		return err
	} else if !fi.Mode().IsRegular() {
		return errors.New("not regular file")
	}
	return nil
}

func PrintStdout(file string, match string, incR io.Reader) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = isRegular(f); err != nil {
		return err
	}

	buf, err := include(f, match, incR)
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, buf)
	return err
}

func ForceInclude(file string, match string, incR io.Reader) error {
	f, err := os.OpenFile(file, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = isRegular(f); err != nil {
		return err
	}

	buf, err := include(f, match, incR)
	if err != nil {
		return err
	}

	// impl write backup?
	if err = f.Truncate(0); err != nil {
		return err
	}
	if _, err = f.Seek(0, 0); err != nil {
		return err
	}
	_, err = io.Copy(f, buf)
	return err
}

var opt struct {
	help    bool
	version bool

	force bool

	file  string
	match string

	includeString string
	includeFile   string // TODO
}

// Name string for specify command name
func makeUsage(w *io.Writer) func() {
	return func() {
		flag.CommandLine.SetOutput(*w)
		fmt.Fprintf(*w, "Description:\n")
		fmt.Fprintf(*w, "  Short description\n\n")
		fmt.Fprintf(*w, "Usage:\n")
		fmt.Fprintf(*w, "  %s [Options]\n\n", Name)
		fmt.Fprintf(*w, "Options:\n")
		flag.PrintDefaults()
		examples := `
Examples:
  $ ` + Name + ` -help # Display help message
  $ ` + Name + ` -file /path/file -match "Tag:" -string "include strings" # include strings after matched
  $ ` + Name + ` -file /path/file -match "Tag:" -inc /path/file           # include file contents after matched
`
		fmt.Fprintf(*w, "%s\n", examples)
	}
}

func main() {
	var usageWriter io.Writer = os.Stdout
	usage := makeUsage(&usageWriter)
	flag.Usage = usage
	flag.BoolVar(&opt.help, "help", false, "Display this message")
	flag.BoolVar(&opt.version, "version", false, "Display version")
	flag.BoolVar(&opt.force, "force", false, "Enable Overwrite")
	flag.StringVar(&opt.match, "match", "", "Specify target match string")
	flag.StringVar(&opt.includeString, "string", "", "Specify include string")
	flag.StringVar(&opt.includeFile, "inc", "", "Specify path to include file")
	flag.StringVar(&opt.file, "file", "", "Specify path to source")
	flag.Parse()
	if n := flag.NArg(); n == 1 && opt.file == "" {
		opt.file = flag.Arg(0)
	} else if n != 0 {
		usageWriter = os.Stderr
		flag.Usage()
		fmt.Fprintf(os.Stderr, "invalid arguments: %v\n", flag.Args())
		os.Exit(1)
	}

	switch {
	case opt.help:
		flag.Usage()
		os.Exit(0)
	case opt.version:
		fmt.Printf("%s %s\n", Name, Version)
		os.Exit(0)
	}

	if opt.file == "" {
		fmt.Fprintln(os.Stderr, "not specified target file")
		os.Exit(1)
	}

	if opt.includeString == "" && opt.includeFile == "" {
		fmt.Fprintln(os.Stderr, "not specified include target")
		os.Exit(1)
	}
	if opt.includeString != "" && opt.includeFile != "" {
		fmt.Fprintln(os.Stderr, "duplicate specified include target")
		os.Exit(1)
	}

	var incR io.Reader
	if opt.includeString != "" {
		incR = bytes.NewBufferString(opt.includeString)
	} else if opt.includeFile != "" {
		b, err := ioutil.ReadFile(opt.includeFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		incR = bytes.NewBuffer(b)
	}

	if opt.force {
		ForceInclude(opt.file, opt.match, incR)
	} else {
		PrintStdout(opt.file, opt.match, incR)
	}
}
