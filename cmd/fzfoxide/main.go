package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/petemango/fzfoxide/pkg/dir"
)

func usage() string {
	return "Usage: \n  --run [path]\n  --record [path]\n  --query [path]\n"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage())
		os.Exit(1)
	}

	acceptedCommands := []string{
		"--run",
		"--record",
		"--query",
	}

	command := os.Args[1]
	if !slices.Contains(acceptedCommands, command) {
		fmt.Printf("invalid command.\n%s", usage())
		os.Exit(1)
	}

	switch command {
	case "--run":
		path := os.Args[2]

		dir, err := dir.Run(path)
		if err != nil {
			fmt.Println("could not find the proper path")
			os.Exit(1)
		}
		fmt.Println(dir)
	default:
		fmt.Printf("invalid command.\n%s", usage())
		os.Exit(1)
	}
}
