package main

import (
	"Docker_Study/cgroup"
	"Docker_Study/namespace"
	"time"
)

func main() {
	//stress命令很特殊， 有2个pid,我们拿到的内阁Pid反而不是需要限制资源的Pid
	stdinPipe, cmd, err := namespace.CreateNewProcess("sh", make([]string, 0)...)

	defer func() {
		err = cmd.Wait()
		if err != nil {
			println(err.Error())
		}
		println("Really Done")
	}()

	if err != nil {
		println("error occurs")
		println(err.Error())
	}
	println("ready to run limit setting, now_pid is ", cmd.Process.Pid)

	//设定200m内存限制
	err = cgroup.SetMemoryLimit(cmd.Process.Pid, 200)
	if err != nil {
		println("error occurs")
		println(err.Error())
	}

	//pid在新的命名空间应该为1
	println("Please check if pid == 1")
	_, err = stdinPipe.Write([]byte("echo $$\n"))
	if err != nil {
		println(err.Error())
	}
	time.Sleep(5 * time.Second)
	println("")

	//正式运行，先运行400， 应该会获得报错 got signal 9 获取资源失败
	println("Except this failed")
	_, err = stdinPipe.Write([]byte("stress --vm-bytes 400m --vm-keep -m 1\n"))
	if err != nil {
		println(err.Error())
	}

	// 然后我们跑100M的压测应该能正常运行
	time.Sleep(5 * time.Second)
	println("")
	println("Except this normal, use top to check it")
	_, err = stdinPipe.Write([]byte("stress --vm-bytes 100m --vm-keep -m 1\n"))
	if err != nil {
		println(err.Error())
	}

	_, err = stdinPipe.Write([]byte("echo $$\n"))
	if err != nil {
		println(err.Error())
	}

}
