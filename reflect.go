package path_helpers

import "reflect"

func PkgPathOf(value interface{}) string {
	var typ reflect.Type
	switch vt := value.(type) {
	case reflect.Type:
		typ = vt
	case reflect.Value:
		typ = vt.Type()
	default:
		typ = reflect.TypeOf(value)
	}

TypeLoop:
	for {
		switch typ.Kind() {
		case reflect.Ptr, reflect.Slice:
			typ = typ.Elem()
		default:
			break TypeLoop
		}
	}

	return TrimGoPathC(typ.PkgPath())
}
