package main

import (
	"github.com/bclicn/color"
)

func PathColored(path string) string { return color.Blue(path) }

var Message = map[byte]string{
	'A': color.Green("ADDED"),
	'D': color.Red("DELETED"),
	'M': color.Yellow("MODIFIED"),
	'R': color.LightGreen("RENAMED"),
	'?': color.Purple("UNTRACKED"),
	'-': color.Green("CLEAN"),
}

var Location = [...]string{
	color.LightCyan("Index"),
	color.LightCyan("Working Tree"),
}

var ErrorSymbol = color.Red("X")
var SkipDir = color.Cyan("Directory to SKIP:")
