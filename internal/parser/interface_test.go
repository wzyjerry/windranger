package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeString(t *testing.T) {
	tp := &Type{
		Raw:     "datetime",
		Name:    "datetime",
		Kind:    KindPrimaryKey,
		Package: "time",
	}
	assert.Equal(t, "[KindPrimaryKey]time.datetime(datetime)", tp.String())
}

func TestFieldString(t *testing.T) {
	field := &Field{
		Name:    "id",
		Comment: "主键",
		Type: &Type{
			Raw:     "datetime",
			Name:    "datetime",
			Kind:    KindPrimaryKey,
			Package: "time",
		},
	}
	assert.Equal(t, "id [KindPrimaryKey]time.datetime(datetime)#主键", field.String())
}

func TestStructureString(t *testing.T) {
	structure := &Structure{
		Name:    "demo",
		Comment: "示例",
		Fields: []*Field{{
			Name:    "id",
			Comment: "主键",
			Type: &Type{
				Raw:  "string",
				Name: "string",
				Kind: KindPrimaryKey,
			},
		}, {
			Name:    "name",
			Comment: "名称",
			Type: &Type{
				Raw:  "string",
				Name: "string",
				Kind: KindNormal,
			},
		}},
	}
	assert.Equal(t,
		`#示例
type demo struct {
	id [KindPrimaryKey]string(string)#主键
	name [KindNormal]string(string)#名称
}`, structure.String())
}

func TestEnumFieldString(t *testing.T) {
	enumField := &EnumField{
		Name:    "array",
		Comment: "数组",
	}
	assert.Equal(t, `array#数组`, enumField.String())
}

func TestEnumString(t *testing.T) {
	enum := &Enum{
		Name:    "kind",
		Comment: "类型",
		EnumFields: []*EnumField{{
			Name:    "array",
			Comment: "数组",
		}, {
			Name:    "optional",
			Comment: "可选",
		}},
	}
	assert.Equal(t,
		`#类型
type kind enum {
	array#数组
	optional#可选
}`, enum.String())
}

func TestPackageString(t *testing.T) {
	pack := &Package{
		Name: "demo",
		Enums: []*Enum{{
			Name:    "kind",
			Comment: "类型",
			EnumFields: []*EnumField{{
				Name:    "array",
				Comment: "数组",
			}, {
				Name:    "optional",
				Comment: "可选",
			}},
		}, {
			Name:    "gender",
			Comment: "性别",
			EnumFields: []*EnumField{{
				Name:    "male",
				Comment: "男",
			}, {
				Name:    "female",
				Comment: "女",
			}},
		}},
		Structures: []*Structure{{
			Name:    "demo",
			Comment: "示例",
			Fields: []*Field{{
				Name:    "id",
				Comment: "主键",
				Type: &Type{
					Raw:  "string",
					Name: "string",
					Kind: KindPrimaryKey,
				},
			}, {
				Name:    "name",
				Comment: "名称",
				Type: &Type{
					Raw:  "string",
					Name: "string",
					Kind: KindNormal,
				},
			}, {
				Name:    "author",
				Comment: "作者列表",
				Type: &Type{
					Raw:  "author",
					Name: "author",
					Kind: KindArray,
				},
			}},
		}, {
			Name:    "author",
			Comment: "作者",
			Fields: []*Field{{
				Name:    "id",
				Comment: "主键",
				Type: &Type{
					Raw:  "string",
					Name: "string",
					Kind: KindPrimaryKey,
				},
			}, {
				Name:    "name",
				Comment: "姓名",
				Type: &Type{
					Raw:  "string",
					Name: "string",
					Kind: KindNormal,
				},
			}, {
				Name:    "gender",
				Comment: "性别",
				Type: &Type{
					Raw:  "gender",
					Name: "gender",
					Kind: KindOptional,
				},
			}},
		}},
		Dependencies: []string{"time"},
	}
	assert.Equal(t,
		`package demo
time
#类型
type kind enum {
	array#数组
	optional#可选
}
#性别
type gender enum {
	male#男
	female#女
}
#示例
type demo struct {
	id [KindPrimaryKey]string(string)#主键
	name [KindNormal]string(string)#名称
	author [KindArray]author(author)#作者列表
}
#作者
type author struct {
	id [KindPrimaryKey]string(string)#主键
	name [KindNormal]string(string)#姓名
	gender [KindOptional]gender(gender)#性别
}`, pack.String())
}
