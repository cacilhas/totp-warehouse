package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Open(filename string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", filename)
	case "windows":
		cmd = exec.Command("rundll32", "filename.dll,FileProtocolHandler", filename)
	case "darwin":
		cmd = exec.Command("open", filename)
	default:
		return nil, fmt.Errorf("unsupported platform")
	}
	return cmd, cmd.Start()
}

func Kill(cmd *exec.Cmd) {
	signal := os.Interrupt
	if runtime.GOOS == "windows" {
		signal = os.Kill
	}
	proc := cmd.Process
	proc.Signal(signal)
	proc.Kill()
	proc.Wait()
}
