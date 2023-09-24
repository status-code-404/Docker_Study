package cgroup

import (
	"fmt"
	"os"
	"path"
)

const CGROUP_PATH = "/sys/fs/cgroup/memory/cgroup_test_memory_limit"
const MEMORY_LIMIT_FILE = "memory.limit_in_bytes"
const TASKS = "tasks"

func SetMemoryLimit(pid int, memoryLimit int) error {
	file, err := os.OpenFile(path.Join(CGROUP_PATH, MEMORY_LIMIT_FILE), os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(fmt.Sprintf("%dM", memoryLimit)))
	if err != nil {
		return err
	}

	file, err = os.OpenFile(path.Join(CGROUP_PATH, TASKS), os.O_WRONLY, 0777)
	_, err = file.Write([]byte(fmt.Sprintf("%d", pid)))
	if err != nil {
		return err
	}
	return nil
}
