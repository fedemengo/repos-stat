package main

import (
	"path/filepath"
	"os"
	. "path"
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

	visitDirs := make([]string, 10)
	excludeDirs := make(map[string]bool)
	for _, dir := range os.Args[1:] {
		if dir[0] == '-' {
			excludeDirs[dir[1:]] = true
		} else {
			visitDirs = append(visitDirs, dir)
		}
	}

	fmt.Println("Directory to SKIP:")
	for d, _ := range excludeDirs {
		fmt.Println("   " + d)
	}
	fmt.Println()

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
				dirErr = statusRepo(path)
			}
			return
		})
	}
}

func statusRepo(path string) error {

	if cdErr := os.Chdir(path); cdErr != nil {
		fmt.Printf("Can't change dir")
		return cdErr
	}

	var out bytes.Buffer
	gitStatus := exec.Command("git", "status", "-s")
	gitStatus.Stdout = &out

	errCode := ""
	if err := gitStatus.Run(); err != nil {
		errCode = "X "
	}
	
	var name bytes.Buffer
	gitName := exec.Command("git", "rev-parse", "--show-toplevel")
	gitName.Stdout = &name
	gitName.Run()

	repoName := Base(name.String())
	if index := strings.Index(repoName, "\n"); index != -1 {
		repoName = repoName[:index]
	}
	fmt.Printf("%v%v%v - %v\n", aurora.Red(errCode), aurora.Blue("GIT"), aurora.Blue(":" + path), repoName)

	files := strings.Split(out.String(), "\n")
	if len(files) == 1 {
		fmt.Println(status["--"])
	} else {
		var messageType string
		for _, file := range files {
			if len(file) < 2 {
				continue
			}
			msgType := file[:2]
			if msgType != messageType {
				fmt.Println(status[msgType])
				messageType = msgType
			}
			fmt.Println("    ", file[2:])
		}
	}
	fmt.Println()
	return filepath.SkipDir
}
