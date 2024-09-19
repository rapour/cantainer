package cantainer

import (
	"os"
	"os/exec"
	"syscall"
)

func NewContainer(name string, root string, call string, options ...string) {

	cmd := exec.Command("/proc/self/exe", append([]string{"child", root, call}, options...)...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	handle, err := GetNetNamespaceHandleFromName(name)
	if err != nil {
		panic(err)
	}

	err = SetNamespace(handle)
	if err != nil {
		panic(err)
	}

	err = cmd.Run()
	if err != nil {
		panic(err)
	}

}

func Child(root string, call string, optinos ...string) {

	oldrootHandle, err := os.Open("/")
	if err != nil {
		panic(err)
	}
	defer oldrootHandle.Close()

	cmd := exec.Command(call, optinos...)

	err = syscall.Chroot(root)
	if err != nil {
		panic(err)
	}

	err = syscall.Chdir("/")
	if err != nil {
		panic(err)
	}

	err = syscall.Mount("proc", "proc", "proc", 0, "")
	if err != nil {
		panic(err)
	}

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	syscall.Sethostname([]byte("cantainer"))

	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	err = syscall.Unmount("proc", 0)
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
