package pawn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBoard(t *testing.T) {
	newBoard := NewBoard()

	assert.Equal(t, len(newBoard.Squares), len(allFiles)*len(allRanks))
}

func TestSquareAtPosition(t *testing.T) {
	newBoard := NewBoard()

	position := allPositions()[0]

	squareAtPosition := newBoard.SquareAtPosition(position)

	assert.Equal(t, squareAtPosition.Position, position)
}
