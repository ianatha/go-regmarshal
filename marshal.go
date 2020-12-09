package regmarshal

import (
	"reflect"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

// Marshal marshals.
func Marshal(v interface{}, key registry.Key, path string) error {
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

	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)
		field := rv.Field(i)

		err := marshalField(regkey, typeField, field)
		if err != nil {
			return err
		}
	}

	return nil
}

func marshalField(regkey registry.Key, typeField reflect.StructField, field reflect.Value) (err error) {
	switch field.Kind() {
	case reflect.String:
		err = regkey.SetStringValue(registryPath(typeField), field.String())
	case reflect.Int:
		err = regkey.SetQWordValue(registryPath(typeField), uint64(field.Int()))
	case reflect.Slice:
		err = regkey.SetBinaryValue(registryPath(typeField), field.Bytes())
	default:
		return errors.Errorf("unexpected type: %s", field.Kind().String())
	}

	return
}
