package commandHandler

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
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

func GetAllFlags() []cli.Flag {
	f := make([]cli.Flag, 0)
	f = append(f, &cli.BoolFlag{
		Name:    "e",
		Aliases: []string{"e"},
		Usage:   "appreance all process, Just use -e on it",
	}, &cli.BoolFlag{
		//后台运行功能，后面再写
		Name:    "d",
		Aliases: []string{"d"},
		Usage:   "show it directly or daemon",
	}, &cli.StringFlag{
		Name:    "f",
		Aliases: []string{"f"},
		Usage:   "appreance all property of process, Just use -f on it",
	},
	)
	return nil
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
