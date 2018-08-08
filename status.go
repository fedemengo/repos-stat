package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bclicn/color"
)

var loc = [...]string{"Index", "Working Tree"}

// GetStatus display the status of each repo
func GetStatus(repoPath string, skipClean, skipBroken bool) error {

	if cdErr := os.Chdir(repoPath); cdErr != nil {
		fmt.Printf("Can't change dir")
		return cdErr
	}

	var out bytes.Buffer
	gitStatus := exec.Command("git", "status", "-s")
	gitStatus.Stdout = &out

	errCode := ""
	if err := gitStatus.Run(); err != nil {
		errCode = "X"
	}

	files := strings.Split(out.String(), "\n")

	if errCode != "" && skipBroken {
		return filepath.SkipDir
	}

	if errCode != "X" && len(files) == 1 && skipClean {
		return filepath.SkipDir
	}

	printRepo(repoPath, errCode, files)
	return filepath.SkipDir
}

type data struct {
	code     byte
	fileName string
}

type heapData []data

func (h heapData) Len() int           { return len(h) }
func (h heapData) Less(i, j int) bool { return h[i].code < h[j].code }
func (h heapData) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push insert an element into the heap
func (h *heapData) Push(x interface{}) {
	*h = append(*h, x.(data))
}

// Pop remove the element on top of the heap according to the heap priority
func (h *heapData) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func getMessage(symbol byte) string {
	switch symbol {
	case 'A':
		return color.Green("ADDED")
	case 'D':
		return color.Red("DELETED")
	case 'M':
		return color.Yellow("MODIFIED")
	case '?':
		return color.Purple("UNTRACKED")
	case '-':
		return color.Green("CLEAN")
	}
	return ""
}

func printRepo(path, repoError string, files []string) {
	if repoError != "" {
		fmt.Printf("%v %v\n", color.Red(repoError), color.Blue(path))
		fmt.Println()
		return
	}

	fmt.Printf("%v%v\n", color.Red(repoError), color.Blue(path))
	if len(files) == 1 {
		fmt.Println(getMessage('-'))
		fmt.Println()
	} else {
		var container [2]heapData
		for _, h := range container {
			h = make(heapData, 10)
			heap.Init(&h)
		}

		for _, file := range files {
			if len(file) < 2 {
				continue
			}

			if file[:2] == "??" {
				file = "? " + file[2:]
			}

			for idx := range container {
				if c := file[idx]; c != ' ' {
					heap.Push(&container[idx], data{code: c, fileName: file[2:]})
				}
			}
		}

		var messageType string
		for idx := range container {
			fmt.Println(color.LightCyan(loc[idx]))
			
			if container[idx].Len() == 0 {
				fmt.Println("  ", getMessage('-'))
			} else {
				for container[idx].Len() > 0 {
					file := heap.Pop(&container[idx])
					msgType := getMessage((file.(data)).code)
					if messageType != msgType {
						fmt.Println("  ", msgType)
						messageType = msgType
					}
					fmt.Println("\t", (file.(data)).fileName)
				}
			}
			fmt.Println()
		}
	}
}
