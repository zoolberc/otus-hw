package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const testDir = "./testdata/env"

func TestReadDir(t *testing.T) {
	expected := Environment{
		"BAR":   EnvValue{"bar", false},
		"EMPTY": EnvValue{"", false},
		"FOO":   EnvValue{"   foo\nwith new line", false},
		"HELLO": EnvValue{"\"hello\"", false},
		"UNSET": EnvValue{"", true},
	}

	t.Run("Get envs", func(t *testing.T) {
		envs, err := ReadDir(testDir)
		require.NoError(t, err)
		require.Equal(t, expected, envs)
	})

	t.Run("dir not exists", func(t *testing.T) {
		_, err := ReadDir("test_envs")
		require.NotEqual(t, nil, err)
	})
}
