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

//type Credential struct {
//// Uid represents the user ID of the process.
//Uid uint32
//// Gid represents the group ID of the process.
//Gid uint32
//// Groups represents the supplementary group IDs of the process.
//Groups []uint32
//// Username represents the username of the process owner.
//Username string
//}
// use pstree -pl see the pid , then use echo $$ in your go environment see the pid in namespace(1) Its ineresting
