# go-addtree
go-addtree is adding 'tree.txt' to 'README.md'

---
## Used addtree example
use `addtree`, includ from tree.txt


```txt:./tree.txt

./
├── LICENSE
├── README.md
├── addtree
│   └── addtree.go
└── tree.txt

1 directory, 4 files

```


---
## Installation
`go get github.com/dorymint/go-addtree/addtree`

---
## Usage
```txt:./tree.txt  
$ cd <your repository root>
$ echo '```txt:./tree.txt\n```\n' >> README.md
$ tree ./ > ./tree.txt
$ addtree
$ rm ./tree.txt
```  

1. change current directory to your repository root
2. adding tree-tag your README.md  
tree-tag is this 2 line  
```` ```txt:./tree.txt ````  
```` ``` ````

3. use `tree . > ./tree.txt` generate tree.txt
4. after `addtree`, "tree.txt" is added between tree-tag

---
## Licence
MIT
