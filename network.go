package cantainer

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/netip"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

func VethNameFronNamespace(namespace string) string {
	return namespace
}

type NsHandle int

const bindMountPath = "/run/netns" /* Bind mount path for named netns */

type Socket struct {
	Address string
	Port    uint
}

func (s Socket) ExtendedAddress() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}

func CreateNetworkNamespace(name string) error {
	checkCmd := exec.Command("bash", "-c", fmt.Sprintf("ip netns add %s", name))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {
		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error executing add netns command: %v [%s]", err, errStr)
	}
	return nil
}

func DeleteNetworkNamespace(name string) error {
	checkCmd := exec.Command("bash", "-c", fmt.Sprintf("ip netns delete %s", name))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {
		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error executing add netns command: %v [%s]", err, errStr)
	}
	return nil
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

func CreateBridge(name string) error {

	checkCmd := exec.Command(`bash`, `-c`, fmt.Sprintf(`scripts/bridge.sh %s`, name))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {
		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error executing bridge command: %v [%s]", err, errStr)
	}

	slog.Info(strings.TrimSuffix(stdOut.String(), "\n"))

	return nil
}

func CreateVXLan(name string, id int, bridgeName string) error {

	checkCmd := exec.Command(`bash`, `-c`, fmt.Sprintf(`scripts/vxlan.sh %s %s %d`, name, bridgeName, id))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {
		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error executing vxlan command: %v [%s]", err, errStr)
	}

	slog.Info(strings.TrimSuffix(stdOut.String(), "\n"))

	return nil
}

func AddRemoteToVXLan(name string, address string) error {

	checkCmd := exec.Command(`bash`, `-c`, fmt.Sprintf(`scripts/vxlan-add-remote.sh %s %s`, name, address))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {
		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error executing add to vxlan command: %v [%s]", err, errStr)
	}

	slog.Info("added node to vxlan", slog.String("address", address))
	return nil
}

func RemoveFromVXLan(name string, address string) error {

	checkCmd := exec.Command(`bash`, `-c`, fmt.Sprintf(`scripts/vxlan-remove-remote.sh %s %s`, name, address))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {

		// TODO: fix the error (in some cases MAC doesn't match)
		if exiterr, ok := err.(*exec.ExitError); ok {
			if exiterr.ExitCode() == 255 {
				slog.Warn("the entry didn't found", slog.String("error", exiterr.Error()))
				return nil
			}
		}

		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error executing remove from vxlan command: %v [%s]", err, errStr)
	}

	slog.Info("removed node from vxlan", slog.String("address", address))
	return nil
}

func RemoveVethPair(name string) error {

	checkCmd := exec.Command(`bash`, `-c`, fmt.Sprintf(`scripts/veth-remove.sh %s`, name))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {
		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error executing remove veth command: %v [%s]", err, errStr)
	}

	slog.Info(strings.TrimSuffix(stdOut.String(), "\n"))

	return nil

}

func CreateVethPair(name string) error {

	checkCmd := exec.Command(`bash`, `-c`, fmt.Sprintf(`scripts/veth-add.sh %s`, name))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {
		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error executing create veth command: %v [%s]", err, errStr)
	}

	slog.Info(strings.TrimSuffix(stdOut.String(), "\n"))

	return nil

}

func AddVethToBridge(vethName string, bridgeName string) error {

	checkCmd := exec.Command(`bash`, `-c`, fmt.Sprintf(`scripts/veth-add-to-bridge.sh %s %s`, vethName, bridgeName))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {
		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error adding veth to bridge: %v [%s]", err, errStr)
	}

	slog.Info(strings.TrimSuffix(stdOut.String(), "\n"))

	return nil
}

func AddVethToNamespace(namespace string, vethName string) error {

	checkCmd := exec.Command(`bash`, `-c`, fmt.Sprintf(`scripts/veth-add-to-namespace.sh %s %s`, vethName, vethName))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {
		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error adding veth to namespace: %v [%s]", err, errStr)
	}

	slog.Info(strings.TrimSuffix(stdOut.String(), "\n"))

	return nil

}

func AssignNetworkToNamespace(namespace string, prefix netip.Prefix) error {
	checkCmd := exec.Command(`bash`, `-c`, fmt.Sprintf(`scripts/assign-net-to-namespace.sh %s %s %s`, namespace, VethNameFronNamespace(namespace), prefix.String()))

	var stdErr, stdOut bytes.Buffer

	checkCmd.Stderr = &stdErr
	checkCmd.Stdout = &stdOut

	if err := checkCmd.Run(); err != nil {
		errStr := strings.TrimSuffix(stdErr.String(), "\n")
		return fmt.Errorf("error assigning ip to namespace: %v [%s]", err, errStr)
	}

	slog.Info(strings.TrimSuffix(stdOut.String(), "\n"))

	return nil
}

func ConnectNetworkNamespaceToBridge(namespace string, bridgeName string) error {

	vethName := VethNameFronNamespace(namespace)

	if err := CreateVethPair(vethName); err != nil {
		return err
	}

	if err := AddVethToNamespace(namespace, vethName); err != nil {
		return err
	}

	if err := AddVethToBridge(vethName, bridgeName); err != nil {
		return err
	}

	return nil
}

// GetNetNamespaceHandleFromName gets a handle to a named network namespace such as one
// created by `ip netns add`.
func GetNetNamespaceHandleFromName(name string) (NsHandle, error) {
	return GetNetNamespaceHanddleFromPath(filepath.Join(bindMountPath, name))
}

// GetNetNamespaceHanddleFromPath gets a handle to a network namespace
// identified by the path
func GetNetNamespaceHanddleFromPath(path string) (NsHandle, error) {
	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, err
	}
	return NsHandle(fd), nil
}

// SetNamespace sets the current network namespace to the namespace represented
// by NsHandle.
func SetNamespace(ns NsHandle) error {
	return unix.Setns(int(ns), unix.CLONE_NEWNET)
}
