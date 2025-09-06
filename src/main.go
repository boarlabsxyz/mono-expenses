package main

import (
	"os"
)

const version = "1.0.0"

const minArgsForCommand = 2

func main() {
	if len(os.Args) < minArgsForCommand {
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
		_, _ = os.Stderr.WriteString("Unknown command: " + command + "\n")
		_, _ = os.Stderr.WriteString("Use 'mono-expenses help' to see available commands.\n")
		os.Exit(1)
	}
}

func showHelp() {
	_, _ = os.Stdout.WriteString("mono-expenses - Quick and dirty solution for tracking and analysing expenses using AI\n")
	_, _ = os.Stdout.WriteString("\n")
	_, _ = os.Stdout.WriteString("USAGE:\n")
	_, _ = os.Stdout.WriteString("    mono-expenses [COMMAND]\n")
	_, _ = os.Stdout.WriteString("\n")
	_, _ = os.Stdout.WriteString("COMMANDS:\n")
	_, _ = os.Stdout.WriteString("    help      Show this help message\n")
	_, _ = os.Stdout.WriteString("    version   Show version information\n")
	_, _ = os.Stdout.WriteString("\n")
	_, _ = os.Stdout.WriteString("When no command is provided, the main expense tracking flow is executed.\n")
}

func showVersion() {
	_, _ = os.Stdout.WriteString("mono-expenses " + version + "\n")
}

func runMainFlow() {
	_, _ = os.Stdout.WriteString("mono-expenses - Expense Tracking\n")
	_, _ = os.Stdout.WriteString("=============================\n")
	_, _ = os.Stdout.WriteString("\n")
	_, _ = os.Stdout.WriteString("Main expense tracking flow would be implemented here.\n")
	_, _ = os.Stdout.WriteString("This is where the core functionality for tracking and analyzing expenses will go.\n")
}
