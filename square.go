package pawn

type File string

const (
	A File = "a"
	B      = "b"
	C      = "c"
	D      = "d"
	E      = "e"
	F      = "f"
	G      = "g"
	H      = "h"
)

var allFiles = []File{
	A, B, C, D, E, F, G, H,
}

type Rank uint8

var allRanks = []Rank{
	1, 2, 3, 4, 5, 6, 7, 8,
}

type Position struct {
	File // vertical columns a through h from queenside to kingside
	Rank // horizontal rows 1 to 8 from White's side of the board
}

func allPositions() []Position {
	positions := make([]Position, len(allFiles)*len(allRanks))

	for rank := range allRanks {
		for file := range allFiles {
			positions = append(positions, Position{File: File(file), Rank: Rank(rank)})
		}
	}

	return positions
}

type Square struct {
	Position
	Piece
}

var materialByFile = map[File]Material{
	A: Rook,
	B: Knight,
	C: Bishop,
	D: Queen,
	E: King,
	F: Bishop,
	G: Knight,
	H: Rook,
}

func NewSquare(position Position) Square {
	var color Color
	var material Material

	switch position.Rank {
	case 7, 8:
		color = Black
	case 1, 2:
		color = White
	}

	switch position.Rank {
	case 2, 7:
		material = Pawn
	case 1, 8:
		material = materialByFile[position.File]
	}

	if color != "" && material != 0 {
		return Square{Position: position, Piece: Piece{color, material}}
	} else {
		return Square{Position: position}
	}
}
