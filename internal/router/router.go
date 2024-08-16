// Copyright (c) 2023,2024 by Marko Gaćeša

package router

import (
	"fmt"
	"gamatet/internal/config"
	"gamatet/logic/menu"
)

type Router map[string][]menu.Item

func NewRouter(config *config.Config) Router {
	x := 5
	b := true
	t := "Marko"
	e := "Marko"
	y := 42

	r := make(map[string][]menu.Item)
	r[""] = []menu.Item{
		menu.NewCommand("Start", "A description for the integer input", func() {
			fmt.Println("Hi!")
		}),
		menu.NewInteger(&x, 10, 20, "Level", "A description for the integer input"),
		menu.NewBool(&b, "Yes/No", "A description for the boolean input"),
		menu.NewText(&t, "Player 1", "A description for the text field", 20),
		menu.NewEnum(&e, []string{"Ogi", "Marko", "Ika"}, "Pick one", "A description for string the enum field"),
		menu.NewEnum(&y, []int{1, 7, 42, 66, 108}, "Pick a number", "A description for the int enum field"),
	}

	return r
}
