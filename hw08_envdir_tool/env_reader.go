package main

import (
	"bytes"
	"os"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	envMap := make(Environment)
	for i := 0; i < len(files); i++ {
		fileInfo, err := files[i].Info()
		if err != nil {
			return nil, err
		}
		if fileInfo.IsDir() {
			continue
		}
		if fileInfo.Size() == 0 {
			envMap[fileInfo.Name()] = EnvValue{Value: "", NeedRemove: true}
		}
		envFile, err := os.ReadFile(dir + "/" + fileInfo.Name())
		if err != nil {
			return nil, err
		}
		envMap[strings.TrimRight(fileInfo.Name(), "=")] = EnvValue{
			validateEnvValue(string(envFile)),
			fileInfo.Size() == 0,
		}
	}
	return envMap, nil
}

func validateEnvValue(value string) string {
	fistLine := strings.Split(value, "\n")[0]
	fistLine = string(bytes.ReplaceAll([]byte(fistLine), []byte{0x00}, []byte("\n")))
	return strings.TrimRight(fistLine, " \t\r")
}
