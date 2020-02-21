package main

import (
	"os"
	"os/exec"
	"time"

	"github.com/cacilhas/totp-warehouse/gui"
)

func main() {
	if search(os.Args, "-nofork") {
		gui.Start()

	} else {
		// Fork
		args := append(os.Args[1:], "-nofork")
		exec.Command(os.Args[0], args...).Start()
		time.Sleep(2 * time.Second)
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
