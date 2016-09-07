/*
./tree.txtの内容を./README.mdに追加する

`go run addtree.go`

または、バイナリがあれば

`addtree`

追加位置はREADME.md内の以下の部分

```txt:./tree.txt

<追加位置>

```

README.md内に文字列 ```txt:./tree.txt が見つからなければ書き込まない

追加位置にある文字列は上書きされる
*/
package main

/* TODO LIST */
// DONE:カレントディレクトリからREADME.mdを掴んでbufferを作る
// test:bufferの内容をそのまま出力して確認

// DONE:buffer内を1行ずつ確認
// test:bufferを1行ずつ行番号を付けて出力して確認

// DONE:```txt:./tree.txt の前とその後の ``` 後でbufferを分割する string[2]
// DONE:```txt:./tree.txt ``` のブロックが見つからなければexitする
// test:parseしたbufferを2つずつ分けて出力して確認

// DONE:tree.txtを掴んでbuffer3を作る
// test:buffer3を出力して確認

// DONE:newBuffer := string[0]+blockBegin+buffer3+blockEnd+string[1]
// test:新しいbufferの内容を出力して確認

// DONE:新しいbufferをREADME.mdにテキストで出力する

import (
	"bufio"
	"fmt"
	"os"
)

// default
const (
	readme = "./README.md"
	tree   = "./tree.txt"

	beginBlock = "```txt:./tree.txt"
	endBlock   = "```"
)

type buffer struct {
	readme, tree []string
	block        treeBlock
}

func (buf *buffer) blockBegin() int { return buf.block.beginLine }
func (buf *buffer) blockEnd() int   { return buf.block.endLine }

// README.md 内でtreeBlockを見つけていればtrueを返す
func (buf *buffer) existsBlock() bool { return buf.block.in && buf.block.exit }

// README.md内のtreeBlockの位置
type treeBlock struct {
	beginLine, endLine int
	in, exit           bool
}

// block開始位置の判定と行番号の記録
func (b *treeBlock) setBeginLine(s string, line int) {
	if b.in {
		return
	}
	if s == beginBlock {
		b.in = true
		b.beginLine = line
	}
}

// block終了位置の判定と行番号の記録
func (b *treeBlock) setEndLine(s string, line int) {
	if b.exit {
		return
	}
	if s == endBlock && b.in {
		b.endLine = line
		b.exit = true
	}
}
func (b *treeBlock) searchBlock(s string, line int) {
	b.setBeginLine(s, line)
	b.setEndLine(s, line)
}

// README.md parse and split
func getReadme() ([]string, error) {

	buf := new(buffer)

	// file open reamde
	file, err := os.Open(readme)
	if err != nil {
		return nil, fmt.Errorf("getReadme(): %q\n", err)
	}
	defer file.Close()

	// Scanner readme
	fmt.Println("Scanning README.md")
	for sc, i := bufio.NewScanner(file), 0; sc.Scan(); i++ {
		if err := sc.Err(); err != nil {
			return nil, fmt.Errorf("getReadme(): %q\n", err)
		}

		buf.block.searchBlock(sc.Text(), i)

		buf.readme = append(buf.readme, fmt.Sprintln(sc.Text()))
	}
	if !buf.existsBlock() {
		return nil, fmt.Errorf("getReadme(): not find tree block in README.md\n")
	}

	// return string
	str := make([]string, 2)
	for _, s := range buf.readme[:buf.blockBegin()] {
		str[0] += s
	}
	for _, s := range buf.readme[buf.blockEnd()+1:] {
		str[1] += s
	}
	return str, nil
}

// get tree.txt buffer
func getTree() (string, error) {

	buf := new(buffer)

	// file open tree
	file, err := os.Open(tree)
	if err != nil {
		return "", fmt.Errorf("getTree(): %q\n", err)
	}
	defer file.Close()

	// scanner tree
	fmt.Println("Scanning tree.txt")
	for sc := bufio.NewScanner(file); sc.Scan(); {
		if err := sc.Err(); err != nil {
			return "", fmt.Errorf("getTree(): %q\n", err)
		}
		buf.tree = append(buf.tree, fmt.Sprintln(sc.Text()))
	}

	// return string
	var str string
	for _, s := range buf.tree {
		str += s
	}
	return str, nil
}

// join buffers
func joinText(readme []string, tree string) (string, error) {
	if len(readme) != 2 {
		return "", fmt.Errorf("joinText(): 'readme' invalid lnegth\n")
	}
	str := readme[0]
	str += fmt.Sprintf("\n%s\n\n", beginBlock)
	str += tree
	str += fmt.Sprintf("\n%s\n\n", endBlock)
	str += readme[1]
	return str, nil
}

// Write Readme
func writeReadme(s string) error {

	file, err := os.Create(readme)
	if err != nil {
		return fmt.Errorf("writeReadme(): %q\n", err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	n, err := w.WriteString(s)
	if err != nil {
		return fmt.Errorf("writeReadme(): write bytes=%v\nerror:%q\n", n, err)
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("writeReadme(): %q\n", err)
	}

	return nil
}

// 上書きの確認
func ask() bool {

	fmt.Println("this string to override at README.md")
	fmt.Printf("[yes:no]? >>")
	for sc, i := bufio.NewScanner(os.Stdin), 0; i < 3 && sc.Scan(); i++ {
		switch sc.Text() {
		case "yes":
			return true
		case "no":
			return false
		default:
			fmt.Println(sc.Text())
			fmt.Printf("[yes:no]? >>")
		}
	}

	fmt.Printf("\n\ndon't write ...process exit\n")
	return false
}

func main() {

	bufReadme := make([]string, 2)
	bufReadme, err := getReadme()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bufTree, err := getTree()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	str, err := joinText(bufReadme, bufTree)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(str)
	if !ask() {
		os.Exit(1)
	}

	if err := writeReadme(str); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
