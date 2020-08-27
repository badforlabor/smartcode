/**
 * Auth :   liubo
 * Date :   2020/8/27 11:06
 * Comment:
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
)

var templateFile = flag.String("template", "", "-template=tempo.t")
var codeFile = flag.String("code", "", "-code=code1.co")
var outputFile = flag.String("output", "", "-output=code1.out")

func main() {

	flag.Parse()

	var pack = ParseAST(*codeFile)
	if pack != nil {
		fmt.Println("AST内容:")

		var d, e = json.MarshalIndent(pack, "", "    ")
		if e == nil {
			fmt.Println(string(d))
		} else {
			fmt.Println(*pack)
		}
	}

	// 将解析出的AST，生成代码
	GenerateCode(pack)
}
