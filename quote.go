package quicksql

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var hex = []byte(`0123456789abcdef`)

func QuoteString(s string) string {
	var res = make([]rune, 0, 2+2*len(s))

	res = append(res, '"')

	for _, ch := range s {
		switch ch {
		case 0:
			res = append(res, '\\', '0')
		case '\n':
			res = append(res, '\\', 'n')
		case '\r':
			res = append(res, '\\', 'r')
		case '\\':
			res = append(res, '\\', '\\')
		case '\'':
			res = append(res, '\\', '\'')
		case '"':
			res = append(res, '\\', '"')
		case '\x1a':
			res = append(res, '\\', 'Z')
		default:
			res = append(res, ch)
		}
	}

	res = append(res, '"')

	return string(res)
}

func QuoteBytes(s []byte) string {
	if len(s) == 0 {
		return `""`
	}

	var res = make([]byte, 0, 2+2*len(s))

	res = append(res, '0', 'x')

	for _, ch := range s {
		res = append(res, hex[ch>>4], hex[ch&0xf])
	}

	return string(res)
}

func QuoteIdentifier(s string) (string, error) {
	var res = make([]rune, 0, 2+len(s))

	res = append(res, '`')

	for _, ch := range s {
		switch ch {
		case '\x00':
			return "", ErrInvalidIdentifier

		case '\n':
			return "", ErrInvalidIdentifier

		case '`':
			res = append(res, '`', '`')

		default:
			res = append(res, ch)
		}
	}

	res = append(res, '`')

	return string(res), nil
}

func Quote(data any) string {
	var (
		typ = reflect.TypeOf(data)
		val = reflect.ValueOf(data)
	)

	if data == nil {
		return "NULL"
	}

	switch typ.Kind() {
	case reflect.Slice:
		switch typ.Elem().Kind() {
		case reflect.Uint8:
			return QuoteBytes(val.Bytes())

		case reflect.Uint32:
			if runes, ok := val.Interface().([]rune); ok {
				return QuoteString(string(runes))
			}
		}

	case reflect.String:
		return QuoteString(val.String())

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return fmt.Sprintf("%d", data)

	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", data)
	}

	return QuoteString(fmt.Sprint(data))
}

func QuoteMultipleStruct(data any) string {
	var (
		typ = reflect.TypeOf(data)
		val = reflect.ValueOf(data)
		num = typ.NumField()
		ret = make([]string, 0, num)
	)

	for i := 0; i < num; i++ {
		var (
			fieldTyp = typ.Field(i)
			fieldVal = val.Field(i)
		)

		if !fieldTyp.IsExported() {
			continue
		}

		if quicksql, ok := fieldTyp.Tag.Lookup("quicksql"); ok && quicksql == "-" {
			continue
		}

		ret = append(ret, Quote(fieldVal.Interface()))
	}

	return strings.Join(ret, ", ")
}

func QuoteMultiple(cols ...any) string {
	if len(cols) == 1 {
		typ := reflect.TypeOf(cols[0])
		val := reflect.ValueOf(cols[0])

		for typ.Kind() == reflect.Pointer {
			if val.IsNil() {
				return ""
			}

			typ = typ.Elem()
			val = val.Elem()
		}

		if typ.Kind() == reflect.Struct {
			return QuoteMultipleStruct(val.Interface())
		}
	}

	var args = make([]string, len(cols))

	for i, col := range cols {
		args[i] = Quote(col)
	}

	return strings.Join(args, ", ")
}

var (
	reCanonical  = regexp.MustCompile(`([A-Z][a-z]+)`)
	reMultiUnder = regexp.MustCompile(`[-_]+`)
)

func CanonicalName(s string) string {
	s = reCanonical.ReplaceAllString(s, "_$1-")
	s = strings.ToLower(s)
	s = reMultiUnder.ReplaceAllLiteralString(s, "_")
	return strings.Trim(s, "_")
}

func QuoteColumnNamesStruct(data any) (string, error) {
	var (
		typ = reflect.TypeOf(data)
		num = typ.NumField()
		ret = make([]string, 0, num)
	)

	for i := 0; i < num; i++ {
		var field = typ.Field(i)

		if !field.IsExported() {
			continue
		}

		if quicksql, ok := field.Tag.Lookup("quicksql"); ok {
			if quicksql == "-" {
				continue
			}

			id, err := QuoteIdentifier(quicksql)
			if err != nil {
				return "", err
			}

			ret = append(ret, id)
		} else {
			id, err := QuoteIdentifier(CanonicalName(field.Name))
			if err != nil {
				return "", err
			}

			ret = append(ret, id)
		}
	}

	return strings.Join(ret, ", "), nil
}

func QuoteColumnNames(cols ...any) (string, error) {
	if len(cols) == 1 {
		typ := reflect.TypeOf(cols[0])
		val := reflect.ValueOf(cols[0])

		for typ.Kind() == reflect.Pointer {
			if val.IsNil() {
				return "", ErrNil
			}

			typ = typ.Elem()
			val = val.Elem()
		}

		if typ.Kind() == reflect.Struct {
			return QuoteColumnNamesStruct(val.Interface())
		}
	}

	var (
		ret = make([]string, len(cols))
		err error
	)

	for i, col := range cols {
		ret[i], err = QuoteIdentifier(fmt.Sprint(col))
		if err != nil {
			return "", err
		}
	}

	return strings.Join(ret, ", "), nil
}
