# 架构设计

`windranger`作为一个基于`yaml`的实体框架，总体上分为两个部分：
- 一个基于`yaml`语法解析器的前端
- 若干基于`go`模板的代码生成器后端

其中，每个后端的建议实现也分为两个部分:
- 根据后端特点对接口返回包结构进行链接处理的链接器
- 以链接器输出为输入的代码生成器

## 接口定义
`windranger`输入为若干`yaml`文件，经过前端解析后，产生包含完整信息的`[]*Package`，保证字典序
包类型
> 包类型包含`包名`、`枚举类型`和`结构类型`，`枚举类型`和`结构类型`保证字典序
> ```go
> type Package struct {
>     Name string
>     Enums []*Enum
>     Structures []*Structure
>     Dependencies []string
> }
> ```
枚举类型
> 枚举类型包含`枚举名`、`类型注释`和`枚举字段`
> ```go
> type Enum struct {
>     Name string
>     Comment string
>     EnumFields []*EnumField
> }
> ```
> 枚举字段包括`原始定义`和`字段注释`
> ```go
> type EnumField struct {
>     Name string
>     Comment string
> }
> ```
结构类型
> 结构类型包含`结构名`、`类型注释`和`字段`
> ```go
> type Structure struct {
>    Name string
>    Comment string
>    Fields []*Field
> }
> ```
> 字段包含`字段名`、`字段注释`和`类型`
> ```go
> type Field struct {
>     Name string
>     Comment string
>     Type *Type
> }
> ```
> 类型包含`类型名`、`属性`
> ```go
> type Type struct {
>     Raw string
>     Name string
>     Kind Kind
>     Package string
> }
> ```
> 属性为枚举类型
> ```go
> type Kind uint32
> const (
>     KindNormal Kind = iota
>     KindArray
>     KindOptional
>     KindPrimaryKey
> )
> ```
公共包名称
```go
const CommonPackage = "common"
```
## 前端设计
前端从文件或网络上接收一个或多个`yaml`文件；对每个`yaml`块分别进行解析，输出`Package`结构；最后对公共结构进行合并，生成`Info`

## 后端设计

### 链接器

链接器接收包信息和类型映射，完成类型链接，并返回每个包的依赖包列表，保证字典序
