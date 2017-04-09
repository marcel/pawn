package main

import (
	"fmt"
	"os"

	"github.com/marcel/pawn"
)

func main() {
	pgns := pawn.ParseAllPGNFromFilePath(os.Args[1])

	for _, pgn := range pgns {
		fmt.Println(pgn)
	}
}
