package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fedemengo/go-data-structures/heap"
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

func printRepo(path string, files []string, broken, clean bool) {
	if broken {
		fmt.Printf("%v %v\n\n", PathColored(path), ErrorSymbol)
		return
	}

	fmt.Printf("%v\n", PathColored(path))
	if clean {
		fmt.Printf("%v\n\n", Message['-'])
	} else {
		container := make([]*heap.Heap, 2)
		for i := range container {
			container[i] = heap.NewHeap(func(e1, e2 heap.Elem) bool {
				return e1.Key.(byte) < e2.Key.(byte)
			})
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
					container[idx].Push(heap.Elem{Key: c, Val: file[2:]})
				}
			}
		}

		for idx := range container {
			messageType := ""
			fmt.Printf("%v\n", Location[idx])
			if container[idx].Size() == 0 {
				fmt.Printf("  %v\n", Message['-'])
			} else {
				for container[idx].Size() > 0 {
					file := container[idx].Pop()
					msgType := Message[file.Key.(byte)]
					if messageType != msgType {
						fmt.Printf("  %v\n", msgType)
						messageType = msgType
					}
					fmt.Printf("        %v\n", file.Val.(string))
				}
			}
			fmt.Printf("\n")
		}
	}
}
