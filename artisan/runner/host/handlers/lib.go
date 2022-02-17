package handlers

import (
	"bufio"
	"fmt"
	"net/http"
	"os/exec"
	"syscall"
)

func checkErr(w http.ResponseWriter, msg string, err error) bool {
	if err != nil {
		msg := fmt.Sprintf("%s: %s\n", msg, err)
		fmt.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
	}
	return err != nil
}

func execute(name string, w http.ResponseWriter, args []string) error {
	command := exec.Command(name, args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		fmt.Printf("failed creating command stdoutpipe: %s", err)
		return err
	}
	defer func() {
		_ = stdout.Close()
	}()
	stdoutReader := bufio.NewReader(stdout)

	stderr, err := command.StderrPipe()
	if err != nil {
		fmt.Printf("failed creating command stderrpipe: %s", err)
		return err
	}
	defer func() {
		_ = stderr.Close()
	}()
	stderrReader := bufio.NewReader(stderr)

	if err = command.Start(); err != nil {
		return err
	}

	go handleReader(stdoutReader, w)
	go handleReader(stderrReader, w)

	if err = command.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if _, ok = exitErr.Sys().(syscall.WaitStatus); ok {
				fmt.Printf("run command failed: '%s' - '%s'", name, err)
			}
		}
		return err
	}
	return nil
}

func handleReader(reader *bufio.Reader, w http.ResponseWriter) {
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		w.Write([]byte(str))
		fmt.Printf("! %s\n", []byte(str))
	}
}
