package config

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

func marshalConfig(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("cfg: Marshal requires a struct, got %s", rv.Kind())
	}

	rt := rv.Type()
	var sb strings.Builder

	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("cfg")
		if tag == "" || tag == "-" {
			continue
		}

		var strVal string
		switch val := rv.Field(i).Interface().(type) {
		case time.Time:
			strVal = val.Format(time.DateTime)
		case string:
			strVal = val
		default:
			return nil, fmt.Errorf("cfg: unsupported field type %T for %q", val, tag)
		}

		fmt.Fprintf(&sb, "%s='%s'\n", tag, escapeQuote(strVal))
	}

	return []byte(sb.String()), nil
}

func unmarshalConfig(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cfg: Unmarshal requires a pointer to struct")
	}
	rv = rv.Elem()
	rt := rv.Type()

	tagToField := make(map[string]int)
	for i := 0; i < rt.NumField(); i++ {
		if tag := rt.Field(i).Tag.Get("cfg"); tag != "" && tag != "-" {
			tagToField[tag] = i
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, rawVal, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("cfg: malformed line %q", line)
		}

		idx, found := tagToField[key]
		if !found {
			continue // unknown key — ignore for forward compatibility
		}

		val := unquote(rawVal)
		field := rv.Field(idx)

		switch field.Interface().(type) {
		case time.Time:
			t, err := time.Parse(time.DateTime, val)
			if err != nil {
				return fmt.Errorf("cfg: parsing %q as time: %w", key, err)
			}
			field.Set(reflect.ValueOf(t))
		case string:
			field.SetString(val)
		default:
			return fmt.Errorf("cfg: unsupported field type for %q", key)
		}
	}

	return scanner.Err()
}

func writeFile(path string, v interface{}) error {
	data, err := marshalConfig(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func readFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return unmarshalConfig(data, v)
}

func escapeQuote(s string) string {
	return strings.ReplaceAll(s, "'", `'\''`)
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		s = s[1 : len(s)-1]
	}
	return strings.ReplaceAll(s, `'\''`, "'")
}
