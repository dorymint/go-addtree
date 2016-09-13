// ./tree.txtの内容を./README.mdに追加する
// `go run addtree.go`
// または、バイナリがあれば
// `addtree`
// 追加位置はREADME.md内の以下の部分
//
// ```txt:./tree.txt
// <追加位置>
// ```
//
// README.md内に文字列 ```txt:./tree.txt が見つからなければ書き込まない
// 追加位置にある文字列は上書きされる
package main

// 処理の流れ:
// カレントディレクトリからREADME.mdを掴んでbufferを作る
// buffer内を1行ずつ確認、TREETAGを探す
// ```txt:./tree.txt の前とその後の ``` 後でbufferを分割する [2]string
// treeBlockがREADME.mdに見つからなければexitする
// カレントディレクトリからtree.txtを掴んでtreebufferを作る
// newBuffer := string[0]+TREETAG+treebuffer+TAGCLOSE+string[1]
// README.md上書きの確認を取る
// newBufferをREADME.mdに上書きする

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// default
const (
	README = "./README.md"
	TREE   = "./tree.txt"

	TREETAG  = "```txt:./tree.txt"
	TAGCLOSE = "```"
)

// README.md内treeBlockの位置
type treeBlock struct {
	in, exit bool
}

func (b *treeBlock) exists() bool { return b.in && b.exit }

// ブロック位置の判定
func (b *treeBlock) searchAndSet(s string) {
	b.setIn(s)
	b.setOut(s)
}
func (b *treeBlock) setIn(s string) {
	if b.in {
		return
	}
	if s == TREETAG {
		b.in = true
	}
}
func (b *treeBlock) setOut(s string) {
	if b.exit {
		return
	}
	if s == TAGCLOSE && b.in {
		b.exit = true
	}
}

// README.md parse and split
func getReadme() ([]string, error) {

	block := new(treeBlock)
	str := make([]string, 2)

	// file open reamde
	file, err := os.Open(README)
	if err != nil {
		return nil, fmt.Errorf("getReadme(): %q\n", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Scanner readme
	fmt.Println("Scanning README.md")

	for sc := bufio.NewScanner(file); sc.Scan(); {
		if err := sc.Err(); err != nil {
			return nil, fmt.Errorf("getReadme(): %q\n", err)
		}
		tmp := sc.Text()

		if !block.in {
			str[0] += fmt.Sprintln(tmp)
		}

		block.searchAndSet(tmp)

		if block.exit {
			str[1] += fmt.Sprintln(tmp)
		}
	}
	if !block.exists() {
		return nil, fmt.Errorf("getReadme(): not find tree block in README.md\n")
	}
	return str, nil
}

// get tree.txt buffer
func getTree() (string, error) {
	var tree string
	// file open tree
	file, err := os.Open(TREE)
	if err != nil {
		log.Fatalf("getTree(): %q\n", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// scanner tree
	fmt.Println("Scanning tree.txt")
	for sc := bufio.NewScanner(file); sc.Scan(); {
		if err := sc.Err(); err != nil {
			return "", fmt.Errorf("getTree(): %q\n", err)
		}
		tree += fmt.Sprintln(sc.Text())
	}
	return tree, nil
}

// join buffers
func joinText(readme []string, tree string) (string, error) {
	if len(readme) != 2 {
		return "", fmt.Errorf("joinText(): parameter []string is invalid length\n")
	}
	return readme[0] + tree + readme[1], nil
}

// Write Readme
func writeReadme(s string) (err error) {
	file, err := os.Create(README)
	if err != nil {
		return fmt.Errorf("writeReadme(): %q\n", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// write
	w := bufio.NewWriter(file)
	n, err := w.WriteString(s)
	if err != nil {
		log.Fatalf("writeReadme(): write bytes=%v\nerror:%q\n", n, err)
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("writeReadme(): %q\n", err)
	}
	return nil
}

// 上書きの確認
func ask(str string) {
	fmt.Println(str)
	fmt.Println("this string to override at README.md")
	fmt.Print("[yes:no]? >>")

	L:
	for sc, i := bufio.NewScanner(os.Stdin), 0; i < 3 && sc.Scan(); i++ {
		switch sc.Text() {
		case "yes":
			return
		case "no":
			break L
		default:
			fmt.Println(sc.Text())
			fmt.Print("[yes:no]? >>")
		}
	}
	fmt.Fprintf(os.Stderr, "\n\ndon't write ...process exit\n")
	os.Exit(1)
}

func fatalIF(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

	buf.readme, err = getReadme()
	fatalIF(err)

	buf.tree, err = getTree()
	fatalIF(err)

	buf.newReadme, err = joinText(buf.readme, buf.tree)
	fatalIF(err)

	ask(buf.newReadme)

	err = writeReadme(buf.newReadme)
	fatalIF(err)
}
