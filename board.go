package pawn

type Board struct {
	Squares []*Square
}

func NewBoard() Board {
	return Board{Squares: AllSquares()}
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
	Piece
	Position
	Takes bool
}

func (p Piece) Takes(position Position) Move {
	return Move{Piece: p, Position: position, Takes: true}
}

type Game struct {
	Board
	Moves []Move
}

// func (g *Game) MoveFromTo(fromSquare, toSquare *Square) {
//
// }
//
// func (g *Game) addMove(move Move) {
//   g.Moves = append(g.Moves, move)
// }
