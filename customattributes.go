package braze

import "time"

type SliceAttributeAction string

const (
	SliceAttributeActionAdd    SliceAttributeAction = "add"
	SliceAttributeActionRemove SliceAttributeAction = "remove"
)

type CustomAttribute struct {
	key   string
	value interface{}
}

func (a *CustomAttribute) Key() string {
	return a.key
}

func (a *CustomAttribute) Value() interface{} {
	return a.value
}

// BoolAttribute returns a bool-valued attribute.
func BoolAttribute(key string, value bool) CustomAttribute {
	return CustomAttribute{key: key, value: value}
}

// Int64Attribute returns an int64-valued attribute.
func Int64Attribute(key string, value int64) CustomAttribute {
	return CustomAttribute{key: key, value: value}
}

// Float64Attribute returns a float64-valued attribute.
func Float64Attribute(key string, value float64) CustomAttribute {
	return CustomAttribute{key: key, value: value}
}

// StringAttribute returns a string-valued attribute.
func StringAttribute(key string, value string) CustomAttribute {
	return CustomAttribute{key: key, value: value}
}

func DateAttribute(key string, value time.Time) CustomAttribute {
	return CustomAttribute{key: key, value: value.Format(time.RFC3339)}
}

func StringSliceAttribute(key string, value []string) CustomAttribute {
	return CustomAttribute{key: key, value: value}
}

func ModifyStringSliceAttribute(key string, value map[SliceAttributeAction][]string) CustomAttribute {
	return CustomAttribute{key: key, value: value}
}
