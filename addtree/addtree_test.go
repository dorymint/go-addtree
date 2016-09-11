package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// TODO:create Tests

func pwd() (string) {
	dir, err := os.Getwd()
	if err != nil {
		 log.Fatal(fmt.Errorf("pwd error\n"))
	}
	fmt.Println(dir)
	return dir
}

func cdTestdir() {
	pwd()
	if err := os.Chdir("./testdir/"); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error test stoped!!\n:%q\n", err)
		os.Exit(2)
	}
	pwd()
}

// TODO:
func TestGetReadme(t *testing.T) {
	cdTestdir()
	makeReadme()

	// TODO:
	readme := make([]string, 2)
	var err error
	if readme, err = getReadme(); err != nil {
		t.Errorf("test error, %q\n%q\n", readme, err)
	}
	t.Logf("pass %q\n%q\n", readme[0], readme[1])
}

// README.md data define
const (
	before   = "### Readme test\n"
	after    = "after tag"
	tag      = "\n```txt:./tree.txt\n"
	tagclose = "\n```\n"
)

var ReadmeTests = []struct {
	in  string // imitates ./README.md // data write to mock
	out []string
}{
	{before + after, nil},
	{before + tag + after, nil},
	{before + tag + tagclose + after, []string{before , after} },
	{before + tagclose + after, nil},
} // ReadmeTests

func makeReadme() {

	// TODO:
	for _, x := range ReadmeTests {
		fmt.Println(x)
		if x.out != nil {
			fmt.Printf("x print%s%s\n", x.out[0], x.out[1])
		}
	}

}
