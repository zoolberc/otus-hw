package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("change envs", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("FOO", "bar")            // value should change
		os.Setenv("UNSET", "")             // should be removed
		os.Setenv("BAR", "foo")            // value should change to 'bar'
		os.Setenv("EMPTY", "IS_NOT_EMPTY") // should be emptied
		os.Setenv("ZIKA", "zika")          // should be preserved

		s := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"FOO": EnvValue{
				Value:      "\x20\x20\x20\x66\x6f\x6f\n\x77\x69\x74\x68\x20\x6e\x65\x77\x20\x6c\x69\x6e\x65",
				NeedRemove: false,
			},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		reader, writer, _ := os.Pipe()
		defer reader.Close()

		oldStdout := os.Stdout
		os.Stdout = writer

		defer func() {
			os.Stdout = oldStdout
		}()

		rc := RunCmd([]string{"/bin/bash", "-c", "env"}, s)

		output := make(chan string)
		go func() {
			var buf [1024]byte
			n, _ := reader.Read(buf[:])
			output <- string(buf[:n])
		}()

		r := <-output

		for k, v := range s {
			if v.NeedRemove {
				require.NotContains(t, r, k+"=")
			} else {
				require.Contains(t, r, k+"="+v.Value)
			}
		}

		require.Equal(t, 0, rc)
	})
}
