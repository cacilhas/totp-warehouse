package main

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/cacilhas/totp-warehouse/gui"
)

func main() {
	if search(os.Args, "-nofork") {
		gui.Start()

	} else {
		// Fork
		binary, _ := exec.LookPath(os.Args[0])
		if err := syscall.Exec(binary, append(os.Args, "-nofork"), os.Environ()); err != nil {
			panic(err)
		}
		os.Exit(0)
	}
}

func search(args []string, element string) bool {
	for _, value := range args {
		if value == element {
			return true
		}
	}
	return false
}
