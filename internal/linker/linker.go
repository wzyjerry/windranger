package linker

import (
	"sort"
	"strings"

	"github.com/wzyjerry/windranger/internal/parser"
)

type FieldFunc func(string) string

type linker struct {
	typemap   map[string]*parser.Type
	packages  []*parser.Package
	fieldFunc FieldFunc
	errors    []error
}

func NewLinker() *linker {
	return &linker{
		typemap: make(map[string]*parser.Type),
		errors:  make([]error, 0),
	}
}

func (l *linker) AddPackages(packages []*parser.Package) *linker {
	l.packages = packages
	return l
}

func (l *linker) AddTypemap(src string, dest string, pack string) *linker {
	l.typemap[src] = &parser.Type{
		Name:    dest,
		Package: pack,
	}
	return l
}

func (l *linker) SetFieldFunc(f FieldFunc) *linker {
	l.fieldFunc = f
	return l
}

func (l *linker) Link() ([]*parser.Package, []error) {
	for _, pack := range l.packages {
		for _, enum := range pack.Enums {
			for _, field := range enum.EnumFields {
				field.Name = strings.ToUpper(enum.Name + "_" + field.Name)
			}
		}
		depSet := make(map[string]struct{})
		for _, structure := range pack.Structures {
			for _, field := range structure.Fields {
				raw := field.Type.Raw
				if t, ok := l.typemap[raw]; ok {
					field.Type.Name = t.Name
					field.Type.Package = t.Package
					depSet[t.Package] = struct{}{}
				} else {
					field.Type.Name = l.fieldFunc(raw)
				}
			}
		}
		for dep := range depSet {
			if dep != "" {
				pack.Dependencies = append(pack.Dependencies, dep)
			}
		}
		sort.SliceStable(pack.Dependencies, func(i, j int) bool {
			return pack.Dependencies[i] < pack.Dependencies[j]
		})
	}
	return l.packages, l.errors
}
