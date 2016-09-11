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
// askでREADME.md上書きの確認を取る
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
	beginLine, endLine int
	in, exit           bool
}

// README.md 内でtreeBlockを見つけていればtrueを返す
func (b *treeBlock) exists() bool { return b.in && b.exit }

// block開始位置の判定と行番号の記録
func (b *treeBlock) setBeginLine(s string, line int) {
	if b.in {
		return
	}
	if s == TREETAG {
		b.in = true
		b.beginLine = line
	}
}
// block終了位置の判定と行番号の記録
func (b *treeBlock) setEndLine(s string, line int) {
	if b.exit {
		return
	}
	if s == TAGCLOSE && b.in {
		b.endLine = line
		b.exit = true
	}
}
func (b *treeBlock) searchAndSet(s string, line int) {
	b.setBeginLine(s, line)
	b.setEndLine(s, line)
}

// README.md parse and split
// TODO:TEST
func getReadme() ([]string, error) {

	buf := make([]string, 0, 256)
	block := new(treeBlock)
	// file open reamde
	file, err := os.Open(README)
	if err != nil {
		return nil, fmt.Errorf("getReadme(): %q\n", err)
	}
	defer func() {
		err := file.Close()
		if err != nil { log.Fatal(err) }
	}()

	// Scanner readme
	fmt.Println("Scanning README.md")
	for sc, i := bufio.NewScanner(file), 0; sc.Scan(); i++ {
		if err := sc.Err(); err != nil {
			return nil, fmt.Errorf("getReadme(): %q\n", err)
		}

		block.searchAndSet(sc.Text(), i)

		buf = append(buf, fmt.Sprintln(sc.Text()))
	}
	if !block.exists() {
		return nil, fmt.Errorf("getReadme(): not find tree block in README.md\n")
	}

	// return string
	str := make([]string, 2)
	for _, s := range buf[:block.beginLine+1] {
		str[0] += s
	}
	for _, s := range buf[block.endLine:] {
		str[1] += s
	}
	return str, nil
}

// get tree.txt buffer
// TODO:TEST
func getTree() (string, error) {

	// TODO:ERROR出るかも取り敢えずテストまでまつ
	tree := make([]string, 0, 256)
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
		tree = append(tree, fmt.Sprintln(sc.Text()))
	}

	// return string
	var str string
	for _, s := range tree {
		str += s
	}
	return str, nil
}

// join buffers
// TODO:TEST
func joinText(readme []string, tree string) (string, error) {
	if len(readme) != 2 {
		return "", fmt.Errorf("joinText(): parameter []string is invalid length\n")
	}
	str := readme[0]
	str += tree
	str += readme[1]
	return str, nil
}

// Write Readme
// TODO:TEST
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
// TODO:TEST
func ask(str string) {
	fmt.Println(str)
	fmt.Println("this string to override at README.md")
	fmt.Printf("[yes:no]? >>")
	for sc, i := bufio.NewScanner(os.Stdin), 0; i < 3 && sc.Scan(); i++ {
		switch sc.Text() {
		case "yes":
			return
		case "no":
			return
		default:
			fmt.Println(sc.Text())
			fmt.Printf("[yes:no]? >>")
		}
	}
	fmt.Fprintf(os.Stderr, "\n\ndon't write ...process exit\n")
	os.Exit(1)
}

func fatalIF(err error) {
	if err != nil { log.Fatal(err) }
}

func main() {
	type buffer struct {
		readme	[]string
		tree	string
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
