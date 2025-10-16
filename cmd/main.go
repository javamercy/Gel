package main

import (
	"Gel/application/services"
	"Gel/persistence/repositories"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", os.Args[0])
		os.Exit(1)
	}
	command := os.Args[1]
	repository := repositories.NewFilesystemRepository()

	switch command {
	case "init":
		{
			workingDirectory, err := os.Getwd()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
			initService := services.NewInitService(repository)
			message, err := initService.Init(workingDirectory)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", message)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "%s\n", message)
		}
	}
}
