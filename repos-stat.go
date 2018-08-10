package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func help() {
	fmt.Println()
	fmt.Println("** repos-stat **")
	fmt.Println()
	fmt.Println("Usage")
	fmt.Println("    repo-stat [options] directory-to-visit [directory-to-skip]")
	fmt.Println()
	fmt.Println("Options")
	fmt.Println("    --no-clean         Skip clean repository")
	fmt.Println("    --no-broken        Skip broken repository")
	fmt.Println()
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			help()
		}
	}()

	skipClean := false
	skipBroken := false

	visitDirs := make([]string, 10)
	excludeDirs := make(map[string]bool)
	for _, arg := range os.Args[1:] {
		if len(arg) > 2 && arg[:2] == "--" {
			if arg[2:] == "no-clean" {
				skipClean = true
			} else if arg[2:] == "no-broken" {
				skipBroken = true
			} else {
				panic("Option " + arg + " unknown");	
			}
		} else if arg[0] == '-' {
			excludeDirs[arg[1:]] = true
		} else {
			visitDirs = append(visitDirs, arg)
		}
	}

	if len(excludeDirs) > 0 {
		fmt.Println(SkipDir)
		for d := range excludeDirs {
			fmt.Println("   " + d)
		}
		fmt.Println()
	}

	for _, dir := range visitDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) (dirErr error) {	
			if err != nil {
				return
			}

			if _, toSkip := excludeDirs[path]; info.IsDir() && toSkip {
				return filepath.SkipDir
			}

			if _, fileErr := os.Stat(path + "/.git/"); info.IsDir() && fileErr == nil {
				dirErr = GetStatus(path, skipClean, skipBroken)
			}
			return
		})
	}
}
