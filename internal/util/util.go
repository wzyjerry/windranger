package util

import (
	"strings"
	"text/template"
	"unicode"

	"github.com/go-openapi/inflect"
	"github.com/wzyjerry/windranger/internal/parser"
)

var (
	rules    = ruleset()
	acronyms = make(map[string]struct{})
	FuncMap  = template.FuncMap{
		"snake":          Snake,
		"camel":          Camel,
		"pascal":         Pascal,
		"protoPascal":    ProtoPascal,
		"upper":          strings.ToUpper,
		"lower":          strings.ToLower,
		"plural":         Plural,
		"add":            Add,
		"getPackageName": GetPackageName,
		"goType":         GoType,
	}
)

func Add(xs ...int) int {
	result := 0
	for _, x := range xs {
		result += x
	}
	return result
}

func ruleset() *inflect.Ruleset {
	rules := inflect.NewDefaultRuleset()
	// Add common initialisms from golint and more.
	for _, w := range []string{
		"ACL", "API", "ASCII", "AWS", "CPU", "CSS", "DNS", "EOF", "GB", "GUID",
		"HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "KB", "LHS", "MAC", "MB",
		"QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "SSO", "TCP",
		"TLS", "TTL", "UDP", "UI", "UID", "URI", "URL", "UTF8", "UUID", "VM",
		"XML", "XMPP", "XSRF", "XSS",
	} {
		acronyms[w] = struct{}{}
		rules.AddAcronym(w)
	}
	return rules
}

// Snake converts the given struct or field name into a snake_case.
//
//	Username => username
//	FullName => full_name
//	HTTPCode => http_code
func Snake(s string) string {
	var (
		j int
		b strings.Builder
	)
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		// Put '_' if it is not a start or end of a word, current letter is uppercase,
		// and previous is lowercase (cases like: "UserInfo"), or next letter is also
		// a lowercase and previous letter is not "_".
		if i > 0 && i < len(s)-1 && unicode.IsUpper(r) {
			if unicode.IsLower(rune(s[i-1])) ||
				j != i-1 && unicode.IsLower(rune(s[i+1])) && unicode.IsLetter(rune(s[i-1])) {
				j = i
				b.WriteString("_")
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// Camel converts the given name into a camelCase.
//
//	user_info  => userInfo
//	full_name  => fullName
//	user_id    => userID
//	full-admin => fullAdmin
func Camel(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	if len(words) == 1 {
		return strings.ToLower(words[0])
	}
	return strings.ToLower(words[0]) + pascalWords(words[1:])
}

// Pascal converts the given name into a PascalCase.
//
//	user_info 	=> UserInfo
//	full_name 	=> FullName
//	user_id   	=> UserID
//	full-admin	=> FullAdmin
func Pascal(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	return pascalWords(words)
}

func pascalWords(words []string) string {
	for i, w := range words {
		upper := strings.ToUpper(w)
		if _, ok := acronyms[upper]; ok {
			words[i] = upper
		} else {
			words[i] = rules.Capitalize(w)
		}
	}
	return strings.Join(words, "")
}

func isSeparator(r rune) bool {
	return r == '_' || r == '-'
}

// ProtoPascal converts the given name into a proto PascalCase, ignore any acronyms.
//
//	user_info 	=> UserInfo
//	full_name 	=> FullName
//	user_id   	=> UserId
//	full-admin	=> FullAdmin
func ProtoPascal(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	return protoPascalWords(words)
}

func protoPascalWords(words []string) string {
	for i, w := range words {
		words[i] = rules.Capitalize(w)
	}
	return strings.Join(words, "")
}

// plural a name.
func Plural(name string) string {
	p := rules.Pluralize(name)
	if p == name {
		p += "Slice"
	}
	return p
}

// GetPackageName 获取包名
func GetPackageName(name string) string {
	return Camel(name)
}

func withStar(full string, kind parser.Kind) bool {
	if kind == parser.KindOptional {
		return true
	}
	switch full {
	case "string", "int64", "float64", "bool", "time.Time", "primitive.ObjectID":
		return false
	}
	return kind == parser.KindArray
}

func GoType(in parser.Type) string {
	full := in.Package
	if full != "" {
		full += "."
	}
	full += in.Name
	var result string
	if in.Kind == parser.KindArray {
		result += "[]"
	}
	if withStar(full, in.Kind) {
		result += "*"
	}
	result += full
	return result
}
