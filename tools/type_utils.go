package tools

import (
	"reflect"
)

func StructToMap(obj interface{}, nested bool) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	return StructReflectToMap(t, v, nested)
}

func StructReflectToMap(t reflect.Type, v reflect.Value, nested bool) map[string]interface{} {
	if t.Kind() == reflect.Struct {
		res := make(map[string]interface{})
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			vv := v.Field(i)
			m := StructReflectToMap(f.Type, vv, nested)
			if m != nil {
				if nested {
					res[f.Name] = m
				} else {
					for mk, mv := range m {
						res[mk] = mv
					}
				}
			} else {
				res[f.Name] = vv.Interface()
			}
		}
		return res
	}
	return nil
}
