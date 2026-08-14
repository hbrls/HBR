package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "backstage",
	Short:         "Backstage CLI.",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long: `Backstage CLI.

Available commands:
  scrum     Scrum workflow`,
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

// outputError propagates the error to the top-level command.
func outputError(err error) error {
	return err
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

	// 其他类型走 JSON 序列化，SetEscapeHTML(false) 避免 HTML 转义
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		fmt.Printf("%v\n", data)
		return
	}
	fmt.Print(buf.String())
}
