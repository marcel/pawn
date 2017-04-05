package pawn

import "fmt"

type AlgebraicNotation interface {
	AN() string
	FAN() string
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
	return fmt.Sprintf("%s%s", s.Piece.AN(), s.Position.AN())
}
