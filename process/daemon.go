package process

import (
	"os"
	"os/exec"
)

func Daemonize(args ...string) {
	var arg []string
	if len(args) > 1 {
		arg = args[1:]
	}
	cmd := exec.Command(args[0], arg...)
	cmd.Env = os.Environ()
	cmd.Start()
}
