package pawn

type Material int

const (
	Pawn Material = iota
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
	White       = "White"
)

type Piece struct {
	Color
	Material
}
