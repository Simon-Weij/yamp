package components

import (
	"fmt"
)

func HomeView(choices []string, cursor int) string {
	s := "Where would you like to go?\n\n"
	for i, choice := range choices {
		cur := " "
		if cursor == i {
			cur = ">"
		}
		s += fmt.Sprintf("%s %s\n", cur, choice)
	}

	s += "\nPress enter to select. Press q to quit.\n"

	return s
}
