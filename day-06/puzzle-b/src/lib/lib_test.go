package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTurnRight(t *testing.T) {
	assert.Equal(t, Coord{Row: 0, Col: -1}, TurnRight(Coord{Row: 1, Col: 0}))
	assert.Equal(t, Coord{Row: 1, Col: 0}, TurnRight(Coord{Row: 0, Col: 1}))
	assert.Equal(t, Coord{Row: 0, Col: 1}, TurnRight(Coord{Row: -1, Col: 0}))
	assert.Equal(t, Coord{Row: -1, Col: 0}, TurnRight(Coord{Row: 0, Col: -1}))
}
