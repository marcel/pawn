package pawn

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PositionTestSuite struct {
	suite.Suite
}

func TestPositionTestSuite(t *testing.T) {
	suite.Run(t, new(PositionTestSuite))
}

func (s *PositionTestSuite) TestPath() {
	start := E5

	s.Equal(Path{E6, E7, E8}, start.Path(Up))
	s.Equal(Path{E4, E3, E2, E1}, start.Path(Down))
	s.Equal(Path{A5, B5, C5, D5}, start.Path(Left))
	s.Equal(Path{F5, G5, H5}, start.Path(Right))
	s.Equal(Path{D6, C7, B8}, start.Path(UpLeftDiagonal))
	s.Equal(Path{F6, G7, H8}, start.Path(UpRightDiagonal))
	s.Equal(Path{D4, C3, B2, A1}, start.Path(DownLeftDiagonal))
	s.Equal(Path{F4, G3, H2}, start.Path(DownRightDiagonal))

	var invalidDirection Direction = 404
	s.Equal(Path{}, start.Path(invalidDirection))
}

func (s *PositionTestSuite) TestJump() {
	up1, ok := E6.Jump(Up)
	s.Equal(E7, up1)
	s.Nil(ok)

	up2, ok := E6.Jump(Up, Up)
	s.Equal(E8, up2)
	s.Nil(ok)

	_, notok := E6.Jump(Up, Up, Up)
	s.Equal(ErrorInvalidPosition, notok)

	right1, ok := F6.Jump(Right)
	s.Equal(G6, right1)
	s.Nil(ok)

	right2, ok := F6.Jump(Right, Right)
	s.Equal(H6, right2)
	s.Nil(ok)

	_, notok2 := F6.Jump(Right, Right, Right)
	s.Equal(ErrorInvalidPosition, notok2)

	down1, ok := E3.Jump(Down)
	s.Equal(E2, down1)
	s.Nil(ok)

	down2, ok := E3.Jump(Down, Down)
	s.Equal(E1, down2)
	s.Nil(ok)

	_, notok3 := E3.Jump(Down, Down, Down)
	s.Equal(ErrorInvalidPosition, notok3)

	left1, ok := C6.Jump(Left)
	s.Equal(B6, left1)
	s.Nil(ok)

	left2, ok := C6.Jump(Left, Left)
	s.Equal(A6, left2)
	s.Nil(ok)

	_, notok4 := C6.Jump(Left, Left, Left)
	s.Equal(ErrorInvalidPosition, notok4)

	knightJump, ok := E5.Jump(Up, Up, Right)
	s.Equal(F7, knightJump)

	_, notok5 := A1.Jump(Up, Up, Left)
	s.Equal(ErrorInvalidPosition, notok5)
}
