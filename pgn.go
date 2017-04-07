package pawn

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type Outcome string

const (
	WhiteWin Outcome = "1-0"
	BlackWin         = "0-1"
	Draw             = "1/2-1/2"
)

type Tags map[string]string

func (t Tags) String() string {
	str := ""

	for tag, value := range t {
		str += fmt.Sprintf("[%s \"%s\"]\n", tag, value)
	}

	return str
}

type PGN struct {
	Tags
	Movetext []*Movetext
	Outcome
}

func (p PGN) String() string {
	moveText := []string{}

	for _, mv := range p.Movetext {
		moveText = append(moveText, mv.String())
	}

	return fmt.Sprintf("%s\n%s %s\n", p.Tags, strings.Join(moveText, " "), p.Outcome)
}

func (p PGN) updateLastMovetext(update func(*Movetext) *Movetext) {
	lastIndex := len(p.Movetext) - 1
	lastMovetext := p.Movetext[lastIndex]
	p.Movetext[lastIndex] = update(lastMovetext)
}

type Movetext struct {
	Number    uint8
	WhiteMove AlgebraicNotation
	BlackMove AlgebraicNotation
}

func (m Movetext) String() string {
	return fmt.Sprintf("%d.%s %s", m.Number, m.WhiteMove, m.BlackMove)
}

func NewPGN() PGN {
	return PGN{
		Tags:     map[string]string{},
		Movetext: []*Movetext{},
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

// Parses a string containing Portable Game Notation and returns a PGN struct
//
// Spec: https://www.chessclub.com/user/help/PGN-spec
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

	var num, white, black string

	scan := sc.Peek()

	reset := func() {
		num, white, black = "", "", ""
	}

	reset()

	for scan != scanner.EOF {
		switch scan {
		case '{':
			// Scan past comments
			scanUntilPast('}', sc, &scan)
		case '(':
			// Scan past RAVs
			scanUntilPast(')', sc, &scan)
		case '#', '.', '+', '!', '?', '\n', '\r':
			scan = sc.Next()
			scan = sc.Peek()
		default:
			sc.Scan()

			if sc.TokenText() == "{" {
				scan = '{'
				continue
			} else if sc.TokenText() == "(" {
				scan = '('
				continue
			}

			switch {
			case num == "":
				num = sc.TokenText()

				scanForOutcome(&num, sc)

				if reachedOutcome(num) {
					pgn.Outcome = Outcome(num)
					return
				}

				moveNumber, _ := strconv.ParseUint(num, 10, 8)
				pgn.Movetext = append(pgn.Movetext, &Movetext{Number: uint8(moveNumber)})
			case white == "":
				scanMovetextForColor(&white, sc, pgn)
				pgn.updateLastMovetext(func(lastMovetext *Movetext) *Movetext {
					lastMovetext.WhiteMove = AlgebraicNotation(white)
					return lastMovetext
				})
			case black == "":
				scanMovetextForColor(&black, sc, pgn)
				pgn.updateLastMovetext(func(lastMovetext *Movetext) *Movetext {
					lastMovetext.BlackMove = AlgebraicNotation(black)
					return lastMovetext
				})

				reset()
			}

			scan = sc.Peek()
		}
	}
}

func scanUntilPast(r rune, sc *scanner.Scanner, scan *rune) {
	for *scan != r && *scan != scanner.EOF {
		*scan = sc.Next()
	}
}

func scanForOutcome(number *string, sc *scanner.Scanner) {
	for sc.Peek() == '-' {
		for i := 0; i < 2; i++ {
			sc.Scan()
			*number += sc.TokenText()
		}
	}

	for sc.Peek() == '/' {
		for i := 0; i < 6; i++ {
			sc.Scan()
			*number += sc.TokenText()
		}
	}
}

func scanMovetextForColor(color *string, sc *scanner.Scanner, pgn *PGN) {
	*color = sc.TokenText()

	scanForOutcome(color, sc)

	if reachedOutcome(*color) {
		pgn.Outcome = Outcome(*color)
		return
	}

	// Pawn promotion
	if sc.Peek() == '=' {
		for i := 0; i < 2; i++ {
			sc.Scan()
			*color += sc.TokenText()
		}
	}

	// Check or check mate
	if peek := sc.Peek(); peek == '+' || peek == '#' {
		sc.Scan()
		*color += sc.TokenText()
	}
}

// N.B. Using a redudant map for O(1) looks rather than O(n) array lookup
var outcomes = map[string]Outcome{
	string(WhiteWin): WhiteWin,
	string(BlackWin): BlackWin,
	string(Draw):     Draw,
}

func reachedOutcome(str string) bool {
	_, containsKey := outcomes[str]
	return containsKey
}
