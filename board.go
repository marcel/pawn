package pawn

import (
	"errors"
)

var (
	ErrorMoveByWrongColor = errors.New("pawn: move by wrong color")
	ErrorInvalidPosition  = errors.New("pawn: invalid position")
)

type Board struct {
	Squares    []*Square
	turnNumber int
}

func NewBoard() *Board {
	return &Board{Squares: AllSquares()}
}

func (b Board) turnToMove() Color {
	return colors[b.turnNumber%len(colors)]
}

func (b *Board) MoveFromAlgebraic(an AlgebraicNotation) (Move, error) {
	defer b.incrementTurnNumber()

	piece := Piece{b.turnToMove(), an.Material()}
	destinationPosition := an.DestinationPosition()
	destinationSquare := b.SquareAtPosition(destinationPosition)

	squaresForPiece := b.SquaresForPiece(piece)

	var fromSquare *Square
	for _, squareForPiece := range squaresForPiece {
		possiblePaths := squareForPiece.possiblePaths()
		for _, path := range possiblePaths {
			for _, position := range path {
				if position == destinationPosition {
					fromSquare = squareForPiece
					break
				}
			}
		}
	}

	move := Move{
		From:  fromSquare.Position,
		To:    destinationPosition,
		Piece: piece,
	}

	pieceToMove := fromSquare.Piece
	fromSquare.Piece = NoPiece
	destinationSquare.Piece = pieceToMove

	return move, nil
}

func (b *Board) incrementTurnNumber() {
	b.turnNumber++
}

func (b Board) SquareAtPosition(position Position) *Square {
	index := position.File.index()*len(allFiles) + (int(position.Rank) - 1)

	return b.Squares[index]
}

func (b Board) SquaresForPiece(piece Piece) []*Square {
	squares := []*Square{}

	for _, square := range b.Squares {
		if square.Piece == piece {
			squares = append(squares, square)
		}
	}

	return squares
}

type Move struct {
	Piece
	From  Position
	To    Position
	Takes bool
}

type Game struct {
	Board
	Moves []Move
}
