package reflect

import "reflect"

// StructName used to get the struct name from the obj.
func StructName(obj any) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
