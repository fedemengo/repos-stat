package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bclicn/color"
)

func help() {
	fmt.Println("Golang repository status")
	fmt.Println()
	fmt.Println("Usage")
	fmt.Println("    repo-stat [options] directory-to-visit [directory-to-skip]")
	fmt.Println()
	fmt.Println("Options")
	fmt.Println("    --no-clean\tSkip clean repository")
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			help()
		}
	}()

	skipClean := false
	skipBroken := false

	visitDirs := make([]string, 10)
	excludeDirs := make(map[string]bool)
	for _, dir := range os.Args[1:] {
		if len(dir) > 2 && dir[:2] == "--" {
			if dir[2:] == "no-clean" {
				skipClean = true
			} else if dir[2:] == "no-broken" {
				skipBroken = true
			} else {
				panic("Option unknown")
			}
		} else if dir[0] == '-' {
			excludeDirs[dir[1:]] = true
		} else {
			visitDirs = append(visitDirs, dir)
		}
	}

	if len(excludeDirs) > 0 {
		fmt.Println(color.Cyan("Directory to SKIP:"))
		for d := range excludeDirs {
			fmt.Println("   " + d)
		}
		fmt.Println()
	}

	for _, dir := range visitDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) (dirErr error) {
			/* defer func(){
				fmt.Println("EXITING", path)
			}()
			fmt.Println("ENTERING", path) */

			if err != nil {
				//fmt.Println("Error", err, "at", path)
				return
			}
			if _, toSkip := excludeDirs[path]; info.IsDir() && toSkip {
				//fmt.Println("Skipping", path)
				return filepath.SkipDir
			}

			if _, fileErr := os.Stat(path + "/.git/"); info.IsDir() && fileErr == nil {
				dirErr = GetStatus(path, skipClean, skipBroken)
			}
			return
		})
	}
}
