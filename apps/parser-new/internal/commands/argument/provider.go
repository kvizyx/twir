package argument

// Provider provides a storage (without synchronization) from which you can conveniently
// retrieve arguments that will be automatically cast to the required type.
type Provider struct {
	storage map[string]Argument
}

func NewProvider(arguments []Argument) Provider {
	storage := make(map[string]Argument, len(arguments))

	for _, argument := range arguments {
		storage[argument.Name()] = argument
	}

	return Provider{
		storage: storage,
	}
}

func (p *Provider) Int(name string) (int, bool) {
	if value, ok := p.storage[name]; ok {
		argument, parseable := value.Value().(int)
		return argument, parseable
	}

	return 0, false
}

func (p *Provider) String(name string) (string, bool) {
	if value, ok := p.storage[name]; ok {
		argument, parseable := value.Value().(string)
		return argument, parseable
	}

	return "", false
}

func (p *Provider) VariadicString(name string) ([]string, bool) {
	if value, ok := p.storage[name]; ok {
		argument, parseable := value.Value().([]string)
		return argument, parseable
	}

	return nil, false
}
