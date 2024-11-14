package argument

type ArgType int

const (
	ArgTypeInt ArgType = iota
	ArgTypeString
	ArgTypeVariadicString
)

type Argument interface {
	Name() string
	Type() ArgType
	Value() any
	Optional() bool
}
