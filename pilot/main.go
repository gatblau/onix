package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	var (
		cmdStr = os.Args[1]
		vars   string
	)
	for i := 2; i < len(os.Args); i++ {
		vars += os.Args[i] + " "
	}
	child := exec.Command(cmdStr, strings.Trim(vars, " "))
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	err := child.Start()

	if err != nil {
		log.Fatal(err)
	}

	child.Wait()
}
