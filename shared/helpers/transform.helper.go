package helper

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
)

type transform struct{}

func NewTransform() inf.ITransform {
	return transform{}
}

func (h transform) SrcToDest(src, dest any) error {
	helper := NewParser()

	srcByte, err := helper.Marshal(src)
	if err != nil {
		return err
	}

	if err = helper.Unmarshal(srcByte, dest); err != nil {
		return err
	}

	return nil
}

func (h transform) ReqToRes(src, dest any) error {
	helper := NewParser()

	srcByte, err := helper.Marshal(src)
	if err != nil {
		return err
	}

	if err = helper.Unmarshal(srcByte, dest); err != nil {
		return err
	}

	return nil
}

func (h transform) ResToReq(src, dest any) error {
	helper := NewParser()

	srcByte, err := helper.Marshal(src)
	if err != nil {
		return err
	}

	if err = helper.Unmarshal(srcByte, dest); err != nil {
		return err
	}

	return nil
}

func (h transform) QueryToStruct(query string, dest any) error {
	valueof := reflect.ValueOf(dest)
	if valueof.Kind() != reflect.Pointer {
		return errors.New("dest must be a pointer to a struct")
	}

	if valueof.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to a struct")
	}

	parsed, err := url.ParseQuery(query)
	if err != nil {
		return err
	}

	store := make(map[string]string)
	for key, values := range parsed {
		if len(values) > 0 {
			unescaped, err := url.QueryUnescape(values[0])
			if err != nil {
				return fmt.Errorf("failed to unescape value for key %s: %v", key, err)
			}
			store[key] = unescaped
		}
	}

	structType := valueof.Elem().Type()
	structValue := valueof.Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		stag := field.Tag.Get("query")

		if stag == "" {
			return fmt.Errorf("field %s must have a query tag", field.Name)
		}

		val, ok := store[stag]
		if ok {
			fieldValue := structValue.Field(i)

			switch field.Type.Kind() {
			case reflect.String:
				fieldValue.SetString(val)

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				toInt, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("field %s: invalid integer value", stag)
				}

				fieldValue.SetInt(int64(toInt))

			case reflect.Bool:
				toBool, err := strconv.ParseBool(val)
				if err != nil {
					return fmt.Errorf("field %s: invalid boolean value", stag)
				}

				fieldValue.SetBool(toBool)

			case reflect.Float32, reflect.Float64:
				toFloat, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return fmt.Errorf("field %s: invalid float value", stag)
				}

				fieldValue.SetFloat(toFloat)

			case reflect.Map:
				if field.Type != reflect.TypeOf(map[string]any{}) {
					return fmt.Errorf("field %s: only map[string]any is supported", stag)
				}

				var m map[string]any
				if err := json.Unmarshal([]byte(val), &m); err != nil {
					return fmt.Errorf("field %s: invalid JSON", stag)
				}

				fieldValue.Set(reflect.ValueOf(m))

			default:
				return fmt.Errorf("unsupported type for field %s: %s", stag, field.Type.Kind())
			}
		}
	}

	return nil
}
