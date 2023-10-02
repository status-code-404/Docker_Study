package main

import (
	"Docker_Study/commandHandler"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:     "docker_test",
		Usage:    "This is a tool self-develop for understand how to develop docker",
		Flags:    commandHandler.GetAllFlags(),
		Commands: commandHandler.GetAllCommands(),
	}

	if err := app.Run(os.Args); err != nil {
		println(err.Error())
	}
}
