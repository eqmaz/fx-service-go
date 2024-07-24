package pseudo

import (
	"fmt"
	"reflect"
)

type Number struct {
	value interface{}
}

func (n *Number) From(v interface{}) error {
	switch v := v.(type) {
	case int, int8, int16, int32, int64:
		n.value = reflect.ValueOf(v).Int()
	case uint, uint8, uint16, uint32, uint64:
		n.value = reflect.ValueOf(v).Uint()
	case float32, float64:
		n.value = reflect.ValueOf(v).Float()
	default:
		return fmt.Errorf("unsupported type: %T. Must pass an integer or float type", v)
	}
	return nil
}

func (n *Number) GetType() string {
	switch n.value.(type) {
	case int64:
		return "int64"
	case uint64:
		return "uint64"
	case float64:
		return "float64"
	default:
		return "unknown"
	}
}

func (n *Number) Value() interface{} {
	return n.value
}

func (n *Number) Float64() float64 {
	switch v := n.value.(type) {
	case int64:
		return float64(v)
	case uint64:
		return float64(v)
	case float64:
		return v
	default:
		return 0
	}
}

func (n *Number) Int64() int64 {
	switch v := n.value.(type) {
	case int64:
		return v
	case uint64:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}

func (n *Number) Uint64() uint64 {
	switch v := n.value.(type) {
	case int64:
		return uint64(v)
	case uint64:
		return v
	case float64:
		return uint64(v)
	default:
		return 0
	}
}
