package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

// normalizeChainExe 把 chain 链各层的可执行身份规范化为一条去重后的
// 绝对路径列表，顺序与 layers 一致（从直接父进程向上）。
//
// 每层取值策略：
//   - 优先使用 layer.Comm：
//     darwin  KERN_PROCARGS2 的 exec_path，形如 "/Applications/.../Foo"；
//     windows QueryFullProcessImageNameW，形如 "C:\\Windows\\explorer.exe"。
//     两者都是完整绝对路径，含空格安全、不截断。filepath.IsAbs 同时覆盖
//     两套语义（darwin 识 `/`，windows 识 `X:\`）。
//   - Comm 非绝对（权限受限回退到 Name() 只给 basename 等场景）时，
//     回退到 parseExePath(layer.Command) 的旧启发式切分。
func normalizeChainExe(layers []ChainBlock) []string {
	var normalized []string
	seen := make(map[string]bool)

	for _, layer := range layers {
		var evalPath string
		if filepath.IsAbs(layer.Comm) {
			evalPath = normalizeExePath(layer.Comm)
		} else {
			evalPath = parseExePath(layer.Command)
		}
		if evalPath == "" {
			continue
		}
		if !seen[evalPath] {
			seen[evalPath] = true
			normalized = append(normalized, evalPath)
		}
	}

	// 身份归一：按 basename 查 knownExeAliases，命中则替换为 canonical；
	// 未命中保留原绝对路径。替换后再走一次 set 去重，避免同一家族的多层
	// （例如 Windsurf + Windsurf Helper）折叠后重复出现。
	aliased := make([]string, 0, len(normalized))
	dedup := make(map[string]struct{}, len(normalized))
	for _, p := range normalized {
		out := p
		// 查表键带前置 "/"，既能用作末段边界（"/Windsurf"）也能用作中间段
		// 匹配（"/iTermServer-" 命中 "/.../iTermServer-3.4.19"）。采用子串
		// 包含匹配，由 aliasRules 的长度降序保证"长 key 优先"，从而让更具体
		// 的规则（如 "/Windsurf Helper"）先于泛化规则（"/Windsurf"）命中。
		if canonical, ok := lookupExeAlias(p); ok {
			out = canonical
		}
		// 空串视为显式抑制（见 knownExeAliases 中 "launchd" 条目），
		// 不入集合、不输出。
		if out == "" {
			continue
		}
		if _, hit := dedup[out]; hit {
			continue
		}
		dedup[out] = struct{}{}
		aliased = append(aliased, out)
	}
	return aliased
}

// knownExeAliases 将链路上各可执行路径归一到 canonical 身份名。
// key 为路径片段（macOS 用 "/" 前缀、Windows 用 "\" 前缀），通过子串
// 包含匹配（strings.Contains）命中绝对路径；保留原始大小写与空格
// （如 "/Windsurf Helper"）。key 可以是纯 basename（如 "/Windsurf"）或
// 含目录片段（如 "\nodejs\node.exe"），视区分精度需要而定。匹配顺序由
// exeAliasRules 保证长 key 优先，避免泛化 key 遮蔽更具体的变体。
// value = canonical identity；空串 = 显式抑制，该层不入链路归一输出。
// 未命中时调用方保留原完整路径。
//
// 新增条目时请成对补齐同家族的常见变体（无后缀 / " Helper" / ".exe" 等），
// 避免同一工具因平台差异 / 子进程形态造成链路输出中的身份分裂。
var knownExeAliases = map[string]string{
	// macOS
	"/Windsurf":        "Windsurf",
	"/Windsurf Helper": "Windsurf",

	"/iTerm2":       "iTerm2",
	"/iTermServer-": "iTerm2", // 捕获 "/.../iTermServer-3.4.19" 等版本变体（子串匹配）

	"/launchd": "", // macOS 系统 PID 1，所有用户态进程的终极祖先，无身份区分价值

	// Windows
	`\Microsoft VS Code\Code.exe`: "VSCode",
	`\Windsurf.exe`:               "Windsurf",

	`\Android Studio\bin`: "AndroidStudio",

	`\PyCharm\bin\pycharm64.exe`: "PyCharm",

	`\Git\usr\bin\bash.exe`: "GitBash",
	`\Git\bin\bash.exe`:     "GitBash",

	`\powershell.exe`:                     "PowerShell",
	`\WindowsApps\Microsoft.PowerShell_7`: "PowserShell7",
	`\WindowsTerminal.exe`:                "WindowsTerminal",
	`\System32\cmd.exe`:                   "WindowsCmd",

	`\.vscode\extensions\kilocode`: "KiloCode-vsix",
	`\plugins\cline\core`:          "Cline-jbp",

	`\python.exe`:               "Python",
	`\.venv\Scripts\python.exe`: "Python-venv",

	`\nodejs\node.exe`:                        "Node",
	`\@kilocode\cli-windows-x64\bin\kilo.exe`: "Kilo-bin",

	`\explorer.exe`: "", // Windows 资源管理器，无身份区分价值
}

// exeAliasRule 是 knownExeAliases 展平后的一条匹配规则。
// 用 slice 而非 map 直接迭代，以便按 key 长度排序，保证"长 key 优先"
// 的确定性命中顺序（Go map 迭代顺序本身是随机的）。
type exeAliasRule struct {
	key       string
	canonical string
}

// exeAliasRules 在 init 中由 knownExeAliases 展平并按 len(key) 降序排列，
// 让更具体的 key（"/Windsurf Helper"）先于泛化 key（"/Windsurf"）被尝试。
var exeAliasRules []exeAliasRule

func init() {
	exeAliasRules = make([]exeAliasRule, 0, len(knownExeAliases))
	for k, v := range knownExeAliases {
		exeAliasRules = append(exeAliasRules, exeAliasRule{key: k, canonical: v})
	}
	sort.Slice(exeAliasRules, func(i, j int) bool {
		return len(exeAliasRules[i].key) > len(exeAliasRules[j].key)
	})
}

// lookupExeAlias 在 exeAliasRules 上做子串包含匹配；按长度降序，
// 首个命中即返回。未命中返回 ("", false)。
func lookupExeAlias(path string) (string, bool) {
	for _, r := range exeAliasRules {
		if strings.Contains(path, r.key) {
			return r.canonical, true
		}
	}
	return "", false
}

// parseExePath 从 Command 字符串中提取第一个可执行路径并标准化。
// 处理带引号的路径，解析符号链接、.. 和 .。
func parseExePath(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return ""
	}

	var exePath string
	if cmd[0] == '"' {
		if end := strings.Index(cmd[1:], "\""); end >= 0 {
			exePath = cmd[1 : end+1]
		}
	} else {
		if idx := strings.Index(cmd, " "); idx >= 0 {
			exePath = cmd[:idx]
		} else {
			exePath = cmd
		}
	}
	if exePath == "" {
		return ""
	}
	return normalizeExePath(exePath)
}

// normalizeExePath 将可执行路径做 Abs + EvalSymlinks + Clean 规范化；
// 任一步失败都容忍降级，保证始终返回非空可读路径。
func normalizeExePath(exePath string) string {
	if exePath == "" {
		return ""
	}
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		absPath = exePath
	}
	evalPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		evalPath = absPath
	}
	return filepath.Clean(evalPath)
}

func printBlockByReflect(block interface{}, logLevel LogLevel, blockName string) {
	fmt.Fprintf(os.Stderr, "---- %s ----\n", blockName)

	v := reflect.ValueOf(block)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("loglevel")

		fieldLevel := parseLogLevel(tag)
		if logLevel >= fieldLevel {
			fieldName := field.Name
			fieldValue := v.Field(i).Interface()
			fmt.Fprintf(os.Stderr, "  %s: %v\n", fieldName, fieldValue)
		}
	}
}

func parseLogLevel(s string) LogLevel {
	switch s {
	case "info":
		return LogLevelInfo
	case "debug":
		return LogLevelDebug
	case "dangerous":
		return LogLevelDangerous
	default:
		return LogLevelInfo
	}
}
