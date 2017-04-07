package pawn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSquare(t *testing.T) {
	expectations := map[Position]Square{
		// White positions
		Position{A, 1}: Square{Position{A, 1}, Piece{White, Rook}},
		Position{A, 2}: Square{Position{A, 2}, Piece{White, Pawn}},
		Position{D, 1}: Square{Position{D, 1}, Piece{White, Queen}},
		Position{D, 2}: Square{Position{D, 2}, Piece{White, Pawn}},

		// Empty positions
		Position{A, 3}: Square{Position: Position{A, 3}},
		Position{D, 3}: Square{Position: Position{D, 3}},

		// Black positions
		Position{A, 8}: Square{Position{A, 8}, Piece{Black, Rook}},
		Position{A, 7}: Square{Position{A, 7}, Piece{Black, Pawn}},
		Position{D, 8}: Square{Position{D, 8}, Piece{Black, Queen}},
		Position{D, 7}: Square{Position{D, 7}, Piece{Black, Pawn}},
	}

	for position, square := range expectations {
		assert.Equal(t, square, *NewSquare(position))
	}
}

func TestAllSquares(t *testing.T) {
	allSquares := AllSquares()

	assert.Equal(t, len(allFiles)*len(allRanks), len(allSquares))
}
