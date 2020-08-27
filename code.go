/**
 * Auth :   liubo
 * Date :   2020/9/1 10:05
 * Comment: 根据解析出的AST，读取Tempate，生成对应的代码
 */

package main

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func Tab() string {
	return "\t"
}
func NewLine() string {
	return "\n"
}
func NIL() string {
	return ""
}
func ToUpper(s string) string  {
	return strings.ToUpper(s)
}

// 裁减掉前N个字符
func StripWord(s string, n int) string {
	return s[n:]
}
func ToSnake(args ...interface{}) string {
	return strcase.ToSnake(args[0].(string))
}

func GenerateCode(pack *PackageInfo) {
	if pack == nil {
		panic("非法的pack")
	}

	var filename = *templateFile

	t := template.New("cpp template")

	// 扩展模板
	t = t.Funcs(template.FuncMap{"Tab": Tab})
	t = t.Funcs(template.FuncMap{"NewLine": NewLine})
	t = t.Funcs(template.FuncMap{"NIL": NIL})
	t = t.Funcs(template.FuncMap{"ToUpper": ToUpper})
	t = t.Funcs(template.FuncMap{"StripWord": StripWord})

	// 美化代码
	t = t.Funcs(template.FuncMap{"ToSnake": ToSnake})	// 格式为：any_kind_of_string
	t = t.Funcs(template.FuncMap{"ToCamel": strcase.ToCamel})	// 格式为：AnyKindOfString
	t = t.Funcs(template.FuncMap{"ToLowerCamel": strcase.ToLowerCamel})	// 格式为：anyKindOfString



	var code = `
template<>
inline bool IsRepEqual<{{.Name}}>(const {{.Name}}& A, const {{.Name}}& B)
{
{{ $length := len .MemberList -}}
{{- if lt $length 3 -}}
    {{Tab}}return false;
{{- else -}}
    {{- Tab}}return {{NIL -}}
    {{- range $index, $element := .MemberList -}}
        {{-  if eq $index 2 }}IsRepEqual(A.$element, B.$element){{end -}}
        {{-  if gt $index 2  -}}
            {{NewLine}}{{Tab}}{{Tab}}&& IsRepEqual(A.$element, B.$element){{NIL -}}
        {{- end -}}
     {{- end -}};
{{- end -}}
{{- NewLine -}}
}
`

	var d, e = ioutil.ReadFile(filename)
	if e == nil && len(d) > 0 {
		code = string(d)
	}

	t, _ = t.Parse(code)

	var b = bytes.NewBufferString("")
	for _, v := range pack.StructList {
		t.Execute(b, &v)
	}

	if len(*outputFile) > 0 {
		var e = ioutil.WriteFile(*outputFile, b.Bytes(), os.ModePerm)
		if e != nil {
			fmt.Printf("输出文件失败:%s, e=%s\n", *outputFile, e.Error())
		}
	} else {
		fmt.Println(b.String())
	}
}
