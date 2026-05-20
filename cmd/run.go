package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/yourname/agentbridge/internal/output"
	"github.com/yourname/agentbridge/internal/workflow"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run an AI agent web task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task := args[0]
		cfg, _ := cmd.Flags().GetString("config")
		mode, _ := cmd.Flags().GetString("mode")
		storePath, _ := cmd.Flags().GetString("store")
		crawler, _ := cmd.Flags().GetString("crawler")
		mcpEndpoint, _ := cmd.Flags().GetString("mcp-endpoint")

		engine := workflow.NewEngine(workflow.Config{
			Mode:        mode,
			ConfigFile:  cfg,
			StorePath:   storePath,
			Crawler:     crawler,
			MCPEndpoint: mcpEndpoint,
		})
		result, err := engine.RunTask(task)
		if err != nil && result == nil {
			return err
		}
		if result != nil && err != nil {
			result.Errors = append(result.Errors, err.Error())
			if result.Status != "failed" {
				result.Status = "failed"
			}
		}

		format, _ := cmd.Flags().GetString("format")
		if format == "markdown" {
			_, err := cmd.OutOrStdout().Write([]byte(output.FormatMarkdown(result)))
			return err
		}

		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(output.FormatRunResult(result))
	},
}

func init() {
	runCmd.Flags().StringP("mode", "m", "normal", "execution mode (normal, deterministic)")
	runCmd.Flags().String("config", "", "config file path (JSON)")
	runCmd.Flags().String("store", ".agentbridge/runs", "run store directory")
	runCmd.Flags().String("crawler", "mock", "crawler provider (mock, http)")
	runCmd.Flags().String("mcp-endpoint", "", "MCP server endpoint (stub client)")
	runCmd.Flags().String("format", "json", "output format (json, markdown)")
	rootCmd.AddCommand(runCmd)
}
