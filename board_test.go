package pawn

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BoardTestSuite struct {
	suite.Suite
	board *Board
}

func TestBoardTestSuite(t *testing.T) {
	suite.Run(t, new(BoardTestSuite))
}

func (s *BoardTestSuite) SetupTest() {
	s.board = NewBoard()
}

func (s *BoardTestSuite) TestNewBoard() {
	s.Equal(len(s.board.Squares), len(allFiles)*len(allRanks))
	s.Equal(s.board.turnToMove(), White)
}

func (s *BoardTestSuite) TestRows() {
	rows := s.board.Rows()

	s.Equal(8, len(rows))
	for i := range rows {
		s.Equal(8, len(rows[i]))
	}

	s.Equal(A8, rows[0][0].Position)
	s.Equal(H8, rows[0][7].Position)

	s.Equal(A1, rows[7][0].Position)
	s.Equal(H1, rows[7][7].Position)
}

func (s *BoardTestSuite) TestSquareAtPosition() {
	expectations := map[Position]Square{
		A1: Square{Piece: Piece{White, Rook}, Position: A1},
		A2: Square{Piece: Piece{White, Pawn}, Position: A2},
		A7: Square{Piece: Piece{Black, Pawn}, Position: A7},
		A8: Square{Piece: Piece{Black, Rook}, Position: A8},
		H8: Square{Piece: Piece{Black, Rook}, Position: H8},
	}

	for position, expectedSquare := range expectations {
		s.Equal(expectedSquare, *s.board.SquareAtPosition(position))
	}
}

func (s *BoardTestSuite) TestSquaresForPiece() {
	expectations := map[Piece][]*Square{
		Piece{White, Rook}: []*Square{
			s.board.SquareAtPosition(A1),
			s.board.SquareAtPosition(H1),
		},
		Piece{White, Knight}: []*Square{
			s.board.SquareAtPosition(B1),
			s.board.SquareAtPosition(G1),
		},
		Piece{White, Bishop}: []*Square{
			s.board.SquareAtPosition(C1),
			s.board.SquareAtPosition(F1),
		},
		Piece{White, Queen}: []*Square{
			s.board.SquareAtPosition(D1),
		},
		Piece{White, King}: []*Square{
			s.board.SquareAtPosition(E1),
		},
		Piece{White, Pawn}: []*Square{
			s.board.SquareAtPosition(A2),
			s.board.SquareAtPosition(B2),
			s.board.SquareAtPosition(C2),
			s.board.SquareAtPosition(D2),
			s.board.SquareAtPosition(E2),
			s.board.SquareAtPosition(F2),
			s.board.SquareAtPosition(G2),
			s.board.SquareAtPosition(H2),
		},
		Piece{Black, Rook}: []*Square{
			s.board.SquareAtPosition(A8),
			s.board.SquareAtPosition(H8),
		},
		Piece{Black, Knight}: []*Square{
			s.board.SquareAtPosition(B8),
			s.board.SquareAtPosition(G8),
		},
		Piece{Black, Bishop}: []*Square{
			s.board.SquareAtPosition(C8),
			s.board.SquareAtPosition(F8),
		},
		Piece{Black, Queen}: []*Square{
			s.board.SquareAtPosition(D8),
		},
		Piece{Black, King}: []*Square{
			s.board.SquareAtPosition(E8),
		},
		Piece{Black, Pawn}: []*Square{
			s.board.SquareAtPosition(A7),
			s.board.SquareAtPosition(B7),
			s.board.SquareAtPosition(C7),
			s.board.SquareAtPosition(D7),
			s.board.SquareAtPosition(E7),
			s.board.SquareAtPosition(F7),
			s.board.SquareAtPosition(G7),
			s.board.SquareAtPosition(H7),
		},
	}

	for piece, expectedSquares := range expectations {
		s.Equal(expectedSquares, s.board.SquaresForPiece(piece))
	}
}

func (s *BoardTestSuite) TestMoveFromAlgebraic() {
	// 1.c4 c6 2.g3 d5 3.Bg2 Nf6 4.Nf3 Bf5 5.cxd5 cxd5 6.O-O Nc6 7.d3 e6
	// 8.Be3 Be7 9.Qb3 Qd7 10.Nd4 Nxd4 11.Bxd4 O-O 12.Nc3 Bg6 13.Rfd1 Bd6
	// 14.Nb5 Bb8 15.Bxf6 gxf6 16.e4 dxe4 17.dxe4 Qe7 18.Qa3 Qxa3
	// 19.Nxa3 Rc8 20.Rd7 Bc7 21.b3 Rab8 22.Rc1 Bd8 23.Rxc8 Rxc8
	// 24.Bf3 Rc1+ 25.Kg2 Bb6 26.Rxb7 Ra1 27.Nc4 Bd4 28.Rd7 e5 29.g4 Rxa2
	// 30.h4 Bxf2 31.h5 Be3+ 32.Kh1 Ra1+ 33.Kg2 Ra2+ 34.Kh1 Ra1+
	// 35.Kg2 Ra2+  1/2-1/2

	// 1.
	AlgebraicMoveAssertion{s,
		"c4", Piece{White, Pawn}, C2, C4,
	}.assert()

	AlgebraicMoveAssertion{s,
		"c6", Piece{Black, Pawn}, C7, C6,
	}.assert()

	// 2.
	AlgebraicMoveAssertion{s,
		"g3", Piece{White, Pawn}, G2, G3,
	}.assert()

	AlgebraicMoveAssertion{s,
		"d5", Piece{Black, Pawn}, D7, D5,
	}.assert()

	// 3.
	AlgebraicMoveAssertion{s,
		"Bg2", Piece{White, Bishop}, F1, G2,
	}.assert()

	AlgebraicMoveAssertion{s,
		"Nf6", Piece{Black, Knight}, G8, F6,
	}.assert()

	// 4.
	AlgebraicMoveAssertion{s,
		"Nf3", Piece{White, Knight}, G1, F3,
	}.assert()

	AlgebraicMoveAssertion{s,
		"Bf5", Piece{Black, Bishop}, C8, F5,
	}.assert()

	// 5.
	AlgebraicMoveAssertion{s,
		"cxd5", Piece{White, Pawn}, C4, D5,
	}.assert()

	AlgebraicMoveAssertion{s,
		"cxd5", Piece{Black, Pawn}, C6, D5,
	}.assert()

	// 6.
	AlgebraicMoveAssertion{s,
		"O-O", Piece{White, King}, E1, G1,
	}.assert()

	AlgebraicMoveAssertion{s,
		"Nc6", Piece{Black, Knight}, B8, C6,
	}.assert()

	// 7.
	AlgebraicMoveAssertion{s,
		"d3", Piece{White, Pawn}, D2, D3,
	}.assert()

	AlgebraicMoveAssertion{s,
		"e6", Piece{Black, Pawn}, E7, E6,
	}.assert()

	// 8.
	AlgebraicMoveAssertion{s,
		"Be3", Piece{White, Bishop}, C1, E3,
	}.assert()

	AlgebraicMoveAssertion{s,
		"Be7", Piece{Black, Bishop}, F8, E7,
	}.assert()

	// 9.
	AlgebraicMoveAssertion{s,
		"Qb3", Piece{White, Queen}, D1, B3,
	}.assert()

	AlgebraicMoveAssertion{s,
		"Qd7", Piece{Black, Queen}, D8, D7,
	}.assert()

	// 10.
	AlgebraicMoveAssertion{s,
		"Nd4", Piece{White, Knight}, F3, D4,
	}.assert()

	AlgebraicMoveAssertion{s,
		"Nxd4", Piece{Black, Knight}, C6, D4,
	}.assert()

	// 11.
	AlgebraicMoveAssertion{s,
		"Bxd4", Piece{White, Bishop}, E3, D4,
	}.assert()

	AlgebraicMoveAssertion{s,
		"O-O", Piece{Black, King}, E8, G8,
	}.assert()

	// TODO
	// 12.Nc3 Bg6
	// 13.Rfd1 Bd6
	// 14.Nb5 Bb8
	// 15.Bxf6 gxf6
	// 16.e4 dxe4
	// 17.dxe4 Qe7
	// 18.Qa3 Qxa3
	// 19.Nxa3 Rc8
	// 20.Rd7 Bc7
	// 21.b3 Rab8
	// 22.Rc1 Bd8
	// 23.Rxc8 Rxc8
	// 24.Bf3 Rc1+
	// 25.Kg2 Bb6
	// 26.Rxb7 Ra1
	// 27.Nc4 Bd4
	// 28.Rd7 e5
	// 29.g4 Rxa2
	// 30.h4 Bxf2
	// 31.h5 Be3+
	// 32.Kh1 Ra1+
	// 33.Kg2 Ra2+
	// 34.Kh1 Ra1+
	// 35.Kg2 Ra2+  1/2-1/2
}

type AlgebraicMoveAssertion struct {
	suite *BoardTestSuite
	an    AlgebraicNotation
	piece Piece
	from  Position
	to    Position
}

func (a AlgebraicMoveAssertion) assert() {
	turnNumberBeforeMove := a.suite.board.turnNumber
	originSquareBeforeMove := *a.suite.board.SquareAtPosition(a.from)
	destinationSquareBeforeMove := *a.suite.board.SquareAtPosition(a.to)

	a.suite.Equal(a.piece, originSquareBeforeMove.Piece)

	if a.an.Takes() {
		a.suite.NotEqual(NoPiece, destinationSquareBeforeMove.Piece)
	} else {
		a.suite.Equal(NoPiece, destinationSquareBeforeMove.Piece)
	}

	move, ok := a.suite.board.MoveFromAlgebraic(a.an)

	a.suite.Equal(
		Move{
			From:  a.from,
			To:    a.to,
			Piece: a.piece,
		},
		move,
	)

	turnNumberAfterMove := a.suite.board.turnNumber
	originSquareAfterMove := *a.suite.board.SquareAtPosition(a.from)
	destinationSquareAfterMove := *a.suite.board.SquareAtPosition(a.to)

	a.suite.Equal(turnNumberBeforeMove+1, turnNumberAfterMove)
	a.suite.Equal(NoPiece, originSquareAfterMove.Piece)
	a.suite.Equal(originSquareBeforeMove.Piece, destinationSquareAfterMove.Piece)

	a.suite.Nil(ok)
}
