// addtree is append ./tree.txt to ./README.md
//
// ```txt:./tree.txt
//
// <include tree.txt>
//
// ```
//
package main

import (
	"bufio"
	"fmt"
	"os"
)

func fatalIF(point string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s:%q", point, err)
		os.Exit(1)
	}
}

// default
const (
	// filename
	README = "./README.md"
	TREE   = "./tree.txt"

	// blocktag
	TREETAG  = "```txt:./tree.txt"
	TAGCLOSE = "```"
)

// README.md内treeBlockの判定
type treeBlock struct {
	in, exit bool
}

// ブロックの存在判定
func (b *treeBlock) exists() bool { return b.in && b.exit }

// ブロック位置の判定
func (b *treeBlock) search(s string) {
	b.setIn(s)
	b.setExit(s)
}
func (b *treeBlock) setIn(s string) {
	if b.in {
		return
	}
	if s == TREETAG {
		b.in = true
	}
}
func (b *treeBlock) setExit(s string) {
	if b.exit {
		return
	}
	if b.in && s == TAGCLOSE {
		b.exit = true
	}
}

// parse and split
func parseReadme(filename string) ([]string, error) {
	// file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("parseReadme:%q", err)
	}
	defer func() {
		fatalIF("parseReadme", file.Close())
	}()

	// Scan and parse
	block := new(treeBlock)
	str := make([]string, 2)
	for sc := bufio.NewScanner(file); sc.Scan(); {
		if err := sc.Err(); err != nil {
			return nil, fmt.Errorf("paseReadme:%q", err)
		}
		if !block.in {
			str[0] += fmt.Sprintln(sc.Text())
		}
		block.search(sc.Text())
		if block.exit {
			str[1] += fmt.Sprintln(sc.Text())
		}
	}
	if !block.exists() {
		return nil, fmt.Errorf("parseReadme:do not find tree block in README.md")
	}
	return str, nil
}

// get tree.txt buffer
func getTree(filename string) (string, error) {
	// file
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("getTree:%q", err)
	}
	defer func() {
		fatalIF("getTree", file.Close())
	}()

	// export
	var str string
	for sc := bufio.NewScanner(file); sc.Scan(); {
		if err := sc.Err(); err != nil {
			return "", fmt.Errorf("scanTree:%q", err)
		}
		str += fmt.Sprintln(sc.Text())
	}
	return str, nil
}

// join buffers
func joinText(readme []string, tree string) (string, error) {
	if len(readme) != 2 {
		return "", fmt.Errorf("joinText:parameter []string is invalid length")
	}
	return readme[0] + tree + readme[1], nil
}

// write
func writeReadme(filename string, s string) error {
	// file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("writeReadme:%q", err)
	}
	defer func() {
		fatalIF("writeReadme", file.Close())
	}()

	// writer
	w := bufio.NewWriter(file)
	n, err := w.WriteString(s)
	if err != nil {
		// TODO:backup and recover
		return fmt.Errorf("writeReadme:write bytes=%v\nerror:%q", n, err)
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("writeReadme:%q", err)
	}
	return nil
}

// 上書きの確認
func ask(str string) error {
	fmt.Println("----------| New README.md |----------")
	fmt.Println(str)
	fmt.Println("----------| New README.md |----------")
	fmt.Println("this string to override at README.md")
	fmt.Print("[yes:no]? >>")

L:
	for sc, i := bufio.NewScanner(os.Stdin), 0; i < 3 && sc.Scan(); i++ {
		fatalIF("ask", sc.Err())
		switch sc.Text() {
		case "yes":
			return nil
		case "no":
			break L
		default:
			fmt.Println(sc.Text())
			fmt.Print("[yes:no]? >>")
		}
	}
	fmt.Println()
	return fmt.Errorf("ask:don't write ...process exit")
}

func main() {
	type buffer struct {
		readme    []string
		tree      string
		newReadme string
	}
	var (
		buf buffer
		err error
	)

	buf.readme, err = parseReadme(README)
	fatalIF("main", err)

	buf.tree, err = getTree(TREE)
	fatalIF("main", err)

	buf.newReadme, err = joinText(buf.readme, buf.tree)
	fatalIF("main", err)

	err = ask(buf.newReadme)
	fatalIF("main", err)

	err = writeReadme(README, buf.newReadme)
	fatalIF("main", err)
}

