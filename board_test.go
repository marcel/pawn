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
	// 1.e4 Nf6 2.e5 Nd5 3.d4 d6 4.Nf3 Bg4 5.Bc4 e6 6.O-O Nb6

	AlgebraicMoveAssertion{
		"e4",
		Piece{White, Pawn},
		E2,
		E4,
		s,
	}.assert()

	AlgebraicMoveAssertion{
		"Nf6",
		Piece{Black, Knight},
		G8,
		F6,
		s,
	}.assert()

	AlgebraicMoveAssertion{
		"e5",
		Piece{White, Pawn},
		E4,
		E5,
		s,
	}.assert()

	AlgebraicMoveAssertion{
		"Nd5",
		Piece{Black, Knight},
		F6,
		D5,
		s,
	}.assert()
}

type AlgebraicMoveAssertion struct {
	an    AlgebraicNotation
	piece Piece
	from  Position
	to    Position
	suite *BoardTestSuite
}

func (a AlgebraicMoveAssertion) assert() {
	turnNumberBeforeMove := a.suite.board.turnNumber
	originSquareBeforeMove := *a.suite.board.SquareAtPosition(a.from)
	destinationSquareBeforeMove := *a.suite.board.SquareAtPosition(a.to)

	a.suite.Equal(a.piece, originSquareBeforeMove.Piece)
	a.suite.Equal(NoPiece, destinationSquareBeforeMove.Piece)

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
