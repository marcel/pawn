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
	destinationPositionPattern, _ := regexp.Compile("[a-h][1-8]")
	if matches := destinationPositionPattern.FindAllString(string(an), -1); len(matches) > 0 {
		return matches[len(matches)-1]
	}
	return ""
}

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
	return materialFigurines[p.Color][p.Material]
}

func (p Position) AN() string {
	return fmt.Sprintf("%s%d", p.File, p.Rank)
}

func (s Square) AN() string {
	return s.Piece.AN() + s.Position.AN()
}
