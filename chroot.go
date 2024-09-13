package cantainer

import (
	"os"
	"os/exec"
	"syscall"
)

func Chroot(root string, call string, optinos ...string) {

	oldrootHandle, err := os.Open("/")
	if err != nil {
		panic(err)
	}
	defer oldrootHandle.Close()

	cmd := exec.Command(call, optinos...)

	err = syscall.Chdir(root)
	if err != nil {
		panic(err)
	}

	err = syscall.Chroot(root)
	if err != nil {
		panic(err)
	}

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	err = syscall.Fchdir(int(oldrootHandle.Fd()))
	if err != nil {
		panic(err)
	}

	err = syscall.Chroot(".")
	if err != nil {
		panic(err)
	}
}
