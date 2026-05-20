package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ofirbts/agentbridge/internal/workflow"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "explain execution plan for a task without running it",
	RunE: func(cmd *cobra.Command, args []string) error {
		task, _ := cmd.Flags().GetString("task")
		if task == "" && len(args) > 0 {
			task = args[0]
		}
		if task == "" {
			return fmt.Errorf("task is required (--task or positional arg)")
		}
		mode, _ := cmd.Flags().GetString("mode")
		crawler, _ := cmd.Flags().GetString("crawler")
		mcpEndpoint, _ := cmd.Flags().GetString("mcp-endpoint")
		mcpTransport, _ := cmd.Flags().GetString("mcp-transport")
		plan := workflow.ExplainPlan(task, mode, crawler, mcpEndpoint, mcpTransport)

		format, _ := cmd.Flags().GetString("format")
		if format == "text" {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "task=%s mode=%s crawler=%s steps=%d behavior=%s\n",
				plan.Task, plan.Mode, plan.Crawler, len(plan.Steps), plan.Behavior)
			return err
		}

		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(plan)
	},
}

func init() {
	explainCmd.Flags().String("task", "", "task to explain")
	explainCmd.Flags().StringP("mode", "m", "normal", "execution mode (normal, deterministic)")
	explainCmd.Flags().String("crawler", "mock", "crawler provider (mock, http)")
	explainCmd.Flags().String("mcp-endpoint", "", "MCP server endpoint URL")
	explainCmd.Flags().String("mcp-transport", "", "MCP transport (http, stub)")
	explainCmd.Flags().String("format", "json", "output format (json, text)")
	rootCmd.AddCommand(explainCmd)
}
