package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourname/agentbridge/internal/failure"
)

var simulateFailureCmd = &cobra.Command{
	Use:   "simulate-failure",
	Short: "inject simulated failures for reliability testing",
	RunE: func(cmd *cobra.Command, args []string) error {
		kind, _ := cmd.Flags().GetString("kind")
		seed, _ := cmd.Flags().GetInt64("seed")
		count, _ := cmd.Flags().GetInt("count")
		if count < 1 {
			count = 1
		}

		sim := failure.NewSimulator(seed)
		results := make([]map[string]any, 0, count)
		failures := 0

		for i := 0; i < count; i++ {
			err := sim.Inject(kind)
			entry := map[string]any{
				"attempt": i + 1,
				"kind":    kind,
				"success": err == nil,
			}
			if err != nil {
				entry["error"] = err.Error()
				failures++
				fmt.Fprintln(cmd.ErrOrStderr(), err.Error())
			}
			results = append(results, entry)
		}

		out := map[string]any{
			"kind":     kind,
			"count":    count,
			"failures": failures,
			"results":  results,
		}
		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(out)
	},
}

func init() {
	simulateFailureCmd.Flags().String("kind", "random", "failure kind (random, timeout, rate_limit, server_error)")
	simulateFailureCmd.Flags().Int64("seed", 0, "random seed for deterministic simulation")
	simulateFailureCmd.Flags().Int("count", 1, "number of simulations to run")
	rootCmd.AddCommand(simulateFailureCmd)
}
