package functions

import "github.com/fatih/color"

// ErrorColor returns a color for printing
// errors
func ErrorColor() *color.Color {
	return color.New(color.FgHiRed).Add(color.Bold)
}
