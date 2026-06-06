package cmd

import (
	"fmt"
	"strings"

	internalBlame "com.lisitede.ai/internal/blame"
	"github.com/spf13/cobra"
)

var blameRouter Router

// blameFlags 用于 blame 命令的 flags
var blameFlags struct {
	title   string
	context string
}

var blameCmd = &cobra.Command{
	Use:   "blame <path>",
	Short: "Create a blame on a plan.",
	Long: `Create a blame on a plan.

Examples:
  ai blame /cms-mgr/PLAN-102 --title "Problem" --context "Full Description"`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var method string
		var path string
		if len(args) == 2 {
			method = strings.ToUpper(args[0])
			path = args[1]
		} else {
			method = ""
			path = args[0]
		}

		// 构建 args 映射
		argsMap := map[string]string{}
		if blameFlags.title != "" {
			argsMap["title"] = blameFlags.title
		}
		if blameFlags.context != "" {
			argsMap["context"] = blameFlags.context
		}

		// 调用路由
		result, err := blameRouter.Invoke(method, path, argsMap)
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	blameRouter.Verb("", "/:appId/:planId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["title"] == "" {
			return nil, fmt.Errorf("blame /:appId/:planId requires --title flag")
		}
		if args["context"] == "" {
			return nil, fmt.Errorf("blame /:appId/:planId requires --context flag")
		}
		return internalBlame.CreateBlame(params["appId"], params["planId"], args["title"], args["context"])
	})

	blameCmd.Flags().SortFlags = false
	blameCmd.Flags().StringVar(&blameFlags.title, "title", "", "blame title (required)")
	blameCmd.Flags().StringVar(&blameFlags.context, "context", "", "blame context, or path to a context markdown file (required)")

	rootCmd.AddCommand(blameCmd)
}
