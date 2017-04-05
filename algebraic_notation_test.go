package pawn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlgebraicNotation(t *testing.T) {
	whiteRook := Piece{White, Rook}

	assert.Equal(t, whiteRook.Material.AN(), "R")
	assert.Equal(t, whiteRook.AN(), "R")
	assert.Equal(t, whiteRook.FAN(), "♖")

	e6 := Square{"e", 6}

	assert.Equal(t, e6.AN(), "e6")
}
