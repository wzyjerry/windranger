package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindConflict(t *testing.T) {
	{
		fields := make([]*Field, 0)
		fields = append(fields, &Field{
			Name: "A",
		})
		fields = append(fields, &Field{
			Name: "B",
		})
		fields = append(fields, &Field{
			Name: "A",
		})
		conflict := findConflict(fields, func(field *Field) string {
			return field.Name
		})
		assert.Equal(t, []string{"A"}, conflict)
	}
	{
		fields := make([]*Field, 0)
		fields = append(fields, &Field{
			Name: "A",
		})
		fields = append(fields, &Field{
			Name: "A",
		})
		fields = append(fields, &Field{
			Name: "A",
		})
		conflict := findConflict(fields, func(field *Field) string {
			return field.Name
		})
		assert.Equal(t, []string{"A", "A"}, conflict)
	}
}

// func TestAddYaml(t *testing.T) {
// 	parser := NewParser()
// 	parser.AddYamlPath("https://dev.aminer.cn/gopkg/luna/-/raw/master/schema.yaml")
// 	_, err := parser.Parse()
// 	assert.Nil(t, err)
// }

func TestEmptyDoc(t *testing.T) {
	{
		parser := NewParser()
		parser.AddYaml([]byte{})
		packages, err := parser.Parse()
		assert.Nil(t, err)
		assert.Equal(t, `[]`, fmt.Sprintf("%v", packages))
	}
	{
		parser := NewParser()
		parser.AddYaml([]byte(`---`))
		packages, err := parser.Parse()
		assert.Nil(t, err)
		assert.Equal(t, `[]`, fmt.Sprintf("%v", packages))
	}
	{
		parser := NewParser()
		parser.AddYaml([]byte(
`---
---`))
		packages, err := parser.Parse()
		assert.Nil(t, err)
		assert.Equal(t, `[]`, fmt.Sprintf("%v", packages))
	}
	{
		parser := NewParser()
		parser.AddYaml([]byte(`23333`))
		packages, err := parser.Parse()
		assert.Nil(t, err)
		assert.Equal(t, `[]`, fmt.Sprintf("%v", packages))
	}
}

func TestEnum(t *testing.T) {
	{
		parser := NewParser()
		parser.AddYaml([]byte(
`# 性别
gender:
  - male # 男
  - female # 女
`))
		packages, err := parser.Parse()
		assert.Nil(t, err)
		assert.Equal(t, 
`[package common
#性别
type gender enum {
	male#男
	female#女
}]`, fmt.Sprintf("%v", packages))
	}
	{
		parser := NewParser()
		parser.AddYaml([]byte(
`# 类型
kind: [normal, array, optional, primary_key]`))
		packages, err := parser.Parse()
		assert.Nil(t, err)
		assert.Equal(t, 
`[package common
#类型
type kind enum {
	normal#
	array#
	optional#
	primary_key#
}]`, fmt.Sprintf("%v", packages))
	}
}

func TestEnumFault(t *testing.T) {
	{
		parser := NewParser()
		parser.AddYaml([]byte(
`# 性别
gender: [a, [b]]
`))
		_, err := parser.Parse()
		assert.Equal(t, []error{fmt.Errorf("枚举类型必须为标量")}, err)
	}
	{
		parser := NewParser()
		parser.AddYaml([]byte(
`# 性别
gender: [a, a]
`))
		_, err := parser.Parse()
		assert.Equal(t, []error{fmt.Errorf("重复的枚举值: a")}, err)
	}
}

func TestMapping(t *testing.T) {
	{
		parser := NewParser()
		parser.AddYaml([]byte(
`# 示例
demo: &table
  id!: string # 主键
  name: string # 名称
  update_at?: datetime # 更新日期
  author[]: # 作者
    id?: string # NAID
    name: string # 姓名
    gender: # 性别
      - unset # 未设置
      - male # 男
      - female # 女
`))
	packages, err := parser.Parse()
	assert.Nil(t, err)
	assert.Equal(t, 
`[package demo
#性别
type gender enum {
	unset#未设置
	male#男
	female#女
}
#作者
type author struct {
	id [KindOptional](string)#NAID
	name [KindNormal](string)#姓名
	gender [KindNormal](gender)#性别
}
#示例
type demo struct {
	id [KindPrimaryKey](string)#主键
	name [KindNormal](string)#名称
	update_at [KindOptional](datetime)#更新日期
	author [KindArray](author)#作者
}]`, fmt.Sprintf("%v", packages))
	}
}

func TestMappingFatual(t *testing.T) {
	parser := NewParser()
	parser.AddYaml([]byte(
`# 示例
demo: &table
  id!: string # 主键
  name: string # 名称
  update_at?: datetime # 更新日期
  author[]: # 作者
    id?: string # NAID
    name: string # 姓名
    name: # 性别
      - unset # 未设置
      - male # 男
      - female # 女
`))
	_, err := parser.Parse()
	assert.Equal(t, []error{fmt.Errorf("重复的字段名: name")}, err)
}

func TestPackage(t *testing.T) {
	parser := NewParser()
	parser.AddYaml([]byte(
`# 性别
gender:
  - unset # 未设置
  - male # 男
  - female # 女
---
# 示例
demo: &table
  id!: string # 主键
  name: string # 名称
  update_at?: datetime # 更新日期
  author[]: # 作者
    id?: string # NAID
    name: string # 姓名
    gender: gender # 作者性别
`))
		packages, err := parser.Parse()
		assert.Nil(t, err)
		assert.Equal(t, 
`[package common
#性别
type gender enum {
	unset#未设置
	male#男
	female#女
} package demo
#作者
type author struct {
	id [KindOptional](string)#NAID
	name [KindNormal](string)#姓名
	gender [KindNormal](gender)#作者性别
}
#示例
type demo struct {
	id [KindPrimaryKey](string)#主键
	name [KindNormal](string)#名称
	update_at [KindOptional](datetime)#更新日期
	author [KindArray](author)#作者
}]`, fmt.Sprintf("%v", packages))
}

func TestPackageMultiTable(t *testing.T) {
	parser := NewParser()
	parser.AddYaml([]byte(
`# 性别
gender:
  - unset # 未设置
  - male # 男
  - female # 女
---
another: &table
  id!: string # 主键
# 示例
demo: &table
  id!: string # 主键
  name: string # 名称
  update_at?: datetime # 更新日期
  author[]: # 作者
    id?: string # NAID
    name: string # 姓名
    gender: gender # 作者性别
`))
	_, err := parser.Parse()
	assert.Equal(t, []error{fmt.Errorf("yaml块中包含多张表")}, err)
}

func TestPackageMultiPackage(t *testing.T) {
	parser := NewParser()
	parser.AddYaml([]byte(
`# 性别
gender:
  - unset # 未设置
  - male # 男
  - female # 女
---
demo: &table
  id!: string # 主键
---
# 示例
demo: &table
  id!: string # 主键
  name: string # 名称
  update_at?: datetime # 更新日期
  author[]: # 作者
    id?: string # NAID
    name: string # 姓名
    gender: gender # 作者性别
`))
	_, err := parser.Parse()
	assert.Equal(t, []error{fmt.Errorf("重复的包: demo")}, err)
}

func TestPackageEnumConflict(t *testing.T) {
	parser := NewParser()
	parser.AddYaml([]byte(
`# 性别
gender:
  - unset # 未设置
  - male # 男
  - female # 女
---
# 性别
gender:
  - unset # 未设置
  - male # 男
  - female # 女
---
# 示例
demo: &table
  id!: string # 主键
  name: string # 名称
  update_at?: datetime # 更新日期
  author[]: # 作者
    id?: string # NAID
    name: string # 姓名
    gender: gender # 作者性别
`))
	_, err := parser.Parse()
	assert.Equal(t, []error{fmt.Errorf("重复的枚举类型: gender")}, err)
}

func TestPackageStructureConflict(t *testing.T) {
	parser := NewParser()
	parser.AddYaml([]byte(
`# 性别
gender:
  - unset # 未设置
  - male # 男
  - female # 女
# 配置
conf:
  id!: string # 主键
---
# 配置
conf:
  id!: string # 主键
---
# 示例
demo: &table
  id!: string # 主键
  name: string # 名称
  update_at?: datetime # 更新日期
  author[]: # 作者
    id?: string # NAID
    name: string # 姓名
    gender: gender # 作者性别
`))
	_, err := parser.Parse()
	assert.Equal(t, []error{fmt.Errorf("重复的结构: conf")}, err)
}

func TestParseWrongYaml(t *testing.T) {
	parser := NewParser()
	parser.AddYaml([]byte(
`&`))
	_, err := parser.Parse()
	assert.NotNil(t, err)
}

func TestParseEarlyReturn(t *testing.T) {
	{
		parser := NewParser()
		parser.AddYamlPath("NOT_REAL_PATH")
		_, err := parser.Parse()
		assert.NotNil(t, err)
	}
	{
		parser := NewParser()
		parser.AddYamlPath("https://NOT_REAL_PATH")
		_, err := parser.Parse()
		assert.NotNil(t, err)
	}
}

func TestAddYamlPath(t *testing.T) {
	{
		parser := NewParser()
		parser.AddYamlPath("https://dev.aminer.cn/gopkg/luna/-/raw/master/schema.yaml")
		_, err := parser.Parse()
		assert.Nil(t, err)
	}
	{
		parser := NewParser()
		parser.AddYamlPath("../../document/demo.yml")
		_, err := parser.Parse()
		assert.Nil(t, err)
	}
}