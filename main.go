package main

import (
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/cacilhas/totp-warehouse/config"
	"github.com/cacilhas/totp-warehouse/gui"
)

func main() {
	if search(os.Args, "-nofork") {
		start()
	} else {
		fork()
	}
}

func start() {
	configdir := config.ConfigDir()
	lockname := path.Join(configdir, "lock")
	if lockfile, err := os.OpenFile(lockname, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600); err == nil {
		lockfile.Close()
	} else {
		return // TOTP Warehouse is running yet
	}
	defer os.Remove(lockname)
	gui.Start()
}

func fork() {
	args := append(os.Args[1:], "-nofork")
	exec.Command(os.Args[0], args...).Start()
	time.Sleep(2 * time.Second)
	os.Exit(0) // force quit
}

func search(args []string, element string) bool {
	for _, value := range args {
		if value == element {
			return true
		}
	}
	return false
}
