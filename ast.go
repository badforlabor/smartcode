/**
 * Auth :   liubo
 * Date :   2020/8/27 14:29
 * Comment: 收集结构体信息
 */

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

// 包
type PackageInfo struct {
	CommentList []string
	Name        string
	StructList  []StructInfo
}

// 结构体
type StructInfo struct {
	Comment    []string
	Name       string
	MemberList []MemberInfo
}

// 成员变量
type MemberInfo struct {
	Comment      []string
	VariableName string
	VariableType string
	Tag          []TagInfo
}

// 成员变量的tag
type TagInfo struct {
	Key   string
	Value string
}

/*
	解析文件
	遍历AST，收集出结构体的名字、注释，结构体的成员变量、扩展属性、注释。
*/
func ParseAST(filename string) *PackageInfo {
	{
		var _, err = os.Stat(filename)
		if err != nil {
			fmt.Println("无法找到文件", filename)
			panic("")
		}
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		log.Println(err)
		panic("")
	}

	// 打印语法树
	//ast.Print(fset, f)

	var pack = &PackageInfo{}
	pack.Name = f.Name.Name
	pack.CommentList = ToComment2(f.Doc)

	for _, d := range f.Decls {
		var decls, ok = d.(*ast.GenDecl)
		if ok {
			var comments = ToComment2(decls.Doc)
			for _, t := range decls.Specs {
				var ts, ok = t.(*ast.TypeSpec)
				if ok {

					var st, ok = ts.Type.(*ast.StructType)
					if ok {

						var newSt StructInfo
						newSt.Comment = comments
						newSt.Name = ts.Name.Name

						for _, f := range st.Fields.List {

							var fComments = CombineCommentList(ToComment2(f.Comment), ToComment2(f.Doc))
							for _, fn := range f.Names {
								var member MemberInfo
								member.Comment = fComments
								member.VariableName = fn.Name
								var mt, ok = f.Type.(*ast.Ident)
								if ok {
									member.VariableType = mt.Name
									if f.Tag != nil {
										member.Tag = ParseTag(f.Tag.Value)
									}

									newSt.MemberList = append(newSt.MemberList, member)
								} else {
									fmt.Println("无法识别:", f.Type)
									panic("")
								}
							}
						}

						pack.StructList = append(pack.StructList, newSt)
					} else {
						fmt.Println("无法识别:", ts.Type)
						panic("")
					}
				} else {
					fmt.Println("无法识别:", t)
					panic("")
				}
			}
		} else {
			fmt.Println("无法识别:", decls)
			panic("")
		}
	}

	return pack
}

func ParseTag(tag string) []TagInfo {
	var ret []TagInfo

	if strings.HasPrefix(tag, "`") && strings.HasSuffix(tag, "`") {
		tag = tag[1 : len(tag)-1]
		var taglist = strings.Split(tag, " ")
		for _, v := range taglist {
			v = strings.TrimSpace(v)
			if len(v) > 0 {

				var kv = strings.Split(v, ":")
				if len(kv) == 2 {
					var one TagInfo
					one.Key = kv[0]
					one.Value = strings.Trim(kv[1], "\"")
					ret = append(ret, one)
				} else {
					fmt.Println("错误的tag:", v)
					panic("")
				}

			}
		}
	}

	return ret
}

func CombineCommentList(com1 []string, com2 []string) []string {
	return append(com1, com2...)
}

func ToComment(g []*ast.CommentGroup) []string {
	var ret []string
	for _, v := range g {
		for _, v2 := range v.List {
			ret = append(ret, v2.Text)
		}
	}
	return ret
}
func ToComment2(g *ast.CommentGroup) []string {
	var ret []string
	if g != nil {
		for _, v := range g.List {
			ret = append(ret, v.Text)
		}
	}
	return ret
}
