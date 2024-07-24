package helpers

import "reflect"

// GetMapKeys gets the keys of a map as a slice of strings
// Note that the order of the keys is not guaranteed
// May need to sort the keys if you need them in a specific order
func GetMapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// CloneMapShallow creates a (shallow) copy of the given map
func CloneMapShallow(input interface{}) interface{} {
	inputVal := reflect.ValueOf(input)
	if inputVal.Kind() != reflect.Map {
		return nil
	}

	outputVal := reflect.MakeMap(inputVal.Type())

	for _, key := range inputVal.MapKeys() {
		outputVal.SetMapIndex(key, inputVal.MapIndex(key))
	}

	return outputVal.Interface()
}

func deepCopyValue(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Map:
		return reflect.ValueOf(CloneMapDeep(value.Interface()))
	case reflect.Slice:
		sliceCopy := reflect.MakeSlice(value.Type(), value.Len(), value.Cap())
		for i := 0; i < value.Len(); i++ {
			sliceCopy.Index(i).Set(deepCopyValue(value.Index(i)))
		}
		return sliceCopy
	case reflect.Ptr:
		ptrCopy := reflect.New(value.Elem().Type())
		ptrCopy.Elem().Set(deepCopyValue(value.Elem()))
		return ptrCopy
	default:
		return value
	}
}

// CloneMapDeep creates a deep copy of the given map
func CloneMapDeep(input interface{}) interface{} {
	inputVal := reflect.ValueOf(input)
	if inputVal.Kind() != reflect.Map {
		return nil
	}

	outputVal := reflect.MakeMap(inputVal.Type())

	for _, key := range inputVal.MapKeys() {
		val := inputVal.MapIndex(key)
		outputVal.SetMapIndex(key, deepCopyValue(val))
	}

	return outputVal.Interface()
}
