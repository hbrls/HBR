package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI CLI.",
	Long: `AI CLI.

Available commands:
  hasshin   Identify the calling shell or environment
  blame     Create a blame on a plan
  get       Get a resource`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Parent() == nil {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", cmd.Long)
			fmt.Fprintf(cmd.OutOrStdout(), "Flags:\n%s\n", cmd.LocalFlags().FlagUsages())
		} else {
			defaultHelp(cmd, args)
		}
	})
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// outputError outputs error in human or JSON format
func outputError(err error) error {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	// 返回 nil 避免 cobra 再打印一次相同的错误
	return nil
}

// printResult 输出结果，支持 string 直接打印和 struct JSON 输出
func printResult(data interface{}) {
	if data == nil {
		fmt.Println("OK")
		return
	}

	// 如果是 string 类型，直接输出
	if s, ok := data.(string); ok {
		fmt.Println(s)
		return
	}

	// 其他类型走 JSON 序列化
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%v\n", data)
		return
	}
	fmt.Println(string(jsonBytes))
}
