package namespace

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func CreateNewProcess() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Run in Linux
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// use pstree -pl see the pid , then use echo $$ in your go environment see the pid in namespace(1) Its ineresting
