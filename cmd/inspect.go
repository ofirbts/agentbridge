package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ofirbts/agentbridge/internal/output"
	"github.com/ofirbts/agentbridge/internal/run"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [run-id]",
	Short: "inspect a previous run by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		storePath, _ := cmd.Flags().GetString("store")
		format, _ := cmd.Flags().GetString("format")
		store := run.NewRunStore(storePath)
		record, err := store.Get(args[0])
		if err != nil {
			return fmt.Errorf("inspect: %w", err)
		}

		if format == "markdown" {
			_, err := cmd.OutOrStdout().Write([]byte(output.FormatMarkdown(record)))
			return err
		}

		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(output.FormatInspect(record))
	},
}

func init() {
	inspectCmd.Flags().String("store", ".agentbridge/runs", "run store directory")
	inspectCmd.Flags().String("format", "json", "output format (json, markdown)")
	rootCmd.AddCommand(inspectCmd)
}
