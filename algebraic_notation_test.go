package pawn

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AlgebraicNotationTestSuite struct {
	suite.Suite
}

func TestAlgebraicNotationTestSuite(t *testing.T) {
	suite.Run(t, new(AlgebraicNotationTestSuite))
}

func (s *AlgebraicNotationTestSuite) TestAN() {
	whiteRook := Piece{White, Rook}

	s.Equal(whiteRook.Material.AN(), "R")
	s.Equal(whiteRook.AN(), "R")
	s.Equal(whiteRook.FAN(), "♜")
	s.Equal(NoPiece.FAN(), " ")

	e6 := Position{E, 6}

	s.Equal(e6.AN(), "e6")

	ra1 := Square{Position{A, 1}, Piece{White, Rook}}

	s.Equal(ra1.AN(), "Ra1")
}

func (s *AlgebraicNotationTestSuite) TestMaterial() {
	expectations := map[AlgebraicNotation]Material{
		"d4":    Pawn,
		"Nf6":   Knight,
		"Bb4":   Bishop,
		"Bxc3+": Bishop,
		"Ra2":   Rook,
		"Qxc1":  Queen,
		"b1=Q+": Pawn, // Promotion to Queen
		"bxc3":  Pawn,
		"cxd4":  Pawn,

		// TODO
		// "O-O"
	}

	for an, expectedMaterial := range expectations {
		s.Equal(an.Material(), expectedMaterial)
	}
}

func (s *AlgebraicNotationTestSuite) TestTakes() {
	expectations := map[AlgebraicNotation]bool{
		"Bxc3+": true,
		"bxc3":  true,
		"cxd5":  true,
		"exd5":  true,
		"Nxc1":  true,
		"Qxc1":  true,
		"axb4":  true,
		"Rxa6":  true,

		"b1=Q+": false,
		"Nf6":   false,
		"c5":    false,
		"O-O":   false,
		"Bg2":   false,
		"Ng3":   false,
	}

	for an, didTake := range expectations {
		s.Equal(an.Takes(), didTake)
	}
}

func (s *AlgebraicNotationTestSuite) TestDestinationPosition() {
	expectations := map[AlgebraicNotation]Position{
		"d4":    Position{D, 4},
		"Nf6":   Position{F, 6},
		"Bxc3+": Position{C, 3},
		"cxd5":  Position{D, 5},
		"Nxc1":  Position{C, 1},
		"b1=Q+": Position{B, 1},
		"Qb6+":  Position{B, 6},
		"f1=Q":  Position{F, 1},

		// TODO
		// O-O
	}

	for an, expectedDestinationPosition := range expectations {
		s.Equal(an.DestinationPosition(), expectedDestinationPosition)
	}
}

func (s *AlgebraicNotationTestSuite) TestPromotion() {
	expectations := map[AlgebraicNotation]Material{
		"d4":    Pawn,
		"Nf6":   Pawn,
		"Bxc3+": Pawn,
		"cxd5":  Pawn,
		"Nxc1":  Pawn,
		"Qb6+":  Pawn,

		"b1=Q+": Queen,
		"b1=R+": Rook,
		"f1=Q":  Queen,
		"f1=R":  Rook,
	}

	for an, expectedMaterial := range expectations {
		s.Equal(an.PromotedTo(), expectedMaterial)

		if expectedMaterial == Pawn {
			s.False(an.IsPromotion())
		} else {
			s.True(an.IsPromotion())
		}
	}
}

func (s *AlgebraicNotationTestSuite) TestIsCheck() {
	isCheck := []AlgebraicNotation{
		"Bxc3+", "Qb6+", "b1=Q+",
	}

	isNotCheck := []AlgebraicNotation{
		"d4", "Nf6", "cxd5", "Nxc1", "f1=Q",
	}

	for _, an := range isCheck {
		s.True(an.IsCheck())
	}

	for _, an := range isNotCheck {
		s.False(an.IsCheck())
	}
}

func (s *AlgebraicNotationTestSuite) TestIsCheckMate() {
	isCheckMate := []AlgebraicNotation{
		"Bxc3#", "Qb6#", "b1=Q#",
	}

	isNotCheckMate := []AlgebraicNotation{
		"d4", "Nf6", "cxd5", "Nxc1", "f1=Q", "Bxc3+", "Qb6+", "b1=Q+",
	}

	for _, an := range isCheckMate {
		s.True(an.IsCheckMate())
	}

	for _, an := range isNotCheckMate {
		s.False(an.IsCheckMate())
	}
}

func (s *AlgebraicNotationTestSuite) TestIsCastle() {
	var kingSide AlgebraicNotation = "O-O"
	var queenSide AlgebraicNotation = "O-O-O"

	s.True(kingSide.IsCastleKingSide())
	s.False(kingSide.IsCastleQueenSide())

	s.True(queenSide.IsCastleQueenSide())
	s.False(queenSide.IsCastleKingSide())

	for _, an := range []AlgebraicNotation{kingSide, queenSide} {
		s.True(an.IsCastle())
	}

	for _, an := range []AlgebraicNotation{"d4", "Nf6", "Bxc3", "cxd4", "Qb26+"} {
		s.False(an.IsCastle())
	}
}

func (s *AlgebraicNotationTestSuite) TestOriginRank() {
	expectations := map[AlgebraicNotation]Rank{
		"5f3":    5,
		"N5f3":   5,
		"N5xf3":  5,
		"Ra8b8":  8,
		"Ra8xb8": 8,

		"Rc1+": NilRank,
		"Rxb7": NilRank,
		"Nc3":  NilRank,
		"Be3+": NilRank,
		"c4":   NilRank,
	}

	for an, expectedRank := range expectations {
		s.Equal(expectedRank, an.OriginRank())
		s.Equal(expectedRank != NilRank, an.OriginDisambiguated())
	}
}

func (s *AlgebraicNotationTestSuite) TestOriginFile() {
	expectations := map[AlgebraicNotation]File{
		"cd5":    C,
		"cxd5":   C,
		"Rfd1":   F,
		"Raxb8":  A,
		"Ra8b8":  A,
		"Ra8xb8": A,

		"Rc1+": NilFile,
		"Rxb7": NilFile,
		"Nc3":  NilFile,
		"Be3+": NilFile,
		"c4":   NilFile,
	}

	for an, expectedFile := range expectations {
		s.Equal(expectedFile, an.OriginFile())
		s.Equal(expectedFile != NilFile, an.OriginDisambiguated())
	}
}

// func TestAlgebraicNotationParse(t *testing.T) {
// 	expectations := map[AlgebraicNotation]Move{
// 		"d4": Move{
// 			Material: Pawn,
// 			Position: Position{D, 4}
// 		}
// 	}

// 	for an, expectedMove := range expectations {
// 		assert.Equal(t, an.Parse(), expectedMove)
// 	}
// }
