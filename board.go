package pawn

import "errors"

var (
	ErrorMoveByWrongColor = errors.New("pawn: move by wrong color")
)

type Board struct {
	Squares []*Square

	turnToMove Color
}

func NewBoard() Board {
	return Board{Squares: AllSquares(), turnToMove: White}
}

func (b *Board) MoveFromAlgebraic(an AlgebraicNotation, color Color) (Move, error) {
	if b.turnToMove != color {
		return Move{}, ErrorMoveByWrongColor
	}

	return Move{}, nil // TODO Placeholder while implementing
}

func (b Board) SquareAtPosition(position Position) (squareAtPosition *Square) {
	for _, square := range b.Squares {
		if square.Position == position {
			squareAtPosition = square
			return
		}
	}

	return
}

type Move struct {
	Material
	Piece
	Position
	Takes bool
}

type Game struct {
	Board
	Moves []Move
}
