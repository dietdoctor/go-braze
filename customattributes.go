package braze

import "time"

type SliceAttributeAction string

const (
	SliceAttributeActionAdd    SliceAttributeAction = "add"
	SliceAttributeActionRemove SliceAttributeAction = "remove"
)

type CustomAttribute struct {
	key   string
	value any
}

func (a *CustomAttribute) Key() string {
	return a.key
}

func (a *CustomAttribute) Value() any {
	return a.value
}

// Attribute returns a custom attribute.
func Attribute[T any](key string, value T) CustomAttribute {
	return CustomAttribute{key: key, value: value}
}

// BoolAttribute returns a bool-valued attribute.
func BoolAttribute(key string, value bool) CustomAttribute {
	return Attribute(key, value)
}

// Int64Attribute returns an int64-valued attribute.
func Int64Attribute(key string, value int64) CustomAttribute {
	return Attribute(key, value)
}

// Float64Attribute returns a float64-valued attribute.
func Float64Attribute(key string, value float64) CustomAttribute {
	return Attribute(key, value)
}

// StringAttribute returns a string-valued attribute.
func StringAttribute(key string, value string) CustomAttribute {
	return Attribute(key, value)
}

func DateAttribute(key string, value time.Time) CustomAttribute {
	return Attribute(key, value.Format(time.RFC3339))
}

func StringSliceAttribute(key string, value []string) CustomAttribute {
	return Attribute(key, value)
}

func ModifyStringSliceAttribute(key string, value map[SliceAttributeAction][]string) CustomAttribute {
	return Attribute(key, value)
}
