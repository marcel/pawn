package pawn

import (
	"errors"
	"sync"
)

var (
	ErrorMoveByWrongColor = errors.New("pawn: move by wrong color")
	ErrorInvalidPosition  = errors.New("pawn: invalid position")
)

type Board struct {
	Squares    []*Square
	turnNumber int
	moveMutex  *sync.Mutex
}

func NewBoard() *Board {
	return &Board{Squares: AllSquares(), moveMutex: &sync.Mutex{}}
}

// Returns 8 rows of 8 squares each starting at the top left and moving down
func (b Board) Rows() [][]*Square {
	rows := [][]*Square{}
	for _, rank := range rankSlice(allRanks[:]).reverse() {
		row := []*Square{}
		for _, file := range allFiles {
			row = append(row, b.SquareAtPosition(Position{file, rank}))
		}
		rows = append(rows, row)
	}

	return rows
}

func (b Board) turnToMove() Color {
	return colors[b.turnNumber%len(colors)]
}

func (b *Board) MoveFromAlgebraic(an AlgebraicNotation) (Move, error) {
	b.moveMutex.Lock()
	defer b.moveMutex.Unlock()
	defer b.incrementTurnNumber()

	switch {
	case an.IsCastle():
		return b.castle(an)
	case an.Takes() && an.Material() == Pawn:
		return b.pawnTakes(an)
	default:
		return b.moveFromAlgebraic(an)
	}
}

func (b *Board) moveFromAlgebraic(an AlgebraicNotation) (Move, error) {
	piece := Piece{b.turnToMove(), an.Material()}
	destinationPosition := an.DestinationPosition()
	destinationSquare := b.SquareAtPosition(destinationPosition)

	squaresForPiece := b.SquaresForPiece(piece)

	var fromSquare *Square

OriginSquareSearch:
	for _, squareForPiece := range squaresForPiece {
		possiblePaths := squareForPiece.possiblePaths()
		for _, path := range possiblePaths {
			for _, position := range path {
				if position == destinationPosition {
					fromSquare = squareForPiece
					if an.OriginDisambiguated() {
						if squareForPiece.File == an.OriginFile() ||
							squareForPiece.Rank == an.OriginRank() {
							break OriginSquareSearch
						}
					} else {
						break OriginSquareSearch
					}
				}
			}
		}
	}

	if fromSquare == nil {
		panic(string(an))
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

func (b *Board) castle(an AlgebraicNotation) (Move, error) {
	var rank Rank

	if b.turnToMove() == White {
		rank = 1
	} else {
		rank = 8
	}

	var kingFrom, kingTo, rookFrom, rookTo Position

	// TODO Moves usually only move one piece. Deal with Move struct
	// encoding that both King and Rook moved?

	switch {
	case an.IsCastleKingSide():
		kingFrom = Position{E, rank}
		kingTo = Position{G, rank}
		rookFrom = Position{H, rank}
		rookTo = Position{F, rank}
	case an.IsCastleQueenSide():
		kingFrom = Position{E, rank}
		kingTo = Position{C, rank}
		rookFrom = Position{A, rank}
		rookTo = Position{D, rank}
	}

	move := Move{
		From:  kingFrom,
		To:    kingTo,
		Piece: Piece{b.turnToMove(), King},
	}

	b.SquareAtPosition(kingFrom).Piece = NoPiece
	b.SquareAtPosition(rookFrom).Piece = NoPiece

	b.SquareAtPosition(kingTo).Piece = Piece{b.turnToMove(), King}
	b.SquareAtPosition(rookTo).Piece = Piece{b.turnToMove(), Rook}

	return move, nil
}

// TODO Handle en passant
func (b *Board) pawnTakes(an AlgebraicNotation) (Move, error) {
	piece := Piece{b.turnToMove(), an.Material()}
	destinationPosition := an.DestinationPosition()
	destinationSquare := b.SquareAtPosition(destinationPosition)

	squaresForPiece := b.SquaresForPiece(piece)

	var fromSquare *Square
OriginSquareSearch:
	for _, squareForPiece := range squaresForPiece {
		pathsToTake := squareForPiece.pathsToTake()
		for _, path := range pathsToTake {
			position := path[0]
			if position == destinationPosition {
				fromSquare = squareForPiece
				if an.OriginDisambiguated() {
					if squareForPiece.File == an.OriginFile() ||
						squareForPiece.Rank == an.OriginRank() {
						break OriginSquareSearch
					}
				} else {
					break OriginSquareSearch
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
