package cantainer

import (
	"bytes"
	"errors"
	"os/exec"
)

func GetAddress() (string, error) {

	cmd := exec.Command("bash", "-c", "sudo ip route | grep default | awk '{print $9s}'")

	var stdErr, stdOut bytes.Buffer

	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	if errStr := stdErr.String(); errStr != "" {
		return "", errors.New(errStr)
	}

	return stdOut.String(), nil
}
