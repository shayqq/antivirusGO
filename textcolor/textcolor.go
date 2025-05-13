package textcolor

import "github.com/fatih/color"

func RedErrorText(text string) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(text)
}

func GreenSuccessText(text string) string {
	green := color.New(color.FgGreen).SprintFunc()
	return green(text)
}
