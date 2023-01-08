package grammars

type compose struct {
	list []ComposeElement
}

func createCompose(
	list []ComposeElement,
) Compose {
	out := compose{
		list: list,
	}

	return &out
}

// List returns the list of elements
func (obj *compose) List() []ComposeElement {
	return obj.list
}
