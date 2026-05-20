package cmd

import (
	"encoding/json"

	"github.com/ofirbts/agentbridge/internal/output"
	"github.com/ofirbts/agentbridge/internal/workflow"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run an AI agent web task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task := args[0]
		mode, _ := cmd.Flags().GetString("mode")
		storePath, _ := cmd.Flags().GetString("store")
		search, _ := cmd.Flags().GetString("search")
		crawler, _ := cmd.Flags().GetString("crawler")
		mcpEndpoint, _ := cmd.Flags().GetString("mcp-endpoint")
		mcpTransport, _ := cmd.Flags().GetString("mcp-transport")

		engine, err := workflow.NewEngine(workflow.Config{
			Mode:         mode,
			StorePath:    storePath,
			Search:       search,
			Crawler:      crawler,
			MCPEndpoint:  mcpEndpoint,
			MCPTransport: mcpTransport,
		})
		if err != nil {
			return err
		}
		record, err := engine.RunTask(task)
		if err != nil && record == nil {
			return err
		}
		if record != nil && err != nil {
			record.Errors = append(record.Errors, err.Error())
			if record.Status != "failed" {
				record.Status = "failed"
			}
		}

		format, _ := cmd.Flags().GetString("format")
		if format == "markdown" {
			_, err := cmd.OutOrStdout().Write([]byte(output.FormatMarkdown(record)))
			return err
		}

		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(output.FormatRunResult(record))
	},
}

func init() {
	runCmd.Flags().StringP("mode", "m", "normal", "execution mode (normal, deterministic)")
	runCmd.Flags().String("store", ".agentbridge/runs", "run store directory")
	runCmd.Flags().String("search", "mock", "search provider (mock, http)")
	runCmd.Flags().String("crawler", "mock", "crawler provider (mock, http)")
	runCmd.Flags().String("mcp-endpoint", "", "MCP server endpoint URL")
	runCmd.Flags().String("mcp-transport", "", "MCP transport (http, stub); default http for http(s) endpoints")
	runCmd.Flags().String("format", "json", "output format (json, markdown)")
	rootCmd.AddCommand(runCmd)
}
