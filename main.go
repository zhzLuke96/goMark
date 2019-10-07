// ![size badge](https://img.shields.io/badge/build-passing-green)
// generating markdown documents from .go files.
//
// # Overview
// > Copyright 2019 zhzluke96.  All rights reserved.
// > Use of this source code is governed by a GPL3.0
// > license that can be found in the LICENSE file.
//
// 类似 godoc 一样的功能，根据源码生成文档
//
// # Usage
/*
$ go install github.com/zhzLuke96/goMark
*/
// ```
// -f string
//    parse file name. (default "./main.go")
// -o string
//    output markdown file name. (default "./README.md")
// -t string
//    markdown title. (default "goMark📑")
// ```
// 需要注意， `/*...*/` 和 `//...` 会区别对待，长注释将认为是 golang 示例代码，反之则是正常的文档
// # TODO
// - [ ] 支持更多模式
// - [ ] 根据类分级
// - [ ] 显示代码位置
// - [ ] 生成API文档
// # LICENSE
// GPL-3.0
//
// <center>👇👇following is automatic content generation👇👇</center>
//

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

var filePth = flag.String("f", "./main.go", "parse file name.")
var output = flag.String("o", "./README.md", "output markdown file name.")
var title = flag.String("t", "goMark📑", "markdown title.")

func init() {
	flag.Parse()
}

func main() {
	file, err := loadFile(*filePth)
	if err != nil {
		fmt.Println(err)
		return
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", file, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	mk := markFile(f, file)
	saveFile([]byte(mk), *output)
}

// Example 是一个示例函数，写在函数之前的注释会被编译到文档中
//
// 它接受一个 name 字符串，并把它作为名字打印出来
// eg.
/*
Example("world")
// => you name is world.

err := Example("")
// =>
fmt.Println(err)
// =>  name string empty!
*/
// 被长注释包括的内容，将视作 golang 代码
func Example(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("Name string empty")
	}
	fmt.Printf("you name is %s.", name)
	return nil
}

func saveFile(fb []byte, pth string) error {
	newFile, err := os.Create(pth)
	if err != nil {
		return err
	}
	defer newFile.Close()
	if _, err := newFile.Write(fb); err != nil {
		return err
	}
	return nil
}

func loadFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func trimSpace(s string) string {
	s = strings.Trim(s, "\n")
	if len(s) == 0 {
		return s
	}
	if s[:1] == " " {
		return s[1:]
	}
	return s
}

func markDoc(ds []*ast.Comment) string {
	ret := ""
	for _, v := range ds {
		if strings.HasPrefix(v.Text, "//") {
			d := strings.TrimPrefix(v.Text, "//")
			d = trimSpace(d)
			ret += d
			if !strings.HasSuffix(ret, "<br>\n") && !strings.HasSuffix(d, "```") {
				ret += "<br>\n"
			} else {
				ret += "\n"
			}
		}
		if strings.HasPrefix(v.Text, "/*") && strings.HasSuffix(v.Text, "*/") {
			d := strings.TrimPrefix(v.Text, "/*")
			d = strings.TrimSuffix(d, "*/")
			d = trimSpace(d)
			ret += "\n```golang\n" + d + "\n```\n"
		}
	}
	return ret
}

func markFile(f *ast.File, fbs []byte) string {
	head := ""
	if *title == "." {
		head = "# " + f.Name.Name + "\n"
	} else {
		head = "# " + *title + "\n"
	}

	if f.Comments != nil {
		if f.Comments[0].Pos() == 1 {
			head += markDoc(f.Comments[0].List)
		}
	} else {
		head += "\n"
	}

	typeList := []string{}
	funcList := []string{}

	body := ""
	for _, v := range f.Decls {
		switch vt := v.(type) {
		case *ast.GenDecl:
			if vt.Tok.String() != "type" || vt.Doc == nil {
				continue
			}
			name := vt.Specs[0].(*ast.TypeSpec).Name.Name
			title := "[type] " + name
			typeList = append(typeList, title)
			body += `# ` + title + "\n\n"
			body += "```golang\n" + string(fbs[vt.Pos()-1:vt.End()-1]) + "\n```\n"
			body += markDoc(vt.Doc.List)
		case *ast.FuncDecl:
			if vt.Doc == nil {
				continue
			}
			var name, fnt string
			if vt.Name != nil {
				name = vt.Name.String()
			}
			if vt.Type != nil {
				fnt = string(fbs[vt.Type.Pos()-1 : vt.Type.End()-1])
			}

			title := "[func] " + name
			funcList = append(funcList, title)
			body += "# " + title + "\n" + ">" + fnt + "\n\n"
			body += markDoc(vt.Doc.List)
		default:
			fmt.Println("pass", vt)
		}
		body += "[go top](#index)\n\n"
	}

	if len(typeList) != 0 || len(funcList) != 0 {
		head += "# Index\n"
	}
	if len(typeList) != 0 {
		head += "TYPE: \n"
		for _, idx := range typeList {
			titleCur := strings.Replace(idx, " ", "-", -1)
			head += fmt.Sprintf("- [%s](#%s)\n", idx, titleCur)
		}
	}
	if len(funcList) != 0 {
		head += "\nFUNC:\n"
		for _, idx := range funcList {
			titleCur := strings.Replace(idx, " ", "-", -1)
			head += fmt.Sprintf("- [%s](#%s)\n", idx, titleCur)
		}
	}

	return head + "\n\n" + body
}
