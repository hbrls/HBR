package cmd

import (
	"strings"

	"com.lisitede.backstage/internal/scrum"
	"github.com/spf13/cobra"
)

var scrumRouter Router
var scrumVision bool

var scrumCmd = &cobra.Command{
	Use:   "scrum <resource>",
	Short: "Scrum workflow.",
	Long: `Scrum workflow.

Examples:
  backstage scrum pull-task
  backstage scrum pull-task --vision`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var method string
		var raw string
		if len(args) == 2 {
			method = strings.ToUpper(args[0])
			raw = args[1]
		} else {
			method = ""
			raw = args[0]
		}

		pathname := raw
		if !strings.HasPrefix(pathname, "/") {
			pathname = "/" + pathname
		}

		argsMap := map[string]string{}
		if scrumVision {
			argsMap["vision"] = "true"
		}

		result, err := scrumRouter.Invoke(method, pathname, argsMap)
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	scrumRouter.Verb("", "/pull-task", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		filter := "TO DO"
		status := "IN DO"
		if args["vision"] == "true" {
			filter = "TO VISION"
			status = "IN VISION"
		}
		return scrum.PullTask(filter, status)
	})

	scrumCmd.Flags().BoolVar(&scrumVision, "vision", false, "pull from TO VISION and move to IN VISION")
	rootCmd.AddCommand(scrumCmd)
}
