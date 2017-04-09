package pawn

import (
	"fmt"
	"io"
	"os"
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
	Movetext
	Outcome
}

type Movetext struct {
	Moves []*MovetextMove
}

type MovetextMove struct {
	Number    uint8
	WhiteMove AlgebraicNotation
	BlackMove AlgebraicNotation
}

func (p PGN) updateLastMove(update func(*MovetextMove)) {
	lastIndex := len(p.Moves) - 1
	lastMove := p.Moves[lastIndex]
	update(lastMove)
}

func (p PGN) String() string {
	return fmt.Sprintf("%s\n%s %s\n", p.Tags, p.Movetext, p.Outcome)
}

func (m Movetext) String() string {
	moves := []string{}

	for _, move := range m.Moves {
		moves = append(moves, move.String())
	}

	return strings.Join(moves, " ")
}

func (m MovetextMove) String() string {
	return fmt.Sprintf("%d.%s %s", m.Number, m.WhiteMove, m.BlackMove)
}

func NewPGN() PGN {
	return PGN{
		Tags:     map[string]string{},
		Movetext: Movetext{Moves: []*MovetextMove{}},
	}
}

type PGNParser struct {
	pgn PGN
	sc  scanner.Scanner
}

func NewPGNParser() *PGNParser {
	return &PGNParser{pgn: NewPGN()}
}

func NewPGNParserFromReader(r io.Reader) *PGNParser {
	scanner := scanner.Scanner{}
	scanner.Init(r)

	return &PGNParser{pgn: NewPGN(), sc: scanner}
}

func (p *PGNParser) ParseFromString(str string) PGN {
	reader := strings.NewReader(str)
	scanner := scanner.Scanner{}
	scanner.Init(reader)

	p.sc = scanner

	return p.parse()
}

func (p *PGNParser) parse() PGN {
	p.parseTags()
	p.parseMoves()

	return p.pgn
}

func (p *PGNParser) hasNext() bool {
	for p.sc.Peek() == ' ' || p.sc.Peek() == '\n' {
		p.sc.Next()
	}

	return p.sc.Peek() != scanner.EOF
}

func (p *PGNParser) parseAll() []PGN {
	pgns := []PGN{}

	for p.hasNext() {
		pgns = append(pgns, p.parse())
		p.pgn = NewPGN()
	}

	return pgns
}

// Parses a string containing Portable Game Notation and returns a PGN struct
//
// Spec: https://www.chessclub.com/user/help/PGN-spec
func ParsePGN(str string) PGN {
	return NewPGNParser().ParseFromString(str)
}

func ParseAllPGNFromFilePath(str string) []PGN {
	file, _ := os.Open(str) // TODO error handling
	defer file.Close()

	return NewPGNParserFromReader(file).parseAll()
}

func (p *PGNParser) parseTags() {
	scan := p.sc.Peek()

	for scan != scanner.EOF {
		switch scan {
		case '[', ']', '\n', '\r':
			scan = p.sc.Next()
		case '1': // First move; all tags have been read
			return
		default:
			p.sc.Scan()
			tag := p.sc.TokenText()
			p.sc.Scan()
			value := p.sc.TokenText()

			p.pgn.Tags[tag] = strings.Trim(value, "\"")
		}
		scan = p.sc.Peek()
	}
}

func (p *PGNParser) parseMoves() {
	p.sc.Mode = scanner.ScanIdents | scanner.ScanChars | scanner.ScanInts | scanner.ScanStrings

	var num, white, black string

	scan := p.sc.Peek()

	reset := func() {
		num, white, black = "", "", ""
	}

	reset()

	for scan != scanner.EOF {
		switch scan {
		case '{':
			// Scan past comments
			p.scanUntilPast('}', &scan)
		case '(':
			// Scan past RAVs
			p.scanUntilPast(')', &scan)
		case '#', '.', '+', '!', '?', '\n', '\r':
			scan = p.sc.Next()
			scan = p.sc.Peek()
		default:
			p.sc.Scan()

			if p.sc.TokenText() == "{" {
				scan = '{'
				continue
			} else if p.sc.TokenText() == "(" {
				scan = '('
				continue
			}

			switch {
			case num == "":
				num = p.sc.TokenText()

				p.scanForOutcome(&num)

				if reachedOutcome(num) {
					p.pgn.Outcome = Outcome(num)
					return
				}

				moveNumber, _ := strconv.ParseUint(num, 10, 8)
				p.pgn.Moves = append(p.pgn.Moves, &MovetextMove{Number: uint8(moveNumber)})
			case white == "":
				p.scanMovetextForColor(&white)

				if p.pgn.Outcome != "" {
					return
				}

				p.pgn.updateLastMove(func(lastMove *MovetextMove) {
					lastMove.WhiteMove = AlgebraicNotation(white)
				})

			case black == "":
				p.scanMovetextForColor(&black)

				if p.pgn.Outcome != "" {
					return
				}

				p.pgn.updateLastMove(func(lastMove *MovetextMove) {
					lastMove.BlackMove = AlgebraicNotation(black)
				})

				reset()
			}

			scan = p.sc.Peek()
		}
	}
}

func (p *PGNParser) scanUntilPast(r rune, scan *rune) {
	for *scan != r && *scan != scanner.EOF {
		*scan = p.sc.Next()
	}
}

func (p *PGNParser) scanForOutcome(number *string) {
	for p.sc.Peek() == '-' {
		for i := 0; i < 2; i++ {
			p.sc.Scan()
			*number += p.sc.TokenText()
		}
	}

	for p.sc.Peek() == '/' {
		for i := 0; i < 6; i++ {
			p.sc.Scan()
			*number += p.sc.TokenText()
		}
	}
}

func (p *PGNParser) scanMovetextForColor(color *string) {
	*color = p.sc.TokenText()

	p.scanForOutcome(color)

	if reachedOutcome(*color) {
		p.pgn.Outcome = Outcome(*color)
		return
	}

	// Pawn promotion
	if p.sc.Peek() == '=' {
		for i := 0; i < 2; i++ {
			p.sc.Scan()
			*color += p.sc.TokenText()
		}
	}

	// Check or check mate
	if peek := p.sc.Peek(); peek == '+' || peek == '#' {
		p.sc.Scan()
		*color += p.sc.TokenText()
	}
}

// N.B. Using a redundant map for O(1) looks rather than O(n) array lookup
var outcomes = map[string]Outcome{
	string(WhiteWin): WhiteWin,
	string(BlackWin): BlackWin,
	string(Draw):     Draw,
}

func reachedOutcome(str string) bool {
	_, containsKey := outcomes[str]
	return containsKey
}
