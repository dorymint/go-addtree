go-include
=========

Include string after matched line

Usage
-----

```sh
# display to stdout
goinc -file /path/file -match "Tag:" -string "include word"

# overwrite
poinc -file /path/file -match "Tag:" -string "include word" -force
```

Install
-------

```sh
go get github.com/yaeshimo/go-include/goinc
```

License
-------

MIT
