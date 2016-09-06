# go-addtree
---
go-addtree is adding 'tree.txt' to 'README.md'

## Tree
use `addtree`. included from tree.txt

```txt:./tree.txt

./
├── LICENSE
├── README.md
├── addtree
│   └── addtree.go
└── tree.txt

1 directory, 4 files

```

## Installation
---
`go get github.com/dorymint/go-addtree/addtree`

## Usage
---
```txt:./tree.txt  
$ cd <your repository root>
$ echo '```txt:./tree.txt\n```\n' >> README.md
$ tree ./ > ./tree.txt
$ addtree
$ rm ./tree.txt
```  

1. change current directory to your repository root
2. adding tree text tag your README.md  
tree text tag is this 2line  
~~~
\`\`\`txt:./tree.txt
\`\`\`
~~~

3. use `tree`. generate tree.txt
4. after `addtree`. "tree.txt" is added between tree text tag

## Licence
---
MIT
