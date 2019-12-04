package snake

type deathBoundary map[xPos]map[yPos]wallExists

func (b deathBoundary) Add(x, y int) {
	_, ok := b[xPos(x)]
	if !ok {
		b[xPos(x)] = map[yPos]wallExists{}
	}

	b[xPos(x)][yPos(y)] = wallExists{}
}

func (b deathBoundary) Remove(x, y int) {
	_, found := b[xPos(x)]
	if !found {
		return // easy peasy
	}

	delete(b[xPos(x)], yPos(y))
}

func (b deathBoundary) IsBoundary(x, y int) bool {
	_, exists := b[xPos(x)]
	if !exists {
		return false
	}

	_, dead := b[xPos(x)][yPos(y)]
	return dead
}

func (b deathBoundary) Copy() deathBoundary {
	newB := deathBoundary{}

	for x, yRow := range b {
		for y, exists := range yRow {
			if newB[x] == nil {
				newB[x] = map[yPos]wallExists{}
			}
			newB[x][y] = exists
		}
	}

	return newB
}
