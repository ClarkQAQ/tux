package rest

import "reflect"

func rangeStructMethod(val any, cb func(v reflect.Value, m reflect.Method) error) error {
	v := reflect.ValueOf(val)
	t := v.Type()

	if v.Kind() != reflect.Pointer && v.Kind() != reflect.Struct {
		return ErrValueNotStructOrPointer
	}

	for i := 0; i < v.NumMethod(); i++ {
		v := v.Method(i)
		m := t.Method(i)

		if !v.IsValid() || !v.CanInterface() {
			continue
		}

		if e := cb(v, m); e != nil {
			return e
		}
	}

	return nil
}
