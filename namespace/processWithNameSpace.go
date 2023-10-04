package namespace

import (
	"io"
	"os"
	"os/exec"
	"syscall"
)

func CreateNewProcess(path string, inputCommand ...string) (*exec.Cmd, io.Writer, error) {
	if len(path) == 0 {
		path = "sh"
	}

	realCommand := make([]string, 0)

	if len(inputCommand) > 0 {
		if inputCommand[0] == "-c" {
			realCommand = inputCommand
		} else {
			realCommand = append(realCommand, "-c")
			realCommand = append(realCommand, inputCommand...)
		}
	}

	cmd := exec.Command(path, realCommand...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Run in Linux
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			syscall.SysProcIDMap{
				//Container Id 本Namespace中采用的user id 使用Id命令能访问到, 0 是root
				ContainerID: 0,
				// Host Id  (调用跟这个程序的）当前用户id
				HostID: 0,
				Size:   1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			syscall.SysProcIDMap{
				ContainerID: 0,
				HostID:      0,
				Size:        1,
			},
		},
	}

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, stdinPipe, nil
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
