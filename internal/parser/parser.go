package parser

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// windranger windranger结构
type windranger struct {
	Version   string   `yaml:"version"`
	Kind      string   `yaml:"kind"`
	Resources []string `yaml:"resources"`
}

// model model结构
type model struct {
	Version  string `yaml:"version"`
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec yaml.Node `yaml:"spec"`
}

// findConflict 查找lst中重复的id
func findConflict[T any](lst []T, idFunc func(item T) string) []string {
	set := make(map[string]struct{})
	conflict := make([]string, 0)
	for i := range lst {
		id := idFunc(lst[i])
		if _, ok := set[id]; ok {
			conflict = append(conflict, id)
		}
		set[id] = struct{}{}
	}
	return conflict
}

type parser struct {
	contents [][]byte
	errors   []error
	// 当前块内信息
	tableName  string
	table      *Structure
	structures []*Structure
	enums      []*Enum
}

func NewParser() *parser {
	return &parser{
		contents: make([][]byte, 0),
		errors:   make([]error, 0),
	}
}

// AddYaml 添加yaml
func (p *parser) AddYaml(yaml []byte) *parser {
	p.contents = append(p.contents, yaml)
	return p
}

// AddYamlPath 添加yaml文件
func (p *parser) AddYamlPath(uri string) *parser {
	// 如果path为url，克隆目录
	root := uri
	u, err := url.Parse(uri)
	if err == nil && strings.HasPrefix(u.Scheme, "http") {
		panic("克隆git")
		// resp, err := http.Get(u.String())
		// if err != nil {
		// 	p.errors = append(p.errors, err)
		// 	return p
		// }
		// defer resp.Body.Close()
		// content, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	p.errors = append(p.errors, err)
		// 	return p
		// }
		// p.contents = append(p.contents, content)
	}
	// 解析配置文件
	content, err := os.ReadFile(path.Join(root, "windranger.yaml"))
	if err != nil {
		p.errors = append(p.errors, err)
		return p
	}
	var cfg windranger
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		p.errors = append(p.errors, err)
		return p
	}
	if cfg.Version != "v1" {
		p.errors = append(p.errors, fmt.Errorf("未知版本号: %s", cfg.Version))
		return p
	}
	if cfg.Kind != "Windranger" {
		p.errors = append(p.errors, fmt.Errorf("未知资源类型: %s", cfg.Kind))
		return p
	}
	for _, sub := range cfg.Resources {
		content, err := os.ReadFile(path.Join(root, sub))
		if err != nil {
			p.errors = append(p.errors, err)
			return p
		}
		p.contents = append(p.contents, content)
	}
	return p
}

// trimComment 格式化注释
func trimComment(comment string) string {
	return strings.Join(strings.Fields(strings.Trim(comment, "#")), " ")
}

// parseComment 格式化首个非空注释
func parseComment(comments ...string) string {
	for _, comment := range comments {
		if tc := trimComment(comment); tc != "" {
			return tc
		}
	}
	return ""
}

// parseSequence 解析枚举类型
func (p *parser) parseSequence(node *yaml.Node) []*EnumField {
	fields := make([]*EnumField, 0, len(node.Content))
	for _, enum := range node.Content {
		if enum.Kind != yaml.ScalarNode {
			p.errors = append(p.errors, fmt.Errorf("枚举类型必须为标量"))
			continue
		}
		fields = append(fields, &EnumField{
			Name:    enum.Value,
			Comment: parseComment(enum.HeadComment, enum.LineComment),
		})
	}
	for _, id := range findConflict(fields, func(field *EnumField) string {
		return field.Name
	}) {
		p.errors = append(p.errors, fmt.Errorf("重复的枚举值: %v", id))
	}
	return fields
}

// parseMapping 解析字典类型
func (p *parser) parseMapping(node *yaml.Node) []*Field {
	fields := make([]*Field, 0, len(node.Content)>>1)
	// 遍历kv-pair
	var key, value *yaml.Node
	for i := 0; i < len(node.Content)>>1; i++ {
		key, value = node.Content[i<<1], node.Content[i<<1|1]
		name := key.Value
		kind := KindNormal
		switch {
		case strings.HasSuffix(name, "[]"):
			kind = KindArray
			name = name[:len(name)-2]
		case strings.HasSuffix(name, "?"):
			kind = KindOptional
			name = name[:len(name)-1]
		case strings.HasSuffix(name, "!"):
			kind = KindPrimaryKey
			name = name[:len(name)-1]
		}
		field := &Field{
			Name: name,
			Type: &Type{
				Kind: kind,
			},
		}
		// 解析值类型
		switch value.Kind {
		case yaml.SequenceNode:
			subFields := p.parseSequence(value)
			if name == "type" {
				name = name + ""
			}
			enum := &Enum{
				Name:       name,
				Comment:    parseComment(key.HeadComment, key.LineComment, value.LineComment),
				EnumFields: subFields,
			}
			p.enums = append(p.enums, enum)
			// 添加字段
			field.Type.Raw = enum.Name
			field.Comment = enum.Comment
			fields = append(fields, field)
		case yaml.MappingNode:
			subFields := p.parseMapping(value)
			structure := &Structure{
				Name:    name,
				Comment: parseComment(key.HeadComment, key.LineComment),
				Fields:  subFields,
			}
			if name == p.tableName {
				p.table = structure
			}
			p.structures = append(p.structures, structure)
			// 添加字段
			field.Type.Raw = name
			field.Comment = parseComment(key.HeadComment, structure.Comment)
			fields = append(fields, field)
		case yaml.ScalarNode:
			// 添加字段
			field.Type.Raw = value.Value
			field.Comment = parseComment(key.HeadComment, value.LineComment)
			fields = append(fields, field)
		}
	}
	for _, id := range findConflict(fields, func(field *Field) string {
		return field.Name
	}) {
		p.errors = append(p.errors, fmt.Errorf("重复的字段名: %v", id))
	}
	return fields
}

// parseDoc 解析yaml块
func (p *parser) parseDoc(node *yaml.Node) *Package {
	// 准备块内缓存
	p.table = nil
	p.enums = make([]*Enum, 0)
	p.structures = make([]*Structure, 0)
	p.parseMapping(node)
	pack := &Package{
		Name:         CommonPackage,
		Enums:        p.enums,
		Structures:   p.structures,
		Dependencies: make([]string, 0),
	}
	if p.table != nil {
		pack.Name = p.table.Name
	}
	return pack
}

// normalize 标准化包结构
func (p *parser) normalize(pack *Package) {
	for _, id := range findConflict(pack.Enums, func(enum *Enum) string {
		return enum.Name
	}) {
		p.errors = append(p.errors, fmt.Errorf("重复的枚举类型: %v", id))
	}
	for _, id := range findConflict(pack.Structures, func(structure *Structure) string {
		return structure.Name
	}) {
		p.errors = append(p.errors, fmt.Errorf("重复的结构: %v", id))
	}
	sort.SliceStable(pack.Enums, func(i, j int) bool {
		return pack.Enums[i].Name < pack.Enums[j].Name
	})
	sort.SliceStable(pack.Structures, func(i, j int) bool {
		return pack.Structures[i].Name < pack.Structures[j].Name
	})
}

// link 链接yaml块，构建输出
func (p *parser) link(packages []*Package) ([]*Package, []error) {
	common := &Package{
		Name:         CommonPackage,
		Enums:        make([]*Enum, 0),
		Structures:   make([]*Structure, 0),
		Dependencies: make([]string, 0),
	}
	linked := make([]*Package, 0)
	var hasCommon bool
	for _, pack := range packages {
		if pack.Name == CommonPackage {
			common.Enums = append(common.Enums, pack.Enums...)
			common.Structures = append(common.Structures, pack.Structures...)
			hasCommon = true
		} else {
			linked = append(linked, pack)
			p.normalize(pack)
		}
	}
	if hasCommon {
		linked = append(linked, common)
		p.normalize(common)
	}
	for _, id := range findConflict(linked, func(pack *Package) string {
		return pack.Name
	}) {
		p.errors = append(p.errors, fmt.Errorf("重复的包: %v", id))
	}
	sort.SliceStable(linked, func(i, j int) bool {
		return linked[i].Name < linked[j].Name
	})
	if len(p.errors) != 0 {
		return nil, p.errors
	}
	return linked, nil
}

// Parse 解析生成Info结构
func (p *parser) Parse() ([]*Package, []error) {
	// 提前返回
	if len(p.errors) != 0 {
		return nil, p.errors
	}
	var (
		err error
	)
	packages := make([]*Package, 0)
	for _, content := range p.contents {
		var cfg model
		err = yaml.Unmarshal(content, &cfg)
		if err != nil {
			p.errors = append(p.errors, err)
			break
		}
		if cfg.Version != "v1" {
			p.errors = append(p.errors, fmt.Errorf("未知版本号: %s", cfg.Version))
			break
		}
		if cfg.Kind != "Model" {
			p.errors = append(p.errors, fmt.Errorf("未知资源类型: %s", cfg.Kind))
			break
		}
		p.tableName = cfg.Metadata.Name
		packages = append(packages, p.parseDoc(&cfg.Spec))
	}
	return p.link(packages)
}
