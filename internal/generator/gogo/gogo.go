package gogo

import (
	"bytes"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/wzyjerry/windranger/internal/linker"
	"github.com/wzyjerry/windranger/internal/parser"
	tmpl "github.com/wzyjerry/windranger/internal/template"
	"github.com/wzyjerry/windranger/internal/util"
)

// InfoGogo Go模板信息
type InfoGogo struct {
	// PackageName go文件包名
	PackageName string
	// Imports 导入的go文件
	Imports []string

	// Enums 枚举类型
	Enums []*parser.Enum
	// Structures 结构
	Structures []*parser.Structure
}

func Generate(packages []*parser.Package, out string) error {
	l := linker.NewLinker().AddPackages(packages).SetFieldFunc(util.ProtoPascal)
	l.
		AddTypemap("int", "int64", "").
		AddTypemap("float", "double", "").
		AddTypemap("bool", "bool", "").
		AddTypemap("string", "string", "").
		AddTypemap("datetime", "Time", "time").
		AddTypemap("objectid", "ObjectID", "primitive")
	packages, errs := l.Link()
	if len(errs) != 0 {
		return errs[0]
	}
	for _, pack := range packages {
		imports := make([]string, len(pack.Dependencies))
		for i, dep := range pack.Dependencies {
			switch dep {
			case "time":
				imports[i] = "time"
			case "primitive":
				imports[i] = "go.mongodb.org/mongo-driver/bson/primitive"
			default:
				panic("跨文件引用")
			}
		}
		sort.SliceStable(imports, func(i, j int) bool {
			return imports[i] < imports[j]
		})
		// 准备生成信息
		_, folder := path.Split(out)
		info := &InfoGogo{
			PackageName: util.Camel(folder),
			Imports:     imports,
			Enums:       pack.Enums,
			Structures:  pack.Structures,
		}
		// 准备生成目录
		err := os.MkdirAll(out, os.ModePerm)
		if err != nil {
			return err
		}
		// 准备模板
		name := "gogo.tmpl"
		t, err := template.New("gogo").Funcs(util.FuncMap).ParseFS(tmpl.FS, path.Join("gogo", name))
		if err != nil {
			return err
		}
		// 生成
		buffer := bytes.NewBuffer(nil)
		err = t.ExecuteTemplate(buffer, name, info)
		if err != nil {
			return err
		}
		// 写文件
		err = os.WriteFile(path.Join(out, util.Camel(pack.Name)+".go"), buffer.Bytes(), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
