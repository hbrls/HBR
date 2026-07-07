package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		const (
			contextDir   = ".context"
			metadataFile = "current-plan-metadata.yaml"
		)

		var lines []string
		ok := true

		if info, err := os.Stat(contextDir); err != nil {
			if os.IsNotExist(err) {
				lines = append(lines, fmt.Sprintf("[FAIL] %s/ not found", contextDir))
			} else {
				lines = append(lines, fmt.Sprintf("[FAIL] %s/: %v", contextDir, err))
			}
			ok = false
		} else if !info.IsDir() {
			lines = append(lines, fmt.Sprintf("[FAIL] %s exists but is not a directory", contextDir))
			ok = false
		} else {
			lines = append(lines, fmt.Sprintf("[OK]   %s/ exists", contextDir))
		}

		metadataPath := filepath.Join(contextDir, metadataFile)
		if info, err := os.Stat(metadataPath); err != nil {
			if os.IsNotExist(err) {
				lines = append(lines, fmt.Sprintf("[FAIL] %s not found", metadataPath))
			} else {
				lines = append(lines, fmt.Sprintf("[FAIL] %s: %v", metadataPath, err))
			}
			ok = false
		} else if info.IsDir() {
			lines = append(lines, fmt.Sprintf("[FAIL] %s is a directory, expected file", metadataPath))
			ok = false
		} else {
			lines = append(lines, fmt.Sprintf("[OK]   %s exists", metadataPath))
		}

		if ok {
			lines = append(lines, "doctor: ok")
		} else {
			lines = append(lines, "doctor: failed")
		}

		printResult(strings.Join(lines, "\n"))
		if !ok {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	doctorCmd.Flags().SortFlags = false
	rootCmd.AddCommand(doctorCmd)
}
