package pawn

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePGN(t *testing.T) {
	assert := assert.New(t)

	pgnStrings := []string{
		win, draw, finalMoveByWhite, movetextWithRAV, embeddedComments, checkMate,
	}

	for _, pgnString := range pgnStrings {
		pgn := ParsePGN(pgnString)

		assert.Equal(len(pgn.Tags), strings.Count(pgnString, "["))
		assert.True(len(pgn.Movetext.Moves) > 0)
		assert.True(pgn.Outcome != "")

		trimmedPGNString := strings.TrimRight(pgnString, " \n\r")

		assert.True(
			strings.HasSuffix(trimmedPGNString, string(pgn.Outcome)),
			fmt.Sprintf(
				"'%s' does not contain '%s'",
				trimmedPGNString[len(trimmedPGNString)-10:],
				string(pgn.Outcome),
			),
		)
	}
}

func TestParseAll(t *testing.T) {
	file, err := os.Open("Carlsen.pgn")
	if err != nil {
		assert.Fail(t, "Could not open file")
	}

	parser := NewPGNParserFromReader(file)

	pgns := parser.parseAll()

	invalidPGNs := 0

	for _, pgn := range pgns {
		for _, moveText := range pgn.Movetext.Moves {
			if moveText.Number == 0 {
				invalidPGNs++
				fmt.Println(pgn)
				break
			}
		}
	}

	assert.Equal(t, invalidPGNs, 0)
}

func TestMultipleEntries(t *testing.T) {
	parser := NewPGNParserFromReader(strings.NewReader(multipleEntries))

	pgns := parser.parseAll()

	assert.Equal(t, len(pgns), 2)
}

func TestPGNString(t *testing.T) {
	pgn := ParsePGN(win)

	assert.True(t, strings.Contains(win, pgn.Tags.String()[:10]))
	assert.True(t, strings.Contains(win, pgn.Movetext.String()[:10]))

	assert.Equal(t, strings.Count(win, "["), strings.Count(pgn.String(), "["))
}

var win = `
[Event "WCh 2013"]
[Site "Chennai IND"]
[Date "2013.11.21"]
[Round "9"]
[White "Anand,V"]
[Black "Carlsen,M"]
[Result "0-1"]
[WhiteElo "2775"]
[BlackElo "2870"]
[ECO "E25"]

1.d4 Nf6 2.c4 e6 3.Nc3 Bb4 4.f3 d5 5.a3 Bxc3+ 6.bxc3 c5 7.cxd5 exd5 8.e3 c4
9.Ne2 Nc6 10.g4 O-O 11.Bg2 Na5 12.O-O Nb3 13.Ra2 b5 14.Ng3 a5 15.g5 Ne8 16.e4 Nxc1
17.Qxc1 Ra6 18.e5 Nc7 19.f4 b4 20.axb4 axb4 21.Rxa6 Nxa6 22.f5 b3 23.Qf4 Nc7
24.f6 g6 25.Qh4 Ne8 26.Qh6 b2 27.Rf4 b1=Q+ 28.Nf1 Qe1  0-1
`
var draw = `
[Event "WCh 2013"]
[Site "Chennai IND"]
[Date "2013.11.22"]
[Round "10"]
[White "Carlsen,M"]
[Black "Anand,V"]
[Result "1/2-1/2"]
[WhiteElo "2870"]
[BlackElo "2775"]
[ECO "B51"]

1.e4 c5 2.Nf3 d6 3.Bb5+ Nd7 4.d4 cxd4 5.Qxd4 a6 6.Bxd7+ Bxd7 7.c4 Nf6 8.Bg5 e6
9.Nc3 Be7 10.O-O Bc6 11.Qd3 O-O 12.Nd4 Rc8 13.b3 Qc7 14.Nxc6 Qxc6 15.Rac1 h6
16.Be3 Nd7 17.Bd4 Rfd8 18.h3 Qc7 19.Rfd1 Qa5 20.Qd2 Kf8 21.Qb2 Kg8 22.a4 Qh5
23.Ne2 Bf6 24.Rc3 Bxd4 25.Rxd4 Qe5 26.Qd2 Nf6 27.Re3 Rd7 28.a5 Qg5 29.e5 Ne8
30.exd6 Rc6 31.f4 Qd8 32.Red3 Rcxd6 33.Rxd6 Rxd6 34.Rxd6 Qxd6 35.Qxd6 Nxd6
36.Kf2 Kf8 37.Ke3 Ke7 38.Kd4 Kd7 39.Kc5 Kc7 40.Nc3 Nf5 41.Ne4 Ne3 42.g3 f5
43.Nd6 g5 44.Ne8+ Kd7 45.Nf6+ Ke7 46.Ng8+ Kf8 47.Nxh6 gxf4 48.gxf4 Kg7 49.Nxf5+ exf5
50.Kb6 Ng2 51.Kxb7 Nxf4 52.Kxa6 Ne6 53.Kb6 f4 54.a6 f3 55.a7 f2 56.a8=Q f1=Q
57.Qd5 Qe1 58.Qd6 Qe3+ 59.Ka6 Nc5+ 60.Kb5 Nxb3 61.Qc7+ Kh6 62.Qb6+ Qxb6+ 
63.Kxb6 Kh5 64.h4 Kxh4 65.c5 Nxc5  1/2-1/2
`

var finalMoveByWhite = `
[Event "World Blitz 2016"]
[Site "Doha QAT"]
[Date "2016.12.30"]
[Round "21.1"]
[White "Leko,P"]
[Black "Carlsen,M"]
[Result "1/2-1/2"]
[WhiteElo "2693"]
[BlackElo "2840"]
[ECO "C50"]

1.e4 e5 2.Nf3 Nc6 3.Bc4 Bc5 4.O-O Nf6 5.d3 a6 6.c3 d6 7.Re1 Ba7 8.Bb3 O-O
9.h3 Ne7 10.Nbd2 Ng6 11.Nf1 c6 12.Ng3 d5 13.exd5 Nxd5 14.d4 exd4 15.Nxd4 h6
16.Bc2 Re8 17.Rxe8+ Qxe8 18.Bb3 Qe5 19.Kh2 Ngf4 20.Bd2 Bd7 21.Qf3 Ne6 22.Nxe6 Bxe6
23.Re1 Qd6 24.Kg1 Nf6 25.Bf4 Qe7 26.Bxe6 fxe6 27.Be3 Bxe3 28.Qxe3 Re8 29.Qe5 Qd7
30.Ne4 Nxe4 31.Qxe4 Qd5 32.a3 Kf7 33.Qf4+ Qf5 34.Qc7+ Re7 35.Qd6 Qc2 36.Qf4+ Kg8
37.Qb4 Rf7 38.Qb6 Qd2 39.Rxe6 Qc1+ 40.Kh2 Qf4+ 41.Kg1 Qc1+ 42.Kh2  1/2-1/2
`

var movetextWithRAV = `
[Event "World Blitz 2016"]
[Site "Doha QAT"]
[Date "2016.12.30"]
[Round "21.1"]
[White "Leko,P"]
[Black "Carlsen,M"]
[Result "1/2-1/2"]
[WhiteElo "2693"]
[BlackElo "2840"]
[ECO "C50"]

1.e4 b5 (1 ... e6 $1 {is of course better.}) 2.Nf3 Nc6 3.Bc4 Bc5 4.O-O Nf6  1/2-1/2
`

var embeddedComments = `
[Event "Baden Baden"]
[Site "Baden Baden DEU"]
[Date "1925.04.17"]
[EventDate "1925.04.16"]
[Round "2"]
[Result "0-1"]
[White "Jan Willem te Kolste"]
[Black "Carlos Torre Repetto"]
[ECO "C12"]
[WhiteElo "?"]
[BlackElo "?"]
[PlyCount "54"]

1. e4 {Notes by Lasker} e6 2. d4 d5 3. Nc3 Nf6 4. Bg5 Bb4
5. Nge2 dxe4 6. a3 Be7 7. Bxf6 gxf6 8. Nxe4 b6 9. N2c3 Bb7
10. Qf3 {Though check with the Knight is now threatened,
White's intention is not aggresive but merely the protection
of the g-pawn.} c6 11. O-O-O f5 12. Ng3 Nd7 13. Bc4 {But here 
White begins to be aggresive and therby he puts himself in the
wrong for his aim is not even remotely warranted by his
advantage. Again he should strive to safeguard his g-pawn by
13.Nh5, to be followed by Rg1, and afterwards possibly by h3
and g4 with some little pressure on the center and on the
h-pawn, proportionate to his advantage in mobility. As he
plays White begins to slide first very slowly; afterwards more 
quickly.} Qc7 14. Rhe1 Nf6 15. Qe2 O-O-O {Black has defended
conscientiously. If White now continues 16.Nxf5 Qf4+ 17.Ne3
Rxd4 18.Ba6 he can still obtain equality. But here is the
parting of the ways.} 16. Bxe6+ {White fancies that he holds
an advantage and attempts to win. The punishment is
immediate.} fxe6 17. Qxe6+ Rd7 18. Nxf5 Bd8 19. Ne4 Nxe4 
20. Rxe4 Kb8 21. g3 Bc8 22. Qh6 Rf7 23. Ne3 Rxf2 24. c3 Rg8
25. d5 Bg5 26. Qxc6 Qxc6 27. dxc6 Bf5 0-1
`
var checkMate = `
[Event "New Orleans, USA"]
[Site "New Orleans, USA"]
[Date "1927.??.??"]
[EventDate "?"] 
[Round "?"] 
[Result "0-1"]
[White "Dupre"]
[Black "Carlos Torre Repetto"]
[ECO "C41"]
[WhiteElo "?"] 
[BlackElo "?"] 
[PlyCount "20"]

1. e4 e5 2. Nf3 d6 3. d4 f5 4. Bc4 exd4 5. exf5 Qe7+ 6. Kd2 g6
7. Re1 Bh6+ 8. Kd3 Bxf5+ 9. Kxd4 Bg7+ 10. Kd5 c6# 0-1
`
var multipleEntries = `
[Event "World Blitz 2016"]
[Site "Doha QAT"]
[Date "2016.12.30"]
[Round "19.1"]
[White "Carlsen,M"]
[Black "Onischuk,V"]
[Result "1-0"]
[WhiteElo "2840"]
[BlackElo "2601"]
[ECO "B07"]

1.e4 d6 2.d4 Nf6 3.Nc3 g6 4.Be3 a6 5.h3 Bg7 6.f4 b5 7.e5 Nfd7 8.Nf3 Nb6 9.Bd3 Bb7
10.O-O e6 11.Be4 Bxe4 12.Nxe4 Nc6 13.Qe2 Qd7 14.Rad1 Nd5 15.c3 f5 16.exf6 Nxf6
17.Nxf6+ Bxf6 18.d5 Ne7 19.f5 exf5 20.Bh6 O-O-O 21.a4 Kb7 22.axb5 axb5 23.Bg5 Bxg5
24.Nxg5 Ra8 25.Ne6 Nc8 26.b3 Nb6 27.c4 Rhe8 28.c5 Nxd5 29.Rxd5 Qxe6 30.Qxb5+ Kc8
31.cxd6 cxd6 32.Qc6+  1-0

[Event "World Blitz 2016"]
[Site "Doha QAT"]
[Date "2016.12.30"]
[Round "20.1"]
[White "Carlsen,M"]
[Black "Anand,V"]
[Result "1-0"]
[WhiteElo "2840"]
[BlackElo "2779"]
[ECO "A45"]

1.d4 Nf6 2.Bf4 d5 3.e3 c5 4.c3 Nc6 5.Nd2 e6 6.Ngf3 Bd6 7.Bg3 O-O 8.Bb5 a6
9.Bxc6 bxc6 10.Qa4 Rb8 11.Qa3 Bxg3 12.hxg3 cxd4 13.cxd4 a5 14.O-O Qb6 15.b3 Ba6
16.Rfc1 Nd7 17.Qd6 Qa7 18.Rxc6 Bb5 19.Rc7 Rb7 20.Rac1 a4 21.Rxb7 Qxb7 22.Rc7 Qb8
23.Rxd7 Bxd7 24.Qxd7  1-0

`
