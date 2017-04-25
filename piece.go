package pawn

type Material uint8

const (
	Pawn Material = iota + 1
	Rook
	Bishop
	Knight
	Queen
	King
)

var materialNames = map[Material]string{
	Pawn:   "",
	Rook:   "R",
	Bishop: "B",
	Knight: "N",
	Queen:  "Q",
	King:   "K",
}

// N.B. The unicode for black is used for white and visa versa
// because against a terminal window black looks white and white looks
// black
var materialFigurines = map[Color]map[Material]string{
	White: map[Material]string{
		Pawn:   "♟",
		Rook:   "♜",
		Bishop: "♝",
		Knight: "♞",
		Queen:  "♛",
		King:   "♚",
	},
	Black: map[Material]string{
		Pawn:   "♙",
		Rook:   "♖",
		Bishop: "♗",
		Knight: "♘",
		Queen:  "♕",
		King:   "♔",
	},
}

type Color string

const (
	Black Color = "Black"
	White Color = "White"
)

var colors = [2]Color{White, Black}

type Piece struct {
	Color
	Material
}

var NoPiece = Piece{}
