package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInclude(t *testing.T) {
	var (
		contents    = "head\nTag:\ntail"
		includeText = "include text"
		match       = "Tag:"
		exp         = "head\nTag:\n" + includeText + "\ntail"
	)

	var srcR io.Reader = bytes.NewBufferString(contents)
	var incR io.Reader = bytes.NewBufferString(includeText)
	buf, err := include(srcR, match, incR)
	if err != nil {
		t.Fatal(err)
	}

	if exp != buf.String() {
		t.Errorf("exp:%v", exp)
		t.Errorf("out:%v", buf.String())
	}
}

// filename, remove := setup(t, contents)
// defer remove()
func setup(t *testing.T, contents []byte) (file string, rm func()) {
	t.Helper()
	dir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	file = filepath.Join(dir, "file.txt")
	if err := ioutil.WriteFile(file, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
	rm = func() {
		os.Remove(file)
		os.RemoveAll(dir)
	}
	return
}

func TestForceInclude(t *testing.T) {
	var (
		contents    = "head\nTag:\ntail"
		includeText = "include text"
		match       = "Tag:"
		exp         = "head\nTag:\n" + includeText + "\ntail"
	)
	file, remove := setup(t, []byte(contents))
	defer remove()

	var incR io.Reader = bytes.NewBufferString(includeText)
	if err := ForceInclude(file, match, incR); err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	out := string(b)
	if exp != out {
		t.Errorf("exp:%q", exp)
		t.Errorf("out:%q", out)
	}
}
