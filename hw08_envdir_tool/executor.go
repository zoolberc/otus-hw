package main

import (
	"fmt"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return -1
	}
	command := exec.Command(cmd[0], cmd[1:]...)
	envs := updateEnv(env)
	command.Env = append(os.Environ(), envs...)
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	if err := command.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return command.ProcessState.ExitCode()
}
func updateEnv(e Environment) (env []string) {
	for k, v := range e {
		if _, ok := os.LookupEnv(k); ok {
			os.Unsetenv(k)
		}
		if !v.NeedRemove {
			s := fmt.Sprintf("%s=%s", k, v.Value)
			env = append(env, s)
		}
	}
	return env
}
