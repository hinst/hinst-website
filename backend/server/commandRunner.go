package server

import (
	"log"
	"os/exec"
)

type commandRunner struct {
	Dir string
}

func (me *commandRunner) command(name string, assertive bool, main string, args ...string) bool {
	var command = exec.Command(main, args...)
	command.Dir = me.Dir
	var output, err = command.CombinedOutput()
	if err != nil {
		if assertive {
			log.Fatalln("Command error:", name, '\n', string(output), '\n', err)
		} else {
			log.Println("Command warning:", name, '\n', string(output), '\n', err)
		}
	}
	return err == nil
}
