package main

import "os"

func main() {
	dir := os.Args[1]
	command := os.Args[2:]
	v, e := ReadDir(dir)
	if e != nil {
		os.Exit(1)
	}
	os.Exit(RunCmd(command, v))
}
