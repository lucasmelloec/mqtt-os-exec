package osExec

import (
	"log"
	"os/exec"
	"strings"
)

func Execute(command []string) {
	for _, value := range command {
		cmdSlice := strings.Fields(value)
		cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)

		if err := cmd.Run(); err != nil {
			log.Println("Error executing command", value)
			return
		}
	}
}
