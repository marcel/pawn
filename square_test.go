package pawn

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNewSquare(t *testing.T) {
	expectations := map[Position]Square{
		// White positions
		Position{A, 1}: Square{Position{A, 1}, Piece{White, Rook}},
		Position{A, 2}: Square{Position{A, 2}, Piece{White, Pawn}},
		Position{D, 1}: Square{Position{D, 1}, Piece{White, Queen}},
		Position{D, 2}: Square{Position{D, 2}, Piece{White, Pawn}},

		// Empty positions
		Position{A, 3}: Square{Position: Position{A, 3}},
		Position{D, 3}: Square{Position: Position{D, 3}},

		// Black positions
		Position{A, 8}: Square{Position{A, 8}, Piece{Black, Rook}},
		Position{A, 7}: Square{Position{A, 7}, Piece{Black, Pawn}},
		Position{D, 8}: Square{Position{D, 8}, Piece{Black, Queen}},
		Position{D, 7}: Square{Position{D, 7}, Piece{Black, Pawn}},
	}

	for position, square := range expectations {
		assert.Equal(t, square, *NewSquare(position))
	}
}

func TestAllSquares(t *testing.T) {
	allSquares := AllSquares()

	assert.Equal(t, len(allFiles)*len(allRanks), len(allSquares))
}

func TestPathsToTake(t *testing.T) {
	whitePawn := Square{Piece: Piece{White, Pawn}}
	blackPawn := Square{Piece: Piece{Black, Pawn}}

	whitePawn.Position = E2

	assert.Equal(t,
		[]Path{
			Path{F3}, Path{D3},
		},
		whitePawn.pathsToTake(),
	)

	blackPawn.Position = E7
	assert.Equal(t,
		[]Path{
			Path{F6}, Path{D6},
		},
		blackPawn.pathsToTake(),
	)

	whitePawn.Position = A2

	assert.Equal(t,
		[]Path{
			Path{B3},
		},
		whitePawn.pathsToTake(),
	)

	blackPawn.Position = A7

	assert.Equal(t,
		[]Path{
			Path{B6},
		},
		blackPawn.pathsToTake(),
	)

	whitePawn.Position = A8

	assert.Equal(t,
		[]Path{},
		whitePawn.pathsToTake(),
	)

	blackPawn.Position = A1

	assert.Equal(t,
		[]Path{},
		blackPawn.pathsToTake(),
	)

	otherMaterial := []Material{
		Rook, Knight, Bishop, Queen, King,
	}

	for _, material := range otherMaterial {
		for _, color := range colors {
			square := Square{E4, Piece{color, material}}
			assert.Equal(t, square.possiblePaths(), square.pathsToTake())
		}
	}
}

type PossiblePathsTestSuite struct {
	suite.Suite
	square Square
}

func (s *PossiblePathsTestSuite) SetupTest() {
	s.square = Square{Piece: Piece{Color: White}}
}

func TestPossiblePathsTestSuite(t *testing.T) {
	suite.Run(t, new(PossiblePathsTestSuite))
}

func (s *PossiblePathsTestSuite) TestPawn() {
	s.square.Material = Pawn
	s.square.Color = White

	s.square.Position = E5

	s.Equal(
		[]Path{
			Path{E6},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = E2

	s.Equal(
		[]Path{
			Path{E3},
			Path{E4},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = E7

	s.Equal(
		[]Path{
			Path{E8},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = E8

	s.Equal(
		[]Path{},
		s.square.possiblePaths(),
	)

	s.square.Color = Black

	s.square.Position = E5

	s.Equal(
		[]Path{
			Path{E4},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = E2

	s.Equal(
		[]Path{
			Path{E1},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = E7

	s.Equal(
		[]Path{
			Path{E6},
			Path{E5},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = E1

	s.Equal(
		[]Path{},
		s.square.possiblePaths(),
	)
}

func (s *PossiblePathsTestSuite) TestKnight() {
	s.square.Material = Knight

	s.square.Position = E5

	s.Equal(
		[]Path{
			Path{F7},
			Path{G6},
			Path{G4},
			Path{F3},
			Path{D3},
			Path{C4},
			Path{C6},
			Path{D7},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = A1

	s.Equal(
		[]Path{
			Path{B3},
			Path{C2},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = B7

	s.Equal(
		[]Path{
			Path{D8},
			Path{D6},
			Path{C5},
			Path{A5},
		},
		s.square.possiblePaths(),
	)
}

func (s *PossiblePathsTestSuite) TestRook() {
	s.square.Material = Rook

	s.square.Position = E5

	s.Equal(
		[]Path{
			Path{E6, E7, E8},
			Path{F5, G5, H5},
			Path{E4, E3, E2, E1},
			Path{A5, B5, C5, D5},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = A1

	s.Equal(
		[]Path{
			Path{A2, A3, A4, A5, A6, A7, A8},
			Path{B1, C1, D1, E1, F1, G1, H1},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = B7

	s.Equal(
		[]Path{
			Path{B8},
			Path{C7, D7, E7, F7, G7, H7},
			Path{B6, B5, B4, B3, B2, B1},
			Path{A7},
		},
		s.square.possiblePaths(),
	)
}

func (s *PossiblePathsTestSuite) TestBishop() {
	s.square.Material = Bishop

	s.square.Position = E5

	s.Equal(
		[]Path{
			Path{F6, G7, H8},
			Path{F4, G3, H2},
			Path{D4, C3, B2, A1},
			Path{D6, C7, B8},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = A1

	s.Equal(
		[]Path{
			Path{B2, C3, D4, E5, F6, G7, H8},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = B7

	s.Equal(
		[]Path{
			Path{C8},
			Path{C6, D5, E4, F3, G2, H1},
			Path{A6},
			Path{A8},
		},
		s.square.possiblePaths(),
	)
}

func (s *PossiblePathsTestSuite) TestQueen() {
	s.square.Material = Queen

	s.square.Position = E5

	s.Equal(
		[]Path{
			Path{E6, E7, E8},
			Path{F6, G7, H8},
			Path{F5, G5, H5},
			Path{F4, G3, H2},
			Path{E4, E3, E2, E1},
			Path{D4, C3, B2, A1},
			Path{A5, B5, C5, D5},
			Path{D6, C7, B8},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = A1

	s.Equal(
		[]Path{
			Path{A2, A3, A4, A5, A6, A7, A8},
			Path{B2, C3, D4, E5, F6, G7, H8},
			Path{B1, C1, D1, E1, F1, G1, H1},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = B7

	s.Equal(
		[]Path{
			Path{B8},
			Path{C8},
			Path{C7, D7, E7, F7, G7, H7},
			Path{C6, D5, E4, F3, G2, H1},
			Path{B6, B5, B4, B3, B2, B1},
			Path{A6},
			Path{A7},
			Path{A8},
		},
		s.square.possiblePaths(),
	)
}

func (s *PossiblePathsTestSuite) TestKing() {
	s.square.Material = King

	s.square.Position = E5

	s.Equal(
		[]Path{
			Path{E6},
			Path{F6},
			Path{F5},
			Path{F4},
			Path{E4},
			Path{D4},
			Path{D5},
			Path{D6},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = A1

	s.Equal(
		[]Path{
			Path{A2},
			Path{B2},
			Path{B1},
		},
		s.square.possiblePaths(),
	)

	s.square.Position = B7

	s.Equal(
		[]Path{
			Path{B8},
			Path{C8},
			Path{C7},
			Path{C6},
			Path{B6},
			Path{A6},
			Path{A7},
			Path{A8},
		},
		s.square.possiblePaths(),
	)
}
