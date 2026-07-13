package middlewares

import "reflect"

// isNilInterface reports whether v is nil or is an interface containing a nil
// pointer/channel/map/slice/function. This is necessary because a Go interface
// holding a typed nil value is not equal to nil.
func isNilInterface(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	}
	return false
}
