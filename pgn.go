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
	sc  scanner.Scanner
}

func NewPGNParser() *PGNParser {
	return &PGNParser{pgn: NewPGN()}
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

// Parses a string containing Portable Game Notation and returns a PGN struct
//
// Spec: https://www.chessclub.com/user/help/PGN-spec
func ParsePGN(str string) PGN {
	return NewPGNParser().ParseFromString(str)
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
				p.pgn.Movetext = append(p.pgn.Movetext, &Movetext{Number: uint8(moveNumber)})
			case white == "":
				p.scanMovetextForColor(&white)
				p.pgn.updateLastMovetext(func(lastMovetext *Movetext) *Movetext {
					lastMovetext.WhiteMove = AlgebraicNotation(white)
					return lastMovetext
				})
			case black == "":
				p.scanMovetextForColor(&black)
				p.pgn.updateLastMovetext(func(lastMovetext *Movetext) *Movetext {
					lastMovetext.BlackMove = AlgebraicNotation(black)
					return lastMovetext
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
