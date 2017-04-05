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

	e6 := Position{E, 6}

	assert.Equal(t, e6.AN(), "e6")

	ra1 := Square{Position{A, 1}, Piece{White, Rook}}

	assert.Equal(t, ra1.AN(), "Ra1")
}
