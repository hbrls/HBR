package cmd

import (
	"fmt"
	"strings"

	internalUse "com.lisitede.ai/internal/use"
	"github.com/spf13/cobra"
)

var useRouter Router

// useFlags 用于 use 命令的 flags
var useFlags struct {
	title   string
	context string
}

var useCmd = &cobra.Command{
	Use:   "use <path>",
	Short: "Create a use on a plan.",
	Long: `Create a use on a plan.

Examples:
  ai use /cms-mgr/PLAN-102 --title "Problem" --context "Full Description"`,
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
		if useFlags.title != "" {
			argsMap["title"] = useFlags.title
		}
		if useFlags.context != "" {
			argsMap["context"] = useFlags.context
		}

		// 调用路由
		result, err := useRouter.Invoke(method, path, argsMap)
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	useRouter.Verb("", "/:appId/:planId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["title"] == "" {
			return nil, fmt.Errorf("use /:appId/:planId requires --title flag")
		}
		if args["context"] == "" {
			return nil, fmt.Errorf("use /:appId/:planId requires --context flag")
		}
		return internalUse.CreateUse(params["appId"], params["planId"], args["title"], args["context"])
	})

	useCmd.Flags().SortFlags = false
	useCmd.Flags().StringVar(&useFlags.title, "title", "", "use title (required)")
	useCmd.Flags().StringVar(&useFlags.context, "context", "", "use context, or path to a context markdown file (required)")

	rootCmd.AddCommand(useCmd)
}
