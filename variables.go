package main

import (
	"github.com/bclicn/color"
)

// PathColored set the path color
func PathColored(path string) string { return color.Blue(path) }

// Message set the color and message for each symbol
var Message = map[byte]string{
	'A': color.Green("ADDED"),
	'D': color.Red("DELETED"),
	'M': color.Yellow("MODIFIED"),
	'R': color.LightGreen("RENAMED"),
	'U': color.Purple("UNMERGED"),
	'?': color.Purple("UNTRACKED"),
	'-': color.Green("CLEAN"),
}

// Location set the color of repository Index and Working Tree
var Location = [...]string{
	color.LightCyan("Index"),
	color.LightCyan("Working Tree"),
}

// ErrorSymbol set the color for error
var ErrorSymbol = color.Red("X")

// SkipDir set the color for skipped directories
var SkipDir = color.Cyan("Directory to SKIP:")
