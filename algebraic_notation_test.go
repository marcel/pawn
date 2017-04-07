package pawn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlgebraicNotation(t *testing.T) {
	assert := assert.New(t)

	whiteRook := Piece{White, Rook}

	assert.Equal(whiteRook.Material.AN(), "R")
	assert.Equal(whiteRook.AN(), "R")
	assert.Equal(whiteRook.FAN(), "♖")

	e6 := Position{E, 6}

	assert.Equal(e6.AN(), "e6")

	ra1 := Square{Position{A, 1}, Piece{White, Rook}}

	assert.Equal(ra1.AN(), "Ra1")
}

func TestAlgebraicNotationMaterial(t *testing.T) {
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
		assert.Equal(t, an.Material(), expectedMaterial)
	}
}

func TestAlgebraicNotationTakes(t *testing.T) {
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
		assert.Equal(t, an.Takes(), didTake)
	}
}

func TestMoveFromAlgebraicNotationDestinationPosition(t *testing.T) {
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
		assert.Equal(t, an.DestinationPosition(), expectedDestinationPosition)
	}
}

func TestAlgebraicNotationPromotion(t *testing.T) {
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
		assert.Equal(t, an.PromotedTo(), expectedMaterial)

		if expectedMaterial == Pawn {
			assert.False(t, an.IsPromotion())
		} else {
			assert.True(t, an.IsPromotion())
		}
	}
}

func TestAlgebraicNotationIsCheck(t *testing.T) {
	isCheck := []AlgebraicNotation{
		"Bxc3+", "Qb6+", "b1=Q+",
	}

	isNotCheck := []AlgebraicNotation{
		"d4", "Nf6", "cxd5", "Nxc1", "f1=Q",
	}

	for _, an := range isCheck {
		assert.True(t, an.IsCheck())
	}

	for _, an := range isNotCheck {
		assert.False(t, an.IsCheck())
	}
}

func TestAlgebraicNotationIsCheckMate(t *testing.T) {
	isCheckMate := []AlgebraicNotation{
		"Bxc3#", "Qb6#", "b1=Q#",
	}

	isNotCheckMate := []AlgebraicNotation{
		"d4", "Nf6", "cxd5", "Nxc1", "f1=Q", "Bxc3+", "Qb6+", "b1=Q+",
	}

	for _, an := range isCheckMate {
		assert.True(t, an.IsCheckMate())
	}

	for _, an := range isNotCheckMate {
		assert.False(t, an.IsCheckMate())
	}
}

func TestAlgebraicNotationIsCastle(t *testing.T) {
	var kingSide AlgebraicNotation = "O-O"
	var queenSide AlgebraicNotation = "O-O-O"

	assert.True(t, kingSide.IsCastleKingSide())
	assert.False(t, kingSide.IsCastleQueenSide())

	assert.True(t, queenSide.IsCastleQueenSide())
	assert.False(t, queenSide.IsCastleKingSide())

	for _, an := range []AlgebraicNotation{kingSide, queenSide} {
		assert.True(t, an.IsCastle())
	}

	for _, an := range []AlgebraicNotation{"d4", "Nf6", "Bxc3", "cxd4", "Qb26+"} {
		assert.False(t, an.IsCastle())
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
