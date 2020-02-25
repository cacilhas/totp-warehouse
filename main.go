package main

import (
	"os"
	"path"

	"github.com/cacilhas/totp-warehouse/config"
	"github.com/cacilhas/totp-warehouse/gui"
	"github.com/cacilhas/totp-warehouse/helpers"
)

func main() {
	if search(os.Args, "-nofork") {
		start()
	} else {
		helpers.Fork(append(os.Args, "-nofork"))
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

func search(args []string, element string) bool {
	for _, value := range args {
		if value == element {
			return true
		}
	}
	return false
}
