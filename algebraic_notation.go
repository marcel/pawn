package pawn

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
