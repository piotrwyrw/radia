package main

import (
	"cli/internal/commands"
	"context"
	"fmt"
	"os"
)

func main() {
	err := commands.InitializeCommands().Run(context.Background(), os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
