// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Mrinal Prakash

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MRINALPRAKASHFSD/reliable-job-engine/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version":
		fmt.Printf("jobctl version %s (commit: %s)\n", config.Version, config.Commit)
	case "status":
		fmt.Println("Job engine status: Phase 1 (Repository Foundation Initialized)")
		fmt.Println("Database schema available in migrations/")
	case "dlq":
		handleDLQSubcommand(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleDLQSubcommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: jobctl dlq <list|replay|purge> [flags]")
		os.Exit(1)
	}

	sub := args[0]
	switch sub {
	case "list":
		fmt.Println("DLQ inspection is planned for Phase 4.")
	case "replay":
		dlqCmd := flag.NewFlagSet("replay", flag.ExitOnError)
		jobID := dlqCmd.String("id", "", "ID of the dead letter job to replay")
		_ = dlqCmd.Parse(args[1:])
		if *jobID == "" {
			fmt.Println("Error: --id flag is required for replay")
			os.Exit(1)
		}
		fmt.Printf("Replaying DLQ job %s (Planned for Phase 4)\n", *jobID)
	case "purge":
		fmt.Println("DLQ purging is planned for Phase 4.")
	default:
		fmt.Fprintf(os.Stderr, "Unknown dlq subcommand: %s\n", sub)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("jobctl - Operational CLI for Reliable Job Engine")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  jobctl <command> [arguments]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  version     Print current version and build commit")
	fmt.Println("  status      Check cluster or local engine status")
	fmt.Println("  dlq         Inspect, replay, or purge dead-letter jobs")
	fmt.Println("  help        Show this help message")
}
