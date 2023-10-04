package commandHandler

import (
	"Docker_Study/namespace"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"reflect"
)

type CommandError struct {
	errorMessage string
}

func (c *CommandError) Error() string {
	return c.errorMessage
}

func GetAllCommands() []*cli.Command {
	r := make([]*cli.Command, 0)
	r = append(r, &cli.Command{
		Name:   "run",
		Usage:  "test docker run",
		Action: handleRun,
	}, &cli.Command{
		Name:   "ps",
		Usage:  "test docker ps",
		Action: handlePS,
	})
	return r

}

// Todo: These flags need to be put into Commands, First we put them here, but after will finish that
func GetAllFlags() []cli.Flag {
	f := make([]cli.Flag, 0)
	f = append(f, &cli.BoolFlag{
		Name:    "echo_process",
		Aliases: []string{"e"},
		Usage:   "appearance all process, Just use -e on it",
	}, &cli.BoolFlag{
		//后台运行功能，后面再写
		Name:    "echo_direct_attribute",
		Aliases: []string{"d"},
		Usage:   "show it directly or daemon",
	}, &cli.StringFlag{
		Name:    "f_test",
		Aliases: []string{"f"},
		Usage:   "appreance all property of process, Just use -f on it",
	}, &cli.StringFlag{
		Name:    "ti",
		Aliases: []string{"t"},
		Usage:   "use tty to run command",
		Value:   "default",
	},
	)
	return f
}

func handlePS(c *cli.Context) error {
	println("handle ps")
	//这里写死了一个路径，使用时需注意
	file, err := os.Open("/home/ubuntu/docker_config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 0)
	if _, err = file.Read(buf); err != nil {
		return err
	}
	var dockerProcess []DockerProcess
	if err = json.Unmarshal(buf, &dockerProcess); err != nil {
		return nil
	}

	for _, value := range dockerProcess {
		if value.Running == true || c.Bool("e") {
			if err = showParams(c.StringSlice("f"), reflect.ValueOf(value)); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleRun(c *cli.Context) error {
	println("handle run")
	// 创造cmd运行docker
	// 前两个是什么不用我多说
	cmd, stdin, err := namespace.CreateNewProcess("sh")
	if err != nil {
		println(err.Error())
		return err
	}
	if err = cmd.Start(); err != nil {
		println(err.Error())
		return err
	}

	// 先初始化容器，目前只是简单的mount proc,然后运行我们输入的命令
	if err = InitDocker(stdin); err != nil {
		println(err.Error())
		return err
	}
	commandLine := ""
	for i := 0; i < c.Args().Len(); i++ {
		commandLine += c.Args().Get(i)
		commandLine += " "
	}
	if len(commandLine) > 0 {
		commandLine += "\n"
		stdin.Write([]byte(commandLine))
	}

	// 之后转到io.stdin上监听我们的输入
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				break
			}
			stdin.Write(buf[:n])
		}
	}()

	defer cmd.Wait()
	return nil
}

// 最开始还是直接挂载proc就算初始化成功了
func InitDocker(stdin io.Writer) error {
	_, err := stdin.Write([]byte("mount -t proc proc /proc\n"))
	if err != nil {
		return err
	}
	return nil
}

func showParams(params []string, container reflect.Value) error {
	if len(params) == 0 {
		params = append(params, "Name", "ProcessId")
	}
	for _, p := range params {
		field := container.FieldByName(p)
		if field.IsValid() {
			print(field.String())
			print(" ")
		} else {
			return &CommandError{fmt.Sprintf("process has no attribute {%s}", p)}
		}
		print("\n")
	}
	return nil
}
