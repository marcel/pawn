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

var materialFigurines = map[Color]map[Material]string{
	Black: map[Material]string{
		Pawn:   "♟",
		Rook:   "♜",
		Bishop: "♝",
		Knight: "♞",
		Queen:  "♛",
		King:   "♚",
	},
	White: map[Material]string{
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
