package main

// Farm is a farm in stardew valley
type Farm struct {
	layout             []byte
	remainingResources *Resources
}

// NewEmptyFarm creates a untouched plot of land with starting resources
func NewEmptyFarm(resources *Resources) *Farm {
	out := make([]byte, FARM_SIZE)
	for i := 0; i < FARM_SIZE; i++ {
		out[i] = '.'
	}
	return &Farm{layout: out, remainingResources: resources}
}

// NewFarm creates a farm with starting resources and layout
func NewFarm(resources *Resources, layout []byte) *Farm {
	return &Farm{layout, resources}
}

// Score represents how much crops we can grow
func (f Farm) Score() int {
	score := 0
	for i := 0; i < FARM_SIZE; i++ {
		if f.layout[i] == 'c' {
			score++
		}
	}
	return score
}

// Render makes a pretty output
func (f Farm) Render() string {
	out := ""
	for i := 0; i < FARM_SIZE; i++ {
		if i%FARM_WIDTH == 0 {
			out += "\n"
		}
		out += string(f.layout[i]) + " "
	}
	return out
}
