package main

// Resources for a farm
type Resources struct {
	sprinklers     int
	qSprinklers    int
	iridSprinklers int
	scarecrows     int
}

// NewResources returns a new list of resources
func NewResources(sprinklers int, qSprinklers int, iridSprinklers int, scarecrows int) *Resources {
	return &Resources{sprinklers, qSprinklers, iridSprinklers, scarecrows}
}

// UseSprinkler returns remaining resources after using a sprinkler
func (r Resources) UseSprinkler() *Resources {
	return NewResources(r.sprinklers-1, r.qSprinklers, r.iridSprinklers, r.scarecrows)
}

// UseQSprinkler returns remaining resources after using a q sprinkler
func (r Resources) UseQSprinkler() *Resources {
	return NewResources(r.sprinklers, r.qSprinklers-1, r.iridSprinklers, r.scarecrows)
}

// UseIridSprinkler returns remaining resources after using a irid sprinkler
func (r Resources) UseIridSprinkler() *Resources {
	return NewResources(r.sprinklers, r.qSprinklers, r.iridSprinklers-1, r.scarecrows)
}

// UseScarecrow returns remaining resources after using a scarecrow
func (r Resources) UseScarecrow() *Resources {
	return NewResources(r.sprinklers, r.qSprinklers, r.iridSprinklers, r.scarecrows-1)
}

// Options are available resources
func (r Resources) Options() []byte {
	out := make([]byte, 0)

	if r.sprinklers > 0 {
		out = append(out, 's')
	}

	if r.qSprinklers > 0 {
		out = append(out, 'q')
	}

	if r.iridSprinklers > 0 {
		out = append(out, 'i')
	}

	if r.scarecrows > 0 {
		out = append(out, 'S')
	}

	return out
}
