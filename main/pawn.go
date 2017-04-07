package main

import (
	"fmt"

	"github.com/marcel/pawn"
)

func main() {
	fmt.Printf("%+v\n", pawn.NewBoard().Squares)
}
