package pawn

import (
	"fmt"
	"regexp"
	"strings"
)

type AlgebraicNotation string

func (an AlgebraicNotation) Material() Material {
	firstCharacter := string(an[0])

	for material, name := range materialNames {
		if firstCharacter == name {
			return material
		}
	}

	return Pawn
}

func (an AlgebraicNotation) Takes() bool {
	return strings.Contains(string(an), "x")
}

func (an AlgebraicNotation) DestinationPosition() Position {
	destinationPositionStr := an.destinationPosition()
	return Position{
		File: File(string(destinationPositionStr[0])),
		Rank: Rank(rankFromByte(destinationPositionStr[1])),
	}
}

func (an AlgebraicNotation) destinationPosition() string {
	destinationPositionPattern := regexp.MustCompile("[a-h][1-8]")
	if matches := destinationPositionPattern.FindAllString(string(an), -1); len(matches) > 0 {
		return matches[len(matches)-1]
	}
	return ""
}

func (an AlgebraicNotation) OriginDisambiguated() bool {
	return an.OriginFile() != NilFile || an.OriginRank() != NilRank
}

func (an AlgebraicNotation) OriginRank() Rank {
	var rank Rank
	rankPattern := regexp.MustCompile("[1-8]")
	if matches := rankPattern.FindAllString(string(an), -1); len(matches) > 1 {
		rank = rankFromByte(matches[0][0])
	}

	return rank
}

func (an AlgebraicNotation) OriginFile() File {
	var file File

	filePattern := regexp.MustCompile("[a-h]")
	if matches := filePattern.FindAllString(string(an), -1); len(matches) > 1 {
		file = File(matches[0])
	}

	return file
}

// TODO
/*
For example, with knights on g1 and d2, either of which might move to f3, the move is specified as Ngf3 or Ndf3, as appropriate. With knights on g5 and g1, the moves are N5f3 or N1f3. As above, an "x" can be inserted to indicate a capture, for example: N5xf3. Another example: two rooks on d3 and h5, either one of which may move to d5. If the rook on d3 moves to d5, it is possible to disambiguate with either Rdd5 or R3d5, but the file takes precedence over the rank, so Rdd5 is correct. (And likewise if the move is a capture, Rdxd5 is correct.)
*/

func (an AlgebraicNotation) IsPromotion() bool {
	return an.PromotedTo() != Pawn
}

func (an AlgebraicNotation) PromotedTo() Material {
	if index := strings.Index(string(an), "="); index != -1 {
		materialCharacter := string(an[index+1])

		return AlgebraicNotation(materialCharacter).Material()
	}

	return Pawn
}

func (an AlgebraicNotation) IsCheck() bool {
	return strings.HasSuffix(string(an), "+")
}

func (an AlgebraicNotation) IsCheckMate() bool {
	return strings.HasSuffix(string(an), "#")
}

func (an AlgebraicNotation) IsCastle() bool {
	return an.IsCastleKingSide() || an.IsCastleQueenSide()
}

func (an AlgebraicNotation) IsCastleKingSide() bool {
	return an == "O-O"
}

func (an AlgebraicNotation) IsCastleQueenSide() bool {
	return an == "O-O-O"
}

type AlgebraiclyNotated interface {
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
	if p == NoPiece {
		return " "
	} else {
		return materialFigurines[p.Color][p.Material]
	}
}

func (p Position) AN() string {
	return fmt.Sprintf("%s%d", p.File, p.Rank)
}

func (s Square) AN() string {
	return s.Piece.AN() + s.Position.AN()
}
