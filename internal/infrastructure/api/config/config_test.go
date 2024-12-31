package config

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetrieveConfigurationDefault(t *testing.T) {
	// when
	applicationConfig := RetrieveConfiguration()

	// then
	require.Equal(t, RetrieveConfiguration(), applicationConfig)
}

func TestRetrieveConfigurationErrorInvalidVariable(t *testing.T) {
	// given
	t.Setenv("SERVER_PORT", "invalid")

	// when
	if os.Getenv("IS_CRASHING") == "1" {
		RetrieveConfiguration()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRetrieveConfigurationErrorInvalidVariable")
	cmd.Env = append(os.Environ(), "IS_CRASHING=1")
	err := cmd.Run()
	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestGetServerAddresses(t *testing.T) {
	// given
	applicationConfig := RetrieveConfiguration()

	// when
	require.Equal(t, "localhost:8081", applicationConfig.GetDebugServerAddress())
	require.Equal(t, "localhost:8080", applicationConfig.GetServerAddress())
}
