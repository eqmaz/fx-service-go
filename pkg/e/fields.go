package e

// Fields type to represent custom fields on an exception.
type Fields map[string]interface{}

// With adds a value to Fields map, and returns the map.
func (f *Fields) With(key string, value interface{}) Fields {
	(*f)[key] = value
	return *f
}
