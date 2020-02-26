package helpers

import (
	"os"
	"os/exec"
	"time"
)

func Fork(args []string) {
	if args == nil {
		args = os.Args
	}
	exec.Command(args[0], args[1:]...).Start()
	time.Sleep(200 * time.Millisecond)
	os.Exit(0) // force quit
}
