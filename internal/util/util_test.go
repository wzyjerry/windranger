package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnake(t *testing.T) {
	validator := require.New(t)
	validator.Equal("username", Snake("Username"))
	validator.Equal("full_name", Snake("FullName"))
	validator.Equal("http_code", Snake("HTTPCode"))
}

func TestCamel(t *testing.T) {
	validator := require.New(t)
	validator.Equal("userInfo", Camel("user_info"))
	validator.Equal("fullName", Camel("full_name"))
	validator.Equal("userID", Camel("user_id"))
	validator.Equal("fullAdmin", Camel("full-admin"))
	validator.Equal("admin", Camel("admin"))
}

func TestPascal(t *testing.T) {
	validator := require.New(t)
	validator.Equal("UserInfo", Pascal("user_info"))
	validator.Equal("FullName", Pascal("full_name"))
	validator.Equal("UserID", Pascal("user_id"))
	validator.Equal("FullAdmin", Pascal("full-admin"))
}

func TestProtoPascal(t *testing.T) {
	validator := require.New(t)
	validator.Equal("UserInfo", ProtoPascal("user_info"))
	validator.Equal("FullName", ProtoPascal("full_name"))
	validator.Equal("UserId", ProtoPascal("user_id"))
	validator.Equal("FullAdmin", ProtoPascal("full-admin"))
}

func TestAdd(t *testing.T) {
	validator := require.New(t)
	x := make([]int, 100)
	for i := range x {
		x[i] = i + 1
	}
	validator.Equal(55, Add(x[:10]...))
	validator.Equal(5050, Add(x...))
}

func TestPlural(t *testing.T) {
	validator := require.New(t)
	validator.Equal("details", Plural("detail"))
	validator.Equal("traces", Plural("trace"))
	validator.Equal("urls", Plural("url"))
	validator.Equal("fishSlice", Plural("fish"))
}

func TestGetPackageName(t *testing.T) {
	validator := require.New(t)
	validator.Equal("common", GetPackageName("common"))
	validator.Equal("publicationNested", GetPackageName("publication"))
	validator.Equal("personNested", GetPackageName("person"))
}
