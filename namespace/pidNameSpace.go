package namespace

import (
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func CreateNewProcess(path string, inputCommand ...string) (io.Writer, *exec.Cmd, error) {
	if len(path) == 0 {
		path = "sh"
	}
	cmd := exec.Command(path, inputCommand...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Run in Linux
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
	}

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	println(cmd.Process.Pid)
	return stdinPipe, cmd, nil
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
