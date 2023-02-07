## windranger - YAML-based model builder

windranger 是一个简单而功能强大的基于 YAML 的实体框架，易于构建和维护数据结构与大数据模型。

## 快速安装

```console
go install github.com/wzyjerry/windranger
```

---

## 语法指南

### 基本类型

- int
- float
- bool
- string
- datetime
- objectid

### 复合类型

- `[]`，集合，标注在 key 上。例如: `urls[]: string`
- 枚举类型，使用 yaml 数组表示。

### 字段标记

- `!`: 标记主键。例如: `id!: string`
- `?`: 标记可空。例如: `name?: string`

### 额外说明

- 禁止使用 common 作为表名（保留字）。
- 任何不包含&table 的文件都被合并至 common 包，否则独立为表名+Nested 包。
- enum 与 type 的查找策略：首先在同文件内查找，未找到则在 common 中查找。

### 最佳实践

1. 使用单复数区分字段和数组
2. 使用**snake 形式**定义名称
