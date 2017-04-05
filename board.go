package pawn

/*
The vertical columns of squares (called files) from White's left (the queenside) to right (the kingside) are labeled a through h.
The horizontal rows of squares (called ranks) are numbered 1 to 8, starting from White's side of the board.
Thus each square has a unique identification of file letter followed by rank number.
(For example, White's king starts the game on square e1; Black's knight on b8 can move to open squares a6 or c6.)
*/

type Board struct {
	Rows []string
}
