package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var getRouter Router

var getDir string

var getCmd = &cobra.Command{
	Use:   "get <resource>",
	Short: "Get a resource.",
	Long: `Get a resource.

Examples:
  ai get first-task
  ai get first-task --dir .context/VISION-001
  ai get first-vision`,
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

		argsMap := map[string]string{
			"dir": getDir,
		}

		result, err := getRouter.Invoke(method, pathname, argsMap)
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	getRouter.Verb("", "/first-task", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		contextDir := args["dir"]
		if contextDir == "" {
			contextDir = ".context"
		}

		entries, err := os.ReadDir(contextDir)
		if err != nil {
			if os.IsNotExist(err) {
				return "No pending tasks.", nil
			}
			return nil, err
		}

		var tasks []string
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if strings.HasPrefix(name, "TASK-") && strings.HasSuffix(name, ".md") {
				tasks = append(tasks, name)
			}
		}

		if len(tasks) == 0 {
			return "No pending tasks.", nil
		}

		sort.Strings(tasks)
		return filepath.Join(contextDir, tasks[0]), nil
	})

	getRouter.Verb("", "/first-vision", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		data, err := os.ReadFile(filepath.Join(".context", "current-vision.yaml"))
		if err != nil {
			if os.IsNotExist(err) {
				return "No visions file.", nil
			}
			return nil, err
		}

		var vf struct {
			ActiveVisions       []string `yaml:"active-visions"`
			LastVision          string   `yaml:"last-vision"`
			LastVisionUpdatedAt string   `yaml:"last-vision-updated-at"`
		}
		if err := yaml.Unmarshal(data, &vf); err != nil {
			return nil, err
		}

		if len(vf.ActiveVisions) == 0 {
			return "No visions configured.", nil
		}

		next := ""
		if vf.LastVision == "" {
			next = vf.ActiveVisions[0]
		} else {
			for i, v := range vf.ActiveVisions {
				if v == vf.LastVision {
					next = vf.ActiveVisions[(i+1)%len(vf.ActiveVisions)]
					break
				}
			}
			if next == "" {
				next = vf.ActiveVisions[0]
			}
		}
		return filepath.Join(".context", next), nil
	})

	getCmd.Flags().SortFlags = false
	getCmd.Flags().StringVar(&getDir, "dir", "", "Directory to search (default: .context)")

	rootCmd.AddCommand(getCmd)
}
