package pawn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlgebraicNotation(t *testing.T) {
	assert := assert.New(t)

	whiteRook := Piece{White, Rook}

	assert.Equal(whiteRook.Material.AN(), "R")
	assert.Equal(whiteRook.AN(), "R")
	assert.Equal(whiteRook.FAN(), "♖")

	e6 := Position{E, 6}

	assert.Equal(e6.AN(), "e6")

	ra1 := Square{Position{A, 1}, Piece{White, Rook}}

	assert.Equal(ra1.AN(), "Ra1")
}

func TestMoveFromAlgebraicNotation(t *testing.T) {
	expectations := map[AlgebraicNotation]Move{
		"e4": Move{
			Piece:    Piece{White, Pawn},
			Position: Position{E, 4},
			Takes:    false,
		},
	}

	for an, expectedMove := range expectations {
		assert.Equal(t, MoveFromAN(an), expectedMove)
	}

}
