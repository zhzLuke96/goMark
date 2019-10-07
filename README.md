# goMark
![size badge](https://img.shields.io/badge/build-passing-green)<br>
> generating markdown documents from .go files.<br>

# Overview<br>
> Copyright 2019 zhzluke96.  All rights reserved.<br>
> Use of this source code is governed by a GPL3.0<br>
> license that can be found in the LICENSE file.<br>

类似 godoc 一样的功能，根据源码生成文档<br>

# Usage<br>

```golang
$ go install github.com/zhzLuke96/goMark
```
```
-f string<br>
___parse file name. (default "./main.go")<br>
-o string<br>
___output markdown file name. (default "./README.md")<br>
-t string<br>
___markdown title. (default "goMark")<br>
```
需要注意, `/*...*/` 和 `//...` 会区别对待，长注释将认为是 golang 示例代码，反之则是正常的文档<br>
# TODO<br>
- [ ] 支持更多模式<br>
- [ ] 根据类分级<br>
- [ ] 显示代码位置<br>
- [ ] 生成API文档<br>
# LICENSE<br>
GPL-3.0<br>

<div style="text-align:center;">👇👇following is automatic content generation👇👇</div><br>

# Index

FUNC:
- [[func] Example](#[func]-Example)


# [func] Example
>func Example(name string) error

Example 是一个示例函数，写在函数之前的注释会被编译到文档中<br>

它接受一个 name 字符串，并把它作为名字打印出来<br>
eg.<br>

```golang
Example("world")
// => you name is world.

err := Example("")
// =>
fmt.Println(err)
// =>  name string empty!
```
被长注释包括的内容，将视作 golang 代码<br>
[go top](#index)

