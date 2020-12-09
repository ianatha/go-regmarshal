package regmarshal

import (
	"reflect"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

// Unmarshal unmarshals.
func Unmarshal(key registry.Key, path string, v interface{}) error {
	pointerRv := reflect.ValueOf(v)
	rv := pointerRv.Elem()
	t := rv.Type()

	if pointerRv.Kind() != reflect.Ptr || pointerRv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	regkey, err := registry.OpenKey(key, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer regkey.Close()

	//fmt.Printf("num: %d\n", rv.Elem())
	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)
		field := rv.Field(i)

		err := unmarshalField(regkey, typeField, field)
		if err != nil {
			return errors.Errorf("unmarshaling %s: %v", typeField.Name, err)
		}
	}

	return nil
}

func unmarshalField(regkey registry.Key, typeField reflect.StructField, field reflect.Value) (err error) {
	switch field.Kind() {
	case reflect.String:
		v, _, err := regkey.GetStringValue(registryPath(typeField))
		if err != nil {
			if err != registry.ErrNotExist {
				return err
			}
		} else {
			field.SetString(v)
		}
	case reflect.Int:
		v, _, err := regkey.GetIntegerValue(registryPath(typeField))
		if err != nil {
			if err != registry.ErrNotExist {
				return err
			}
		} else {
			field.SetInt(int64(v))
		}
	case reflect.Slice:
		// TODO check its a slice of byte

		v, _, err := regkey.GetBinaryValue(registryPath(typeField))
		if err != nil {
			if err != registry.ErrNotExist {
				return err
			}
		} else {
			field.SetBytes(v)
		}
	default:
		return errors.Errorf("unexpected type: %s", field.Kind().String())
	}

	return
}
