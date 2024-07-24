package pseudo

import (
	"fmt"
	"strconv"
)

// IntString Union pseudo-type representing an integer or a string
type IntString struct {
	value interface{}
}

func (u *IntString) From(v interface{}) error {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		u.value = v
		return nil
	case string:
		u.value = v
		return nil
	default:
		return fmt.Errorf("unsupported type: %T. Must pass int, uint, or string", v)
	}
}

func (u *IntString) GetType() string {
	switch u.value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "int"
	case string:
		return "string"
	default:
		return "unknown"
	}
}

func (u *IntString) Value() interface{} {
	return u.value
}

func (u *IntString) String() string {
	switch v := u.value.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case string:
		return v
	default:
		return ""
	}
}

func (u *IntString) Int() int {
	switch v := u.value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
		return 0
	default:
		return 0
	}
}
