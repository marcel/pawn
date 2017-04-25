package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"github.com/marcel/pawn"
	"github.com/olekukonko/tablewriter"
)

// TODO Integrate Flags
var gamePlayer = new(GamePlayer)
var gameMenu = new(GameMenu)
var initialized bool

func (gp *GamePlayer) loadCurrentSelection(gameMenu *GameMenu) bool {
	if gameMenu.currentGame != gameMenu.currentSelection || !initialized {
		selectedPGN := gameMenu.pgns[gameMenu.currentSelection]
		gp.InitWithPGN(&selectedPGN)
		gameMenu.currentGame = gameMenu.currentSelection

		return true
	}

	return false
}

func initializeGamePlayer() {
	pgnFile := os.Args[1]

	file, err := os.Open(pgnFile)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var pgnReader io.Reader

	if strings.HasSuffix(pgnFile, ".gz") {
		pgnReader, err = gzip.NewReader(file)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		pgnReader = file
	}

	gameMenu.pgns = pawn.NewPGNParserFromReader(pgnReader).ParseAll()

	gamePlayer.loadCurrentSelection(gameMenu)

	initialized = true
}

func main() {
	initializeGamePlayer()

	g, _ := gocui.NewGui(gocui.Output256)
	defer g.Close()

	commandHelp := CommandHelp{[]string{
		"↑     Previous Game",
		"↓     Next Game",
		"Enter Select Game",
		"→     Next Move",
		"←     Previous Move",
	}}

	g.SetManager(gameMenu, commandHelp, gamePlayer)

	initKeybindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

type CommandHelp struct {
	commands []string
}

type viewDimensions struct {
	x0 int
	x1 int
	y0 int
	y1 int
}

type viewInitializer func(*gocui.View)

func initializeView(g *gocui.Gui, name string, dimensions viewDimensions, initializer viewInitializer) {
	v, err := g.SetView(name, dimensions.x0, dimensions.y0, dimensions.x1, dimensions.y1)
	if err != nil {
		initializer(v)
	}
}

func (ch CommandHelp) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	dimensions := viewDimensions{
		x0: maxX - 25,
		y0: 0,
		x1: maxX - 1,
		y1: len(ch.commands) + 1,
	}

	initializeView(g, "help", dimensions,
		func(v *gocui.View) {
			for _, command := range ch.commands {
				v.Title = "Controls"
				fmt.Fprintln(v, command)
			}
		},
	)

	return nil
}

func initKeybindings(g *gocui.Gui) error {
	g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			gameMenu.moveSideBarCursor(g, +1)
			return nil
		},
	)

	g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			gameMenu.moveSideBarCursor(g, -1)
			return nil
		},
	)

	g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {

			gamePlayer.loadCurrentSelection(gameMenu)
			return nil
		},
	)

	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		},
	)

	g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			gamePlayer.playPreviousMove()

			return nil
		},
	)
	g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			gamePlayer.playNextMove()

			return nil
		},
	)

	return nil
}

type GameMenu struct {
	pgns             []pawn.PGN
	currentGame      int
	currentSelection int
}

func (gm GameMenu) Layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 25, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for index, pgn := range gameMenu.pgns {
			matchup := pgn.MatchUp()
			fmt.Fprintf(v, "%3d. %s\n", index+1, matchup)
		}
	}

	return nil
}

func (gm *GameMenu) moveSideBarCursor(g *gocui.Gui, increment int) {
	view, _ := g.View("side")
	cx, cy := view.Cursor()
	view.SetCursor(cx, cy+increment)
	// TODO Only update the origin once the cursor gets about half way down the
	// page
	// ox, oy := view.Origin()
	// view.SetOrigin(ox, oy+increment)

	gm.currentSelection += increment
}

type GamePlayer struct {
	board       *pawn.Board
	pgn         *pawn.PGN
	currentTurn int
}

func (gp *GamePlayer) Init() *GamePlayer {
	gp.board = pawn.NewBoard()
	gp.currentTurn = 0
	return gp
}

func (gp *GamePlayer) InitWithPGN(pgn *pawn.PGN) *GamePlayer {
	gp.Init()
	gp.pgn = pgn

	return gp
}

func (gp GamePlayer) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	name := fmt.Sprintf("board-%d-%d", gameMenu.currentGame, gp.currentTurn)
	v, err := g.SetView(name, maxX/2-17, 3, maxX/2+17, 21)
	v.Frame = false
	v.Editable = false

	if err != nil {

		summary := color.New(color.FgWhite, color.Bold, color.Underline).Sprintf("%s vs %s", gp.pgn.Tags["White"], gp.pgn.Tags["Black"])
		siteDate := fmt.Sprintf("%s: %s", gp.pgn.Tags["Site"], gp.pgn.Tags["Date"])

		matchSummaryViewName := fmt.Sprintf("summary-%d", gameMenu.currentGame)
		x0, _, x1, _, _ := g.ViewPosition(name)
		if matchSummaryView, err := g.SetView(matchSummaryViewName, x0-5, 0, x1+5, 3); err != nil {
			matchSummaryView.Frame = false
			fmt.Fprintln(matchSummaryView, tablewriter.Pad(summary, " ", x1-x0+10))
			fmt.Fprintf(matchSummaryView, tablewriter.Pad(siteDate, " ", x1-x0+10))
		}

		table := tablewriter.NewWriter(v)
		table.SetRowLine(true)

		table.SetCenterSeparator("+")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("—")

		table.SetAlignment(tablewriter.ALIGN_CENTER)

		if gp.currentTurn > 0 {
			an := gp.pgn.Turns()[gp.currentTurn-1]
			gp.board.MoveFromAlgebraic(an)
		}

		for _, row := range gp.board.Rows() {
			rowStr := []string{}
			for _, square := range row {
				rowStr = append(rowStr, square.Piece.FAN())
			}

			table.Append(rowStr)
		}
		table.Render()
	}

	g.SetCurrentView(name)
	g.SetViewOnTop(name)

	movesViewName := fmt.Sprintf("moves-%d-%d", gameMenu.currentSelection, gp.currentTurn)
	x0, _, x1, y1, _ := g.ViewPosition(name)
	// TODO Make these dimensions both relative to the board
	// but also a % of the space between the board and the bottom of
	// the window
	movesView, mvE := g.SetView(movesViewName, x0-5, y1+1, x1+5, y1+7)
	movesView.Wrap = true
	movesView.Autoscroll = true
	movesView.Title = "Moves"

	if mvE != nil {
		moves := []string{}
		for index, an := range gp.pgn.Turns()[:gp.currentTurn] {
			if index%2 == 0 {
				moveNumber := index/2 + 1
				moveStr := fmt.Sprintf("\033[3%d;%dm%d.\033[0m", 7, 1, moveNumber)
				moves = append(moves, moveStr)
			}
			if index == gp.currentTurn-1 {
				moveStr := fmt.Sprintf("\033[3%d;%dm%s\033[0m", 7, 7, an)
				moves = append(moves, moveStr)
			} else {
				moves = append(moves, string(an))
			}
		}
		if len(moves) == gp.currentTurn+len(gp.pgn.Moves) {
			moves = append(moves, string(gp.pgn.Outcome))
		}
		fmt.Fprintln(movesView, strings.Join(moves, " "))
	}

	g.SetViewOnTop(movesViewName)

	return nil
}

func (gp *GamePlayer) playNextMove() {
	if gp.currentTurn < len(gp.pgn.Turns())-1 {
		gp.currentTurn++
	}
}

func (gp *GamePlayer) playPreviousMove() {
	if gp.currentTurn > 0 {
		gp.currentTurn--
	}
}
