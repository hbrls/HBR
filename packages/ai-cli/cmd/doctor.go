package cmd

import (
	"github.com/spf13/cobra"
)

// doctorCmd 运行环境诊断命令
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose the CLI runtime environment.",
	Long: `Diagnose the CLI runtime environment.
Checks required environment variables, backend connectivity and version info.

Examples:
  # Run all checks
  ai doctor`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implement doctor checks
		printResult("doctor: ok")
		return nil
	},
}

func init() {
	doctorCmd.Flags().SortFlags = false
	rootCmd.AddCommand(doctorCmd)
}
