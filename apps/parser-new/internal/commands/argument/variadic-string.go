package argument

// VariadicString is an argument of variadic string type.
type VariadicString struct {
	name     string
	value    []string
	optional bool
}

var _ Argument = (*VariadicString)(nil)

func NewVariadicString(name string, optional bool) *VariadicString {
	return &VariadicString{
		name:     name,
		optional: optional,
	}
}

func (vs *VariadicString) Name() string   { return vs.name }
func (vs *VariadicString) Type() ArgType  { return ArgTypeVariadicString }
func (vs *VariadicString) Value() any     { return vs.value }
func (vs *VariadicString) Optional() bool { return vs.optional }
