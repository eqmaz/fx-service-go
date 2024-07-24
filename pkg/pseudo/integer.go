package pseudo

import (
	"fmt"
	"reflect"
)

// Integer pseudo-type representing an integer of any size. Can be used to store any integer type.
type Integer struct {
	value interface{}
}

// From sets the value of the Integer pseudo-type. Accepts any integer type. Int and Uint types are converted to int64 and uint64 internally.
func (a *Integer) From(v interface{}) error {
	switch v := v.(type) {
	case int, int8, int16, int32, int64:
		a.value = reflect.ValueOf(v).Int()
	case uint, uint8, uint16, uint32, uint64:
		a.value = reflect.ValueOf(v).Uint()
	case float32, float64:
		a.value = int64(reflect.ValueOf(v).Float())
	default:
		//a.value = int64(0)
		return fmt.Errorf("unsupported type: %T. Must pass an integer type", v)
	}
	return nil
}

func (a *Integer) GetType() string {
	switch a.value.(type) {
	case uint64:
		return "uint64"
	default:
		return "int64"
	}
}

func (a *Integer) Value() interface{} {
	return a.value
}

func (a *Integer) Int() int {
	switch v := a.value.(type) {
	case int64:
		return int(v)
	case uint64:
		return int(v)
	default:
		return 0
	}
}

func (a *Integer) Int8() int8 {
	switch v := a.value.(type) {
	case int64:
		return int8(v)
	case uint64:
		return int8(v)
	default:
		return 0
	}
}

func (a *Integer) Int16() int16 {
	switch v := a.value.(type) {
	case int64:
		return int16(v)
	case uint64:
		return int16(v)
	default:
		return 0
	}
}

func (a *Integer) Int32() int32 {
	switch v := a.value.(type) {
	case int64:
		return int32(v)
	case uint64:
		return int32(v)
	default:
		return 0
	}
}

func (a *Integer) Int64() int64 {
	switch v := a.value.(type) {
	case int64:
		return v
	case uint64:
		return int64(v)
	default:
		return 0
	}
}

func (a *Integer) Uint() uint {
	switch v := a.value.(type) {
	case int64:
		return uint(v)
	case uint64:
		return uint(v)
	default:
		return 0
	}
}

func (a *Integer) Uint8() uint8 {
	switch v := a.value.(type) {
	case int64:
		return uint8(v)
	case uint64:
		return uint8(v)
	default:
		return 0
	}
}

func (a *Integer) Uint16() uint16 {
	switch v := a.value.(type) {
	case int64:
		return uint16(v)
	case uint64:
		return uint16(v)
	default:
		return 0
	}
}

func (a *Integer) Uint32() uint32 {
	switch v := a.value.(type) {
	case int64:
		return uint32(v)
	case uint64:
		return uint32(v)
	default:
		return 0
	}
}

func (a *Integer) Uint64() uint64 {
	switch v := a.value.(type) {
	case int64:
		return uint64(v)
	case uint64:
		return v
	default:
		return 0
	}
}

func (a *Integer) Increment() {
	switch v := a.value.(type) {
	case int64:
		a.value = v + 1
	case uint64:
		a.value = v + 1
	}
}

func (a *Integer) Decrement() {
	switch v := a.value.(type) {
	case int64:
		a.value = v - 1
	case uint64:
		a.value = v - 1
	}
}

func (a *Integer) Add(x interface{}) error {
	switch x := x.(type) {
	case int, int8, int16, int32, int64:
		val := reflect.ValueOf(x).Int()
		switch v := a.value.(type) {
		case int64:
			a.value = v + val
		case uint64:
			a.value = v + uint64(val)
		}
	case uint, uint8, uint16, uint32, uint64:
		val := reflect.ValueOf(x).Uint()
		switch v := a.value.(type) {
		case int64:
			a.value = v + int64(val)
		case uint64:
			a.value = v + val
		}
	default:
		return fmt.Errorf("unsupported type: %T. Must pass an integer type", x)
	}
	return nil
}

func (a *Integer) Subtract(x interface{}) error {
	switch x := x.(type) {
	case int, int8, int16, int32, int64:
		val := reflect.ValueOf(x).Int()
		switch v := a.value.(type) {
		case int64:
			a.value = v - val
		case uint64:
			a.value = v - uint64(val)
		}
	case uint, uint8, uint16, uint32, uint64:
		val := reflect.ValueOf(x).Uint()
		switch v := a.value.(type) {
		case int64:
			a.value = v - int64(val)
		case uint64:
			a.value = v - val
		}
	default:
		return fmt.Errorf("unsupported type: %T. Must pass an integer type", x)
	}
	return nil
}
