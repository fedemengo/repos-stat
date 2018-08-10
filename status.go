package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetStatus display the status of each repo
func GetStatus(repoPath string, skipClean, skipBroken bool) error {

	if cdErr := os.Chdir(repoPath); cdErr != nil {
		return cdErr
	}

	var out bytes.Buffer
	gitStatus := exec.Command("git", "status", "-s")
	gitStatus.Stdout = &out

	broken := false
	if err := gitStatus.Run(); err != nil {
		broken = true
	}

	files := strings.Split(out.String(), "\n")
	clean := false
	if !broken && len(files) == 1 {
		clean = true
	}

	if (!broken && !clean) || (broken && !skipBroken) || (clean && !skipClean) {
		printRepo(repoPath, files, broken, clean)
	}

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

func printRepo(path string, files []string, broken, clean bool) {
	if broken {
		fmt.Printf("%v %v\n\n", PathColored(path), ErrorSymbol)
		return
	}

	fmt.Printf("%v\n", PathColored(path))
	if clean {
		fmt.Printf("%v\n\n", Message['-'])
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
			fmt.Printf("%v\n", Location[idx])
			
			if container[idx].Len() == 0 {
				fmt.Printf("  %v\n", Message['-'])
			} else {
				for container[idx].Len() > 0 {
					file := heap.Pop(&container[idx])
					msgType := Message[(file.(data)).code]
					if messageType != msgType {
						fmt.Printf("  %v\n", msgType)
						messageType = msgType
					}
					fmt.Printf("        %v\n", (file.(data)).fileName)
				}
			}
			fmt.Printf("\n")
		}
	}
}
