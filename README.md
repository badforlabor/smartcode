# smartcode
智能生成重复代码 generate code





# 思路

- 读取go语法的文件，生成AST
- 根据AST，解析出所有的结构体，以及结构体内部的成员变量
- 根据结构体信息，调用不同的模板，生成自己需要的代码（譬如，协议自动生成工具）





# 协议生成工具

### 使用方法

```
--template=tempo.t --code=code1.co --output=code1.out
```

- --template  代码模板。模板规则参考：https://golang.org/pkg/text/template/
- --code 代码，譬如定义了一堆结构体。代码使用的是Golang语法。
- --output 输出文件。如果为空，那么会输出到控制台



### 扩展模板

- {{NewLine}}  新的一样
- {{NIL}} 空
- {{Tab}} tab
- {{ToUpper}} 转化成**大写**，调用方式为 {{ToUpper .Name}}
- {{StripWord}} **裁减掉N个字母**，调用方式为 {{StripWord .Name 1}}
- {{ToSnake}} 转化成**Snake格式**（如：any_kind_of_string），调用方式为 {{ToSnake .Name}}
- {{ToCamel}} 转化成**Camel格式**（如：AnyKindOfString），调用方式为 {{ToCamel .Name}}
- {{ToLowerCamel}} 转化成**LowerCamel格式**（如：anyKindOfString），调用方式为 {{ToLowerCamel .Name}}