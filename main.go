package main

import (
	"fmt"

	"github.com/grig-iv/anki-card-creator/ld"
	// "github.com/charmbracelet/bubbletea"
)

func main() {
	page, err := ld.ParseUrl("https://www.ldoceonline.com/dictionary/shackle")
	if err != nil {
		panic(err)
	}

	fmt.Println(page)
}
