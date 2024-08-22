package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gagliardetto/solana-go"
)

type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return "string representation of flag"
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func BoolPointer(b bool) *bool {
	return &b
}

func Uint64Ptr(val uint64) *uint64 {
	return &val
}

func ReplaceLastComma(str string, replacement string) string {
	lastCommaIndex := strings.LastIndex(str, ",")
	if lastCommaIndex != -1 {
		str = str[:lastCommaIndex] + replacement + str[lastCommaIndex+1:]
	}

	return str
}

func UnpackStruct(s interface{}) []interface{} {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference pointer
	}

	var result []interface{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.Ptr {
			field = field.Elem() // Dereference pointer
		}

		fieldType := field.Type()

		if fieldType == reflect.TypeOf(solana.PublicKey{}) {
			result = append(result, field.Interface().(solana.PublicKey).String())
		} else {
			result = append(result, field.Interface())
		}

	}

	return result
}

func BuildInsertQuery(i any) string {
	column := "("
	values := " VALUES ("
	typ := reflect.TypeOf(i).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")

		column += fmt.Sprintf("%s,", jsonTag)
		values += "?,"

	}

	column = ReplaceLastComma(column, ")")
	values = ReplaceLastComma(values, ")")

	return column + values
}
