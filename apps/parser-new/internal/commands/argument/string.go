package argument

// String is an argument of string type.
type String struct {
	name     string
	value    string
	optional bool
}

var _ Argument = (*String)(nil)

func NewString(name string, optional bool) *String {
	return &String{
		name:     name,
		optional: optional,
	}
}

func (s *String) Name() string   { return s.name }
func (s *String) Type() ArgType  { return ArgTypeString }
func (s *String) Value() any     { return s.value }
func (s *String) Optional() bool { return s.optional }
