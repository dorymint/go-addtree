package main

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)


// directory
var (
	TESTDIR = "testdir"
	PWD     = ""
)

func init() {
	if err := os.Chdir(TESTDIR); err != nil {
		log.Fatal(err)
	}

	var err error
	PWD, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
}

func tempFile(content string) (string, error) {
	// use os.TempDir()
	f, err := ioutil.TempFile("", "readme")
	if err != nil {
		return "", err
	}
	defer func() {
		if errclose := f.Close(); errclose != nil {
			log.Fatal(errclose)
		}
	}()
	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}
	name := f.Name()
	return name, nil
}

func TestParseReadme(t *testing.T) {
	// README contents
	var (
		before   = "### README"
		tag      = "```txt:./tree.txt"
		inline   = "delete target"
		tagclose = "```"
		after    = "### after"
	)

	fatalTest := func(filename string) {
		readme, errFatalTest := parseReadme(filename)
		if errFatalTest == nil {
			t.Fatalf("expected return error but nil")
		} else {
			t.Logf("Retrun Error:%q", errFatalTest)
		}
		var expected []string
		if !reflect.DeepEqual(readme, expected) {
			t.Fatalf("expected %q but %q", expected, readme)
		}
	}

	fatalList := []string{
		"",
		tag,
		tag + "\n" + tag,
		tag + tagclose,
		tagclose,

		before + tag + "\n" + tagclose,
		tag + "\n" + tagclose + after,

		before + "\n" + tag + "\n" + after,
		before + "\n" + tagclose + "\n" + tag + "\n" + after,
	}

	t.Log("Start Fatal Test")
	for i, x := range fatalList {
		filename, err := tempFile(x)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Fatal Test %v", i+1)
		fatalTest(filename)
	}
	t.Log("End of Fatal Test")

	matchTest := func(filename string, expected []string) {
		out, errMatchTest := parseReadme(filename)
		if errMatchTest != nil {
			t.Fatal(errMatchTest)
		}
		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("\nexpected:\n%q\nbut:\n%q", expected, out)
		}
	}

	matchList := []struct {
		input    string
		expected []string
	}{
		{input: tag + "\n" + tagclose,
			expected: []string{
				tag + "\n", tagclose + "\n"}},

		{input: tag + "\n" + tagclose + "\n",
			expected: []string{
				tag + "\n", tagclose + "\n"}},

		{input: tag + "\n" + tagclose + "\n" + "\n",
			expected: []string{
				tag + "\n", tagclose + "\n" + "\n"}},

		{input: before + "\n" + tag + "\n" + tagclose + "\n" + after,
			expected: []string{
				before + "\n" + tag + "\n", tagclose + "\n" + after + "\n"}},

		{input: tag + "\n" + inline + "\n" + tagclose,
			expected: []string{
				tag + "\n", tagclose + "\n"}},

		{input: tag + "\n" + "\n" + tagclose,
			expected: []string{
				tag + "\n", tagclose + "\n"}},

		{input: tag + "\n" + tag + "\n" + tagclose,
			expected: []string{
				tag + "\n", tagclose + "\n"}},

		{input: tag + "\n" + tagclose + "\n" + tagclose,
			expected: []string{
				tag + "\n", tagclose + "\n" + tagclose + "\n"}},

		{input: tag + "\n" + tag + "\n" + tagclose + "\n" + tagclose,
			expected: []string{
				tag + "\n", tagclose + "\n" + tagclose + "\n"}},
	}

	t.Log("Start Match Test")
	for i, x := range matchList {
		filename, err := tempFile(x.input)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Match Test %v", i+1)
		matchTest(filename, x.expected)
	}
	t.Log("End of Match Test")
}
