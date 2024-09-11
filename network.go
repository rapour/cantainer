package cantainer

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Socket struct {
	Address string
	Port    uint
}

func (s Socket) ExtendedAddress() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}

func GetAddress() (string, error) {

	cmd := exec.Command(`bash`, `-c`, `ip route | grep default | awk '{print $9}'`)

	var stdErr, stdOut bytes.Buffer

	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error executing get address command: %v", err)
	}

	if errStr := stdErr.String(); errStr != "" {
		return "", fmt.Errorf("error executing get address command: %s", errStr)
	}

	return strings.TrimSuffix(stdOut.String(), "\n"), nil
}
