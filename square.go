package pawn

type Square struct {
	Position
	Piece
}

func (s Square) possiblePaths() []Path {
	switch s.Material {
	// TODO Handle capture case? Or only from perspective of the board?
	case Pawn:
		switch s.Rank {
		case 2:
			switch s.Color {
			case White:
				up1, _ := s.Jump(Up)
				up2, _ := s.Jump(Up, Up)

				return []Path{Path{up1}, Path{up2}}
			case Black:
				down, _ := s.Jump(Down)
				return []Path{Path{down}}
			}

		case 7:
			switch s.Color {
			case Black:
				down1, _ := s.Jump(Down)
				down2, _ := s.Jump(Down, Down)

				return []Path{Path{down1}, Path{down2}}
			case White:
				up, _ := s.Jump(Up)
				return []Path{Path{up}}
			}
		default:
			switch s.Color {
			case White:
				if up, err := s.Jump(Up); err == nil {
					return []Path{Path{up}}
				} else {
					return []Path{}
				}
			case Black:
				if down, err := s.Jump(Down); err == nil {
					return []Path{Path{down}}
				} else {
					return []Path{}
				}
			}
		}
	case Knight:
		possibleJumps := [][]Direction{
			[]Direction{Up, Up, Right},
			[]Direction{Right, Right, Up},
			[]Direction{Right, Right, Down},
			[]Direction{Down, Down, Right},
			[]Direction{Down, Down, Left},
			[]Direction{Left, Left, Down},
			[]Direction{Left, Left, Up},
			[]Direction{Up, Up, Left},
		}

		paths := []Path{}

		for _, jumps := range possibleJumps {
			if jump, err := s.Jump(jumps...); err == nil {
				paths = append(paths, Path{jump})
			}
		}

		return paths
	case Rook:
		return paths(
			s.Path(Up),
			s.Path(Right),
			s.Path(Down),
			s.Path(Left),
		)
	case Bishop:
		return paths(
			s.Path(UpRightDiagonal),
			s.Path(DownRightDiagonal),
			s.Path(DownLeftDiagonal),
			s.Path(UpLeftDiagonal),
		)
	case Queen:
		return paths(
			s.Path(Up),
			s.Path(UpRightDiagonal),
			s.Path(Right),
			s.Path(DownRightDiagonal),
			s.Path(Down),
			s.Path(DownLeftDiagonal),
			s.Path(Left),
			s.Path(UpLeftDiagonal),
		)
	case King:
		p := []Path{}

		for _, path := range paths(
			s.Path(Up),
			s.Path(UpRightDiagonal),
			s.Path(Right),
			s.Path(DownRightDiagonal),
			s.Path(Down),
			s.Path(DownLeftDiagonal),
			s.Path(Left),
			s.Path(UpLeftDiagonal),
		) {
			p = append(p, path[:1])
		}

		return p
	}

	return []Path{}
}

// Given a set of Path it simply strips out any that
// are empty
func paths(p ...Path) []Path {
	paths := []Path{}

	for _, path := range p {
		if len(path) > 0 {
			paths = append(paths, path)
		}
	}

	return paths
}

var materialByFile = map[File]Material{
	A: Rook,
	B: Knight,
	C: Bishop,
	D: Queen,
	E: King,
	F: Bishop,
	G: Knight,
	H: Rook,
}

func NewSquare(position Position) *Square {
	var color Color
	emptyColor := color

	var material Material
	emptyMaterial := material

	switch position.Rank {
	case 7, 8:
		color = Black
	case 1, 2:
		color = White
	}

	switch position.Rank {
	case 2, 7:
		material = Pawn
	case 1, 8:
		material = materialByFile[position.File]
	}

	if color != emptyColor && material != emptyMaterial {
		return &Square{Position: position, Piece: Piece{color, material}}
	}

	return &Square{Position: position}
}

func AllSquares() []*Square {
	var squares = []*Square{}

	for _, position := range allPositions {
		squares = append(squares, NewSquare(position))
	}

	return squares
}
