package main

import (
	"path/filepath"
	"os"
	"os/exec"
	"fmt"
	"bytes"
	"strings"
	"github.com/logrusorgru/aurora"
)

var status = map[string]aurora.Value {
	//"M ":	"updated in index",
	//"MM":	"updated in index",
	//"MD":	"updated in index",
	//"A ":	"added to index",
	//"AM":	"added to index",
	//"AD":	"added to index",
	//"D":		"deleted from index",
	//"R ":	"renamed in index",
	//"RM":	"renamed in index",
	//"RD":	"renamed in index",
	//"C ":	"copied in index",
	//"CM":	"copied in index",
	//"CD":	"copied in index",
	" M":	aurora.Brown("MODIFIED"),
	" D":	aurora.Red("DELETED"),
	//"DR":	"renamed in work tree",
	//" DR":	"renamed in work tree",
	//"DC":	"copied in work tree",
	//" DC":	"copied in work tree",
	//"DD":		"unmerged, both deleted",
	//"AU":		"unmerged, added by us",
	//"UD":		"unmerged, deleted by them",
	//"UA":		"unmerged, added by them",
	//"DU":		"unmerged, deleted by us",
	//"AA":		"unmerged, both added",
	//"UU":		"unmerged, both modified",
	"??":		aurora.Magenta("UNTRACKED"),
	//"!!":		"ignored",
	"--":		aurora.Green("CLEAN"),
}

func main() {
	fmt.Println("###############################")
	fmt.Println("## Git repo status in Golang ##")
	fmt.Println("###############################")
	fmt.Println()

	visitDirs := make([]string, 10)
	excludeDirs := make(map[string]bool)
	for _, dir := range os.Args[1:] {
		if dir[0] == '-' {
			excludeDirs[dir[1:]] = true
		} else {
			visitDirs = append(visitDirs, dir)
		}
	}

	fmt.Println("Directory to SKIP")
	for d, _ := range excludeDirs {
		fmt.Println("   " + d)
	}
	fmt.Println()

	for _, dir := range visitDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			/* defer func(){
				fmt.Println("EXITING", path)
			}()
			fmt.Println("ENTERING", path) */

			if err != nil {
				//fmt.Println("Error", err, "at", path)
				return err
			}
			if _, toSkip := excludeDirs[path]; info.IsDir() && toSkip {
				//fmt.Println("Skipping", path)
				return filepath.SkipDir
			}
			
			if _, fileErr := os.Stat(path + "/.git/"); info.IsDir() && fileErr == nil {
				cdErr := os.Chdir(path)
				if cdErr != nil {
					fmt.Printf("Can't change dir")
				}
				cmd := exec.Command("git", "status", "-s")
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()

				errCode := ""
				if err != nil {
					errCode = "X "
				}
				
				cmd = exec.Command("git", "rev-parse", "--show-toplevel")
				var name bytes.Buffer
				cmd.Stdout = &name
				err = cmd.Run()

				repoName := ""
				if err == nil {
					subDir := strings.Split(name.String(), "/")
					repoName = subDir[len(subDir)-1]
				}

				fmt.Printf("%v%v%v - %v", aurora.Red(errCode), aurora.Blue("GIT"), aurora.Blue(":" + path), repoName)
				//fmt.Println(aurora.Blue("GIT") + aurora.Red(errCode) + aurora.Blue(":" + path))
				if err == nil {
					files := strings.Split(out.String(), "\n")
					if len(files) == 1 {
						fmt.Printf("%v\n\n", status["--"])
						return filepath.SkipDir
					}
					var messageType string
					for _, line := range files {
						if len(line) < 2 {
							continue
						}
						msgType := line[:2]
						if msgType != messageType {
							fmt.Println(status[msgType])
							messageType = msgType
						}
						fmt.Printf("    %v\n", line[2:])
					}
				}
				fmt.Println()
				return filepath.SkipDir
			}
			return nil
		})
	}
}
