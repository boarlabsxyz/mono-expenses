package main

import (
	"fmt"
	"os"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		runMainFlow()

		return
	}

	command := os.Args[1]

	switch command {
	case "help", "--help", "-h":
		showHelp()
	case "version", "--version", "-v":
		showVersion()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Use 'mono-track help' to see available commands.")
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("mono-track - Quick and dirty solution for tracking and analysing expenses using AI")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("    mono-track [COMMAND]")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("    help      Show this help message")
	fmt.Println("    version   Show version information")
	fmt.Println()
	fmt.Println("When no command is provided, the main expense tracking flow is executed.")
}

func showVersion() {
	fmt.Printf("mono-track %s\n", version)
}

func runMainFlow() {
	fmt.Println("mono-track - Expense Tracking")
	fmt.Println("=============================")
	fmt.Println()
	fmt.Println("Main expense tracking flow would be implemented here.")
	fmt.Println("This is where the core functionality for tracking and analyzing expenses will go.")
}
