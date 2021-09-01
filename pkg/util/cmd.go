package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func ExecCmd(cmdline string) {
	parts := strings.Split(cmdline, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatalln("failed to exec", parts[0], err)
	}
}

func FindCommand(cmd string) (string, bool) {
	paths := []string{
		fmt.Sprintf("/bin/%s", cmd),
		fmt.Sprintf("/sbin/%s", cmd),
		fmt.Sprintf("/usr/bin/%s", cmd),
		fmt.Sprintf("/usr/sbin/%s", cmd),
		fmt.Sprintf("/usr/local/bin/%s", cmd),
		fmt.Sprintf("/usr/local/sbin/%s", cmd),
	}
	for _, path := range paths {
		_, err := exec.LookPath(path)
		if err == nil {
			return path, true
		}
	}
	return "", false
}
