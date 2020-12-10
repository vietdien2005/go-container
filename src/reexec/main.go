package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)

func main() {
	cmd := reexec.Command("bootstrapProcess")

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
