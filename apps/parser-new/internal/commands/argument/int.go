package argument

// Int is an argument of integer type.
type Int struct {
	name     string
	value    int
	optional bool
}

var _ Argument = (*Int)(nil)

func NewInt(name string, optional bool) *Int {
	return &Int{
		name:     name,
		optional: optional,
	}
}

func (i *Int) Name() string   { return i.name }
func (i *Int) Type() ArgType  { return ArgTypeInt }
func (i *Int) Value() any     { return i.value }
func (i *Int) Optional() bool { return i.optional }
