package otsukai

import (
	require "github.com/alecthomas/assert/v2"
	"github.com/alecthomas/repr"
	"os"
	"testing"
)

func c(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(content)
}

func TestExample_DockerCompose(t *testing.T) {
	content := c(t, "./examples/docker-compose/otsukai.rb")
	ret, err := Parser.ParseString("", content)
	require.NoError(t, err)

	require.Equal(t, 4, len(ret.Statements))
	repr.Println(ret)
}
