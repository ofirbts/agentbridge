package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourname/agentbridge/internal/output"
	"github.com/yourname/agentbridge/internal/run"
	"github.com/yourname/agentbridge/internal/workflow"
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
			wf := &workflow.Result{
				Status:     record.Status,
				Steps:      record.Steps,
				Errors:     record.Errors,
				Retries:    record.Retries,
				DurationMS: record.DurationMS,
				RunID:      record.ID,
				Crawler:    record.Crawler,
				MCP:        record.MCP,
				Result:     record.Result,
			}
			_, err := cmd.OutOrStdout().Write([]byte(output.FormatMarkdown(wf)))
			return err
		}

		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(record)
	},
}

func init() {
	inspectCmd.Flags().String("store", ".agentbridge/runs", "run store directory")
	inspectCmd.Flags().String("format", "json", "output format (json, markdown)")
	rootCmd.AddCommand(inspectCmd)
}
