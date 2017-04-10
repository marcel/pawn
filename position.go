package pawn

import (
	"math"
	"strconv"
)

type File string

const (
	A File = "a"
	B      = "b"
	C      = "c"
	D      = "d"
	E      = "e"
	F      = "f"
	G      = "g"
	H      = "h"
)

var allFiles = [...]File{
	A, B, C, D, E, F, G, H,
}

// TODO Just make File a uint8 like Rank
var allFilesMap = map[File]int{
	A: 0, B: 1, C: 2, D: 3, E: 4, F: 5, G: 6, H: 7,
}

func (f File) index() int {
	i, _ := allFilesMap[f]
	return i
}

type Rank uint8

var allRanks = [...]Rank{
	1, 2, 3, 4, 5, 6, 7, 8,
}

func rankFromByte(b byte) Rank {
	rankInt, _ := strconv.ParseUint(string(b), 10, 8)
	return Rank(rankInt)
}

type Position struct {
	File // vertical columns a through h from queenside to kingside
	Rank // horizontal rows 1 to 8 from White's side of the board
}

var allPositions []Position

func init() {
	for _, file := range allFiles {
		for _, rank := range allRanks {
			allPositions = append(
				allPositions,
				Position{
					File: file,
					Rank: rank,
				})
		}
	}
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
	UpLeftDiagonal
	UpRightDiagonal
	DownLeftDiagonal
	DownRightDiagonal
)

type Path []Position

func (p Position) Jump(directions ...Direction) (Position, error) {
	if len(directions) == 0 {
		return p, nil
	} else {
		var newPosition Position
		empty := newPosition

		switch directions[0] {
		case Up:
			if index := int(p.Rank); index < len(allRanks) {
				newPosition = Position{p.File, allRanks[index]}
			}
		case Right:
			if index := p.File.index() + 1; index < len(allFiles) {
				newPosition = Position{allFiles[index], p.Rank}
			}
		case Down:
			if index := int(p.Rank) - 2; index >= 0 {
				newPosition = Position{p.File, allRanks[index]}
			}
		case Left:
			if index := p.File.index() - 1; index >= 0 {
				newPosition = Position{allFiles[index], p.Rank}
			}
		}

		if newPosition == empty {
			return p, ErrorInvalidPosition
		} else {
			return newPosition.Jump(directions[1:]...)
		}
	}
}

func (p Position) Path(direction Direction) Path {
	switch direction {
	case Up:
		ranks := allRanks[p.Rank:]
		return p.pathByRanks(ranks)
	case Down:
		ranks := rankSlice(allRanks[:p.Rank-1]).reverse()
		return p.pathByRanks(ranks)
	case Left:
		files := allFiles[:p.File.index()]
		return p.pathByFiles(files)
	case Right:
		files := allFiles[p.File.index()+1:]
		return p.pathByFiles(files)
	case UpLeftDiagonal:
		files := fileSlice(allFiles[:p.File.index()]).reverse()
		ranks := allRanks[p.Rank:]
		return p.pathByFilesAndRanks(files, ranks)
	case UpRightDiagonal:
		files := allFiles[p.File.index()+1:]
		ranks := allRanks[p.Rank:]
		return p.pathByFilesAndRanks(files, ranks)
	case DownLeftDiagonal:
		files := fileSlice(allFiles[:p.File.index()]).reverse()
		ranks := rankSlice(allRanks[:p.Rank-1]).reverse()
		return p.pathByFilesAndRanks(files, ranks)
	case DownRightDiagonal:
		files := allFiles[p.File.index()+1:]
		ranks := rankSlice(allRanks[:p.Rank-1]).reverse()
		return p.pathByFilesAndRanks(files, ranks)
	}

	return Path{}
}

func (p Position) pathByRanks(ranks rankSlice) Path {
	path := Path{}

	for _, rank := range ranks {
		path = append(path, Position{p.File, rank})
	}

	return path
}

func (p Position) pathByFiles(files fileSlice) Path {
	path := Path{}
	for _, file := range files {
		path = append(path, Position{file, p.Rank})
	}

	return path
}

// Yak shaving so I can reverse a []File
type fileSlice []File

func (fs fileSlice) reverse() fileSlice {
	reversed := make(fileSlice, len(fs))
	copy(reversed, fs)

	for lhs, rhs := 0, len(reversed)-1; lhs < rhs; lhs, rhs = lhs+1, rhs-1 {
		reversed[lhs], reversed[rhs] = reversed[rhs], reversed[lhs]
	}
	return reversed
}

// Yak shaving so I can reverse a []Rank
type rankSlice []Rank

func (rs rankSlice) reverse() rankSlice {
	reversed := make(rankSlice, len(rs))
	copy(reversed, rs)

	for lhs, rhs := 0, len(reversed)-1; lhs < rhs; lhs, rhs = lhs+1, rhs-1 {
		reversed[lhs], reversed[rhs] = reversed[rhs], reversed[lhs]
	}
	return reversed
}

func (p Position) pathByFilesAndRanks(files fileSlice, ranks rankSlice) Path {
	size := int(math.Min(float64(len(files)), float64(len(ranks))))
	path := Path{}

	for i := 0; i < size; i++ {
		path = append(path, Position{files[i], ranks[i]})
	}

	return path
}

var (
	A1 = Position{A, 1}
	A2 = Position{A, 2}
	A3 = Position{A, 3}
	A4 = Position{A, 4}
	A5 = Position{A, 5}
	A6 = Position{A, 6}
	A7 = Position{A, 7}
	A8 = Position{A, 8}
	B1 = Position{B, 1}
	B2 = Position{B, 2}
	B3 = Position{B, 3}
	B4 = Position{B, 4}
	B5 = Position{B, 5}
	B6 = Position{B, 6}
	B7 = Position{B, 7}
	B8 = Position{B, 8}
	C1 = Position{C, 1}
	C2 = Position{C, 2}
	C3 = Position{C, 3}
	C4 = Position{C, 4}
	C5 = Position{C, 5}
	C6 = Position{C, 6}
	C7 = Position{C, 7}
	C8 = Position{C, 8}
	D1 = Position{D, 1}
	D2 = Position{D, 2}
	D3 = Position{D, 3}
	D4 = Position{D, 4}
	D5 = Position{D, 5}
	D6 = Position{D, 6}
	D7 = Position{D, 7}
	D8 = Position{D, 8}
	E1 = Position{E, 1}
	E2 = Position{E, 2}
	E3 = Position{E, 3}
	E4 = Position{E, 4}
	E5 = Position{E, 5}
	E6 = Position{E, 6}
	E7 = Position{E, 7}
	E8 = Position{E, 8}
	F1 = Position{F, 1}
	F2 = Position{F, 2}
	F3 = Position{F, 3}
	F4 = Position{F, 4}
	F5 = Position{F, 5}
	F6 = Position{F, 6}
	F7 = Position{F, 7}
	F8 = Position{F, 8}
	G1 = Position{G, 1}
	G2 = Position{G, 2}
	G3 = Position{G, 3}
	G4 = Position{G, 4}
	G5 = Position{G, 5}
	G6 = Position{G, 6}
	G7 = Position{G, 7}
	G8 = Position{G, 8}
	H1 = Position{H, 1}
	H2 = Position{H, 2}
	H3 = Position{H, 3}
	H4 = Position{H, 4}
	H5 = Position{H, 5}
	H6 = Position{H, 6}
	H7 = Position{H, 7}
	H8 = Position{H, 8}
)
