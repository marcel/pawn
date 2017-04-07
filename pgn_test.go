package pawn

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestParsePGN(t *testing.T) {
	assert := assert.New(t)

	pgn := ParsePGN(draw)

	assert.Equal(len(pgn.Tags), strings.Count(draw, "["))
	assert.Equal(pgn.Tags["White"], "Carlsen,M")
	assert.Equal(pgn.Tags["Result"], "1/2-1/2")
}
