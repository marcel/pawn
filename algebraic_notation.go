package pawn

import (
	"fmt"
	"strconv"
)

type AlgebraicNotation string

type AlgebraiclyNotated interface {
	AN() string
	FAN() string
}

func MoveFromAN(an AlgebraicNotation) Move {
	var piece Piece
	position := positionFromAN(an)

	if len(an) == 2 {
		piece = Piece{White, Pawn} // TODO How to figure out Color?

		return Move{piece, position, false}
	}

	return Move{}
}

func positionFromAN(an AlgebraicNotation) Position {
	substring := an[len(an)-2:]
	file := File(string(substring[0]))
	rankInt, _ := strconv.ParseUint(string(substring[1]), 10, 8)
	rank := Rank(rankInt)

	return Position{File: file, Rank: rank}
}

func (m Material) AN() string {
	return materialNames[m]
}

func (p Piece) AN() string {
	return p.Material.AN()
}

func (p Piece) FAN() string {
	return materialFigurines[p.Color][p.Material]
}

func (p Position) AN() string {
	return fmt.Sprintf("%s%d", p.File, p.Rank)
}

func (s Square) AN() string {
	return s.Piece.AN() + s.Position.AN()
}

func (m Move) AN() string {
	s := m.Piece.AN()

	if m.Takes {
		s += "x"
	}

	s += m.Position.AN()

	return s
}
