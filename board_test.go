package pawn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	newBoard := NewBoard()

	assert.Equal(t, len(newBoard.Squares), len(allFiles)*len(allRanks))
	assert.Equal(t, newBoard.turnToMove, White)
}

func TestSquareAtPosition(t *testing.T) {
	newBoard := NewBoard()

	position := allPositions[0]

	squareAtPosition := newBoard.SquareAtPosition(position)

	assert.Equal(t, squareAtPosition.Position, position)
}

func TestMoveFromAlgebraic(t *testing.T) {
	newBoard := NewBoard()

	_, err := newBoard.MoveFromAlgebraic("", Black)
	assert.Error(t, err, ErrorMoveByWrongColor)
}
