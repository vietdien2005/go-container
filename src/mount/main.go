package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)

func main() {
	// wget https://raw.githubusercontent.com/vietdien2005/go-container/master/layer/alpine_linux.tar
	// mkdir -p /tmp/go_container/rootfs
	// tar -xvf alpine_linux.tar -C /tmp/go_container/rootfs

	cmd := reexec.Command("bootstrapProcess", "/tmp/go_container/rootfs")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	cmd.Run()
}

func init() {
	reexec.Register("bootstrapProcess", bootstrapProcess)

	if reexec.Init() {
		os.Exit(0)
	}
}

func bootstrapProcess() {
	fmt.Printf("\n>>[Bootstrap Process]\n\n<<")

	mountPath := os.Args[1]

	mountProcess(mountPath)
	pivotRoot(mountPath)

	syscall.Sethostname([]byte("go-container"))

	runProcess()
}

func runProcess() {
	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=[go-container]-$ "}

	cmd.Run()
}

func pivotRoot(newRootFS string) {
	oldRootFS := filepath.Join(newRootFS, "/.pivot_root")

	syscall.Mount(newRootFS, newRootFS, "", syscall.MS_BIND|syscall.MS_REC, "")
	os.MkdirAll(oldRootFS, 0700)

	syscall.PivotRoot(newRootFS, oldRootFS)

	os.Chdir("/")

	oldRootFS = "/.pivot_root"

	syscall.Unmount(oldRootFS, syscall.MNT_DETACH)

	os.RemoveAll(oldRootFS)
}

func mountProcess(newRootFS string) {
	source := "proc"
	target := filepath.Join(newRootFS, "/proc")
	fstype := "proc"
	flags := 0
	data := ""

	os.MkdirAll(target, 0755)

	syscall.Mount(source, target, fstype, uintptr(flags), data)
}
