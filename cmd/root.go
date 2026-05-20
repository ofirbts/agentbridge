package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agentbridge",
	Short: "A reliable execution layer for AI agent web workflows",
	Long: `agentbridge is a lightweight CLI that makes web-aware AI agents predictable, observable, and production-aware.

Features:
- CLI-first
- Deterministic execution mode
- Retry + fallback system
- Structured logs and tracing
- Pluggable search / crawl / extract providers`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
