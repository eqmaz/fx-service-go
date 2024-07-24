package pseudo

type UnionType interface {
	From(interface{}) error
	GetType() string
	Value() interface{}
}

// NewIntString creates a new IntString union pseudo-type representing an integer or a string
func NewIntString[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string](value T) IntString {
	var result IntString
	err := result.From(value)
	if err != nil {
		return IntString{}
	}
	return result
}

// NewNumber creates a new Number pseudo-type, representing any numerical value (integer or float)
func NewNumber[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](num T) Number {
	var n Number
	err := n.From(num)
	if err != nil {
		return Number{}
	}
	return n
}

func NewInteger[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](value T) Integer {
	var result Integer
	err := result.From(value)
	if err != nil {
		return Integer{}
	}
	return result
}
