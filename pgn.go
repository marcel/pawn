package pawn

import (
	"strings"
	"text/scanner"
)

type PGN struct {
	Tags map[string]string
}

func NewPGN() PGN {
	return PGN{
		Tags: map[string]string{},
	}
}

type PGNParser struct {
	pgn PGN
}

func NewPGNParser() PGNParser {
	return PGNParser{pgn: NewPGN()}
}

func (p PGNParser) ParseFromString(str string) PGN {
	reader := strings.NewReader(str)
	scanner := scanner.Scanner{}
	scanner.Init(reader)

	return p.ParseFromScanner(&scanner)
}

func (p PGNParser) ParseFromScanner(sc *scanner.Scanner) PGN {
	ParseTags(sc, &p.pgn)
	ParseMoves(sc, &p.pgn)

	return p.pgn
}

func ParsePGN(str string) PGN {
	return NewPGNParser().ParseFromString(str)
}

func ParseTags(sc *scanner.Scanner, pgn *PGN) {
	scan := sc.Peek()

	for scan != scanner.EOF {
		switch scan {
		case '[', ']', '\n', '\r':
			scan = sc.Next()
		case '1': // First move; all tags have been read
			return
		default:
			sc.Scan()
			tag := sc.TokenText()
			sc.Scan()
			value := sc.TokenText()

			pgn.Tags[tag] = strings.Trim(value, "\"")
		}
		scan = sc.Peek()
	}
}

func ParseMoves(sc *scanner.Scanner, pgn *PGN) {
	sc.Mode = scanner.ScanIdents | scanner.ScanChars | scanner.ScanInts | scanner.ScanStrings
}
