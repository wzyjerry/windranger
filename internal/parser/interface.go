package parser

import "strings"

const CommonPackage = "common"

type Kind uint32
const (
    KindNormal Kind = iota
    KindArray
    KindOptional
    KindPrimaryKey
)
var kindName = [...]string {
    KindNormal: "KindNormal",
    KindArray: "KindArray",
    KindOptional: "KindOptional",
    KindPrimaryKey: "KindPrimaryKey",    
}

type Type struct {
    Raw string
    Name string
    Kind Kind
    Package string
}

func (t *Type) String() string {
    var builder strings.Builder
    builder.WriteByte('[')
    builder.WriteString(kindName[t.Kind])
    builder.WriteByte(']')
    if t.Package != "" {
        builder.WriteString(t.Package)
        builder.WriteByte('.')
    }
    builder.WriteString(t.Name)
    builder.WriteByte('(')
    builder.WriteString(t.Raw)
    builder.WriteByte(')')
    return builder.String()
}

type Field struct {
    Name string
    Comment string
    Type *Type
}

func (f *Field) String() string {
    var builder strings.Builder
    builder.WriteString(f.Name)
    builder.WriteByte(' ')
    builder.WriteString(f.Type.String())
    builder.WriteString("#")
    builder.WriteString(f.Comment)
    return builder.String()
}

type Structure struct {
	Name string
	Comment string
	Fields []*Field
}

func (s *Structure) String() string {
    var builder strings.Builder
    builder.WriteString("#")
    builder.WriteString(s.Comment)
    builder.WriteByte('\n')
    builder.WriteString("type ")
    builder.WriteString(s.Name)
    builder.WriteString(" struct {\n")
    for _, field := range s.Fields {
        builder.WriteString("\t")
        builder.WriteString(field.String())
        builder.WriteByte('\n')
    }
    builder.WriteByte('}')
    return builder.String()
}

type EnumField struct {
    Name string
    Comment string
}

func (f *EnumField) String() string {
    var builder strings.Builder
    builder.WriteString(f.Name)
    builder.WriteString("#")
    builder.WriteString(f.Comment)
    return builder.String()
}

type Enum struct {
    Name string
    Comment string
    EnumFields []*EnumField
}

func (e *Enum) String() string {
    var builder strings.Builder
    builder.WriteString("#")
    builder.WriteString(e.Comment)
    builder.WriteByte('\n')
    builder.WriteString("type ")
    builder.WriteString(e.Name)
    builder.WriteString(" enum {\n")
    for _, field := range e.EnumFields {
        builder.WriteString("\t")
        builder.WriteString(field.String())
        builder.WriteByte('\n')
    }
    builder.WriteByte('}')
    return builder.String()
}

type Package struct {
    Name string
    Enums []*Enum
    Structures []*Structure
    Dependencies []string
}

func (p *Package) String() string {
    var builder strings.Builder
    builder.WriteString("package ")
    builder.WriteString(p.Name)
    builder.WriteByte('\n')
    for _, dep := range p.Dependencies {
        builder.WriteString(dep)
        builder.WriteByte('\n')
    }
    for _, enum := range p.Enums {
        builder.WriteString(enum.String())
        builder.WriteByte('\n')
    }
    for _, structure := range p.Structures {
        builder.WriteString(structure.String())
        builder.WriteByte('\n')
    }
    return strings.TrimSpace(builder.String())
}
