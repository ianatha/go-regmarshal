// Package regmarshal implements marshalling and unmarshalling Go structs
// into a sensible structure in the Windows Registry.
//
// For the time being, the only supported types are `String`, `int`, and `[]byte{}`.
package regmarshal

import "reflect"

// An InvalidUnmarshalError is a description of an unmarshalling error.
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

func registryPath(typeField reflect.StructField) string {
	// TODO use tag
	//if value, ok := typeField.Tag.Lookup("example"); ok {
	//	fmt.Printf("%s - val: %v - rv.field(d): %v\n", typeField.Name,  value, field)
	//	field.SetString("adsf-123")
	//
	//
	//	//rv.Field(i))
	//	//fmt.Printf("%s pretty please\n", rv.Field(i).String())
	//} else {
	//	fmt.Printf("i\n")
	//}

	return typeField.Name
}
