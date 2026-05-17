package tui

type Theme struct {
	borderColour string
}

func currentTheme() Theme {
	return defaultTheme()
}

func defaultTheme() Theme {
	return Theme{
		borderColour: "#27C5F5",
	}
}
