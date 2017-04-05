package pawn

import "fmt"

type Square struct {
	File string // vertical columns a through h from queenside to kingside
	Rank uint8  // horizontal rows 1 to 8 from White's side of the board
}

func (s Square) AN() string {
	return fmt.Sprintf("%s%d", s.File, s.Rank)
}
