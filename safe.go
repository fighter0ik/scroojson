package scroojson

import (
	"encoding/json"
	"reflect"
	"strings"
)

type Safe[T any] struct {
	Content T
	deposit map[string]json.RawMessage
}

func (s *Safe[T]) UnmarshalJSON(bytes []byte) error {
	var content T
	var deposit map[string]json.RawMessage

	if err := json.Unmarshal(bytes, &content); err != nil {
		panic(err) // todo
	}

	t := reflect.TypeOf(content)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		if err := json.Unmarshal(bytes, &deposit); err != nil {
			panic(err) // todo
		}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}

			n := ""
			if tag, ok := f.Tag.Lookup("json"); ok {
				if tag == "-" {
					continue
				}

				for _, s := range strings.Split(tag, ",") {
					if s != "omitempty" {
						n = s
					}
				}
			}
			if len(n) < 1 {
				n = strings.ToLower(f.Name)
			}
			delete(deposit, n)
		}
	}

	s.Content = content
	s.deposit = deposit
	return nil
}

func (s Safe[T]) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(s.Content)
	if err != nil {
		panic(err) // todo
	}

	if len(s.deposit) < 1 {
		return bytes, nil
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &m); err != nil {
		panic(err) // todo
	}
	for k, v := range s.deposit {
		m[k] = v
	}

	bytes, err = json.Marshal(m)
	if err != nil {
		panic(err) // todo
	}
	return bytes, nil
}
