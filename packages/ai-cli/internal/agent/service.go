package agent

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// Hasshin 采集调用方的运行环境特征（Phase 1 信号探索版）。
//
// **平台范围**：Phase 1 聚焦 macOS + Windows；Windows chain 块通过 gopsutil
// 实现完整；runtime 块的 macOS-专属字段留零值、process 块的 uid/euid 由 Go
// stdlib 在 Windows 属性上自然返回 -1。Linux chain 返回 NotImplementedException。
//
// 探索期策略：stderr 按块输出多维观测信号，用于在多种调用环境下采样分析；
// 不做分类，只做特征抽取（见 .context/current-task.md 决策 #7）。stdout
// 暂沿用老 shell 白名单识别（命中 bash/zsh/sh/fish 返回名字，否则 "unknown"），
// 便于现有回归命令可用；归类模型待样本充足后再设计。
//
// 信号分块（各块 schema 定义详见 internal/agent/types.go；当前
// RuntimeBlock / ProcessBlock / ChainBlockLayers / EnvBlockLayers 已落地）：
//
//	runtime  (宿主环境)
//	process  (本进程)
//	chain    (祖先链快照；语义停止条件终止，无硬层数上限)
//	env      (按 family 分组 + uncategorized + raw 全量；TASK-103 旧 D/E 二分已合并)
//
// stdio 信号不再作为独立块，改由特征指纹 FP-01 承载（见 .context/current-task.md 第八节）。
func Hasshin(logLevel LogLevel) {
	rt := dumpRuntimeBlock(logLevel)
	_ = dumpProcessBlock(logLevel)
	chainBlock := dumpChainBlock(logLevel)
	_ = dumpEnvBlock(logLevel)

	// ========== Final ==========
	// fmt.Fprintln(os.Stderr, "[HASSHIN] ---- Final ----")

	chain := normalizeChainExe(chainBlock.Layers)
	fmt.Fprint(os.Stderr, renderHasshinPrompt(rt, chain))
}

// dumpRuntimeBlock 采集并输出 runtime 块（宿主环境快照）：os / arch /
// hostname / cpu / memory / OS 版本信息 / linux 提示 / virtualization。
// 结果写入 stderr，供信号探索采样分析使用。
//
// 平台分派：mac 侧通过 sw_vers / uname / sysctl 采集
// （见 collectDarwinRuntime）；linux / windows 的平台代码留待
// TASK-105 或 Phase 2。在非 mac 平台上运行时，mac 侧专属字段保持
// 零值 / 空串。
func dumpRuntimeBlock(logLevel LogLevel) RuntimeBlock {
	b := RuntimeBlock{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		CPUCount: runtime.NumCPU(),
	}
	if h, err := os.Hostname(); err == nil {
		b.Hostname = h
	} else {
		b.Hostname = fmt.Sprintf("<err: %v>", err)
	}

	switch runtime.GOOS {
	case "darwin":
		collectDarwinRuntime(&b)
	}

	if logLevel >= LogLevelDebug {
		printBlockByReflect(b, logLevel, "Runtime")
	}

	return b
}

// collectDarwinRuntime 填充 RuntimeBlock 的 mac 侧字段，通过 sw_vers / uname /
// sysctl 采集。字符串字段采集失败时写入 "<err: ...>"；整数字段采集失败留零。
func collectDarwinRuntime(b *RuntimeBlock) {
	setExecStr(&b.OSVersion, "sw_vers", "-productVersion")
	setExecStr(&b.OSBuild, "sw_vers", "-buildVersion")
	setExecStr(&b.KernelVersion, "uname", "-r")
	setExecStr(&b.CPUModel, "sysctl", "-n", "machdep.cpu.brand_string")
	setExecInt(&b.CPUCoresPhysical, "sysctl", "-n", "hw.physicalcpu")
	setExecInt64(&b.MemoryBytes, "sysctl", "-n", "hw.memsize")
}

// setExecStr 运行命令并把 trim 后的 stdout 写入 dst；失败时写入 "<err: %v>"。
func setExecStr(dst *string, name string, args ...string) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		*dst = fmt.Sprintf("<err: %v>", err)
		return
	}
	*dst = strings.TrimSpace(string(out))
}

// setExecInt 运行命令并把 trim 后的 stdout 解成十进制 int 写入 dst；失败保持零值。
func setExecInt(dst *int, name string, args ...string) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return
	}
	if n, perr := strconv.Atoi(strings.TrimSpace(string(out))); perr == nil {
		*dst = n
	}
}

// setExecInt64 同 setExecInt，dst 为 *int64。
func setExecInt64(dst *int64, name string, args ...string) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return
	}
	if n, perr := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64); perr == nil {
		*dst = n
	}
}

// dumpProcessBlock 采集并输出 A 块（process 分组，本进程运行身份）：
// pid / ppid / uid / euid / username / cwd / argv。结果写入 stderr，
// 供信号探索采样分析使用。
//
// 平台分派：当前仅填充跨平台 stdlib 能给出的字段；Unix-only 数值 ID
// （pgid / sid / gid / egid）及 Windows 专属字段（SID / SessionID /
// domain / integrity / elevated）留待 TASK-105 / Phase 2 填充。
func dumpProcessBlock(logLevel LogLevel) ProcessBlock {
	b := ProcessBlock{
		PID:  os.Getpid(),
		PPID: os.Getppid(),
		UID:  os.Getuid(),
		EUID: os.Geteuid(),
		Argv: os.Args,
	}
	if u, err := user.Current(); err == nil {
		b.Username = u.Username
	} else {
		b.Username = fmt.Sprintf("<err: %v>", err)
	}
	if cwd, err := os.Getwd(); err == nil {
		b.CWD = cwd
	} else {
		b.CWD = fmt.Sprintf("<err: %v>", err)
	}

	collectPlatformProcess(&b)

	if logLevel >= LogLevelDebug {
		printBlockByReflect(b, logLevel, "Process")
	}

	return b
}

// maxChainDepth 是 chain 遍历的死循环保护上限，**不是业务约束**。
// 正常链通常 3~10 层；触达上限以 Termination="depth_cap" 停止。
const maxChainDepth = 64

// dumpChainBlock 采集 C 块（chain，祖先链快照）并写入 stderr。
//
// 策略：从 os.Getppid() 起沿 ppid 上爬；以语义停止条件终止
// （见 ChainBlockLayers doc）。**本函数只采集 ps 直接返回的原始信息**
// （原始字段 + RawPsLine 全量保留），不做角色分类 / 指纹匹配等
// 二次判定 —— 那些逻辑归特征指纹章节处理（见 .context/current-task.md
// 第八节）。
//
// 平台实现：通过 Go build constraints（//go:build 指令）按 OS 分离，
// 见 chain_darwin.go / chain_windows.go / chain_other.go。
// 编译时按目标平台自动选择对应文件，无需运行时 switch。
//
// 返回值 parentComm 是 layers[0].Comm（即 os.Getppid() 的 ps -o comm=
// 原值，可能带 "-" 前缀），供 Hasshin final whitelist 段沿用老逻辑。
// 链为空时返回 ""。
func dumpChainBlock(logLevel LogLevel) ChainBlockLayers {
	b := ChainBlockLayers{}
	visited := make(map[int]bool)
	cur := os.Getppid()

	for len(b.Layers) < maxChainDepth {
		if cur <= 1 {
			if cur == 1 {
				// 把 init/launchd 这一终止层也记进来（成功与否都记，不丢信息）
				layer, _ := queryChainLayerOS(cur)
				b.Layers = append(b.Layers, layer)
				b.Termination = "init"
			} else {
				// cur == 0：极罕见的进程失孤
				b.Termination = "orphan"
			}
			break
		}
		if visited[cur] {
			b.Termination = "cycle"
			break
		}
		visited[cur] = true

		layer, err := queryChainLayerOS(cur)
		if err != nil {
			// 失败层暂不记录
			// 失败层也记录进来：保留 pid + 错误详情，不丢信息
			// b.Layers = append(b.Layers, layer)
			b.Termination = "query_failed"
			break
		}
		b.Layers = append(b.Layers, layer)

		cur = layer.PPID
		if cur == 0 {
			b.Termination = "orphan"
			break
		}
	}
	if b.Termination == "" {
		b.Termination = "depth_cap"
	}
	b.Depth = len(b.Layers)

	if logLevel >= LogLevelDebug {
		for i, layer := range b.Layers {
			printBlockByReflect(layer, logLevel, fmt.Sprintf("Chain - %d", i+1))
		}
	}
	return b
}

// envExactMap 固定键 → family 名的精确匹配表。优先级高于 envPrefixRules。
//
// 特异 key 放 exact（如 VSCODE_CRASH_REPORTER_PROCESS_TYPE 归 ExtensionHost，
// 而非让前缀匹配吃进 IDE）；同家族的通用 key 走 prefix。
var envExactMap = map[string]string{
	// ShellVersion
	"BASH_VERSION": "ShellVersion",
	"ZSH_VERSION":  "ShellVersion",
	"FISH_VERSION": "ShellVersion",

	// Terminal
	"TERM":                 "Terminal",
	"TERM_PROGRAM":         "Terminal",
	"TERM_PROGRAM_VERSION": "Terminal",
	"COLORTERM":            "Terminal",
	"LC_TERMINAL":          "Terminal",
	"KITTY_WINDOW_ID":      "Terminal",
	"WEZTERM_EXECUTABLE":   "Terminal",
	"ALACRITTY_SOCKET":     "Terminal",
	"TERMINAL_EMULATOR":    "Terminal",

	// Session（比 Terminal 特异：有这些 key 基本能确定是 session manager/remote）
	"SSH_CONNECTION": "Session",
	"SSH_CLIENT":     "Session",
	"SSH_TTY":        "Session",
	"TMUX":           "Session",
	"STY":            "Session",

	// CI
	"CI":             "CI",
	"GITHUB_ACTIONS": "CI",
	"GITLAB_CI":      "CI",
	"CIRCLECI":       "CI",
	"BUILDKITE":      "CI",
	"JENKINS_URL":    "CI",
	"TRAVIS":         "CI",

	// Remote
	"KUBERNETES_SERVICE_HOST": "Remote",
	"IN_NIX_SHELL":            "Remote",

	// ExtensionHost（VSCODE_ 子集，特异性高于 IDE 前缀）
	"VSCODE_CRASH_REPORTER_PROCESS_TYPE": "ExtensionHost",
	"VSCODE_PID":                         "ExtensionHost",
	"VSCODE_IPC_HOOK":                    "ExtensionHost",
	"VSCODE_ESM_ENTRYPOINT":              "ExtensionHost",
	"VSCODE_NLS_CONFIG":                  "ExtensionHost",

	// Locale（LC_TERMINAL 已在上面走 Terminal）
	"LANG": "Locale",

	// Editor
	"EDITOR": "Editor",
	"VISUAL": "Editor",
	"PAGER":  "Editor",

	// User (Unix) — 用户身份相关
	"USER":    "User",
	"LOGNAME": "User",

	// User (Windows) — 用户身份相关
	"USERDOMAIN": "User",
}

// envPrefixRules 前缀匹配规则，按顺序逐条试；exact 不命中才查这里。
// 更窄的前缀应在前（目前 IDE / Locale 无冲突）。
var envPrefixRules = []struct {
	prefix string
	family string
}{
	// IDE 前缀（VSCODE_ 的特异 key 已在 envExactMap 归 ExtensionHost）
	{"VSCODE_", "IDE"},
	{"WINDSURF_", "IDE"},
	{"CURSOR_", "IDE"},
	{"CASCADE_", "IDE"},
	{"JETBRAINS_", "IDE"},
	{"IDEA_", "IDE"},
	{"PYCHARM_", "IDE"},

	// Locale（LC_TERMINAL 已在 envExactMap 归 Terminal）
	{"LC_", "Locale"},
}

// classifyEnvKey 将 env key 归类到某个 family；归不到返回空串（→ Uncategorized）。
// exact 优先 prefix；per-key 决策，每个 key 最多一个 family（单归属）。
//
// Windows 适配：Windows 的 env key 大小写不敏感（`Path` == `PATH` == `path`），
// 但 os.Environ() 返回原始大小写。策略：
//  1. 先按原始 key 试匹配（mac 行为此步即终结，与此前完全一致）
//  2. Windows 下再按 UPPER(key) 再试一次（捕捉 Path / Home 等 native 形式）
//
// mac 不执行第 2 步，原有行为严格不变。
func classifyEnvKey(key string) string {
	// Step 1: 先按原始 key 匹配（mac / Windows 同走此分支）
	if fam, ok := envExactMap[key]; ok {
		return fam
	}
	for _, r := range envPrefixRules {
		if strings.HasPrefix(key, r.prefix) {
			return r.family
		}
	}

	// Step 2: Windows 专用 fallback，按 UPPER(key) 再试一次
	if runtime.GOOS == "windows" {
		upper := strings.ToUpper(key)
		if upper != key { // 已是 UPPER 的 key 第一步已试过，跳过重复查
			if fam, ok := envExactMap[upper]; ok {
				return fam
			}
			for _, r := range envPrefixRules {
				if strings.HasPrefix(upper, r.prefix) {
					return r.family
				}
			}
		}
	}

	return ""
}

// mapForFamily 返回 EnvBlockLayers 里指定 family 对应的 map 指针。
// family 名必须与 envExactMap / envPrefixRules 的 family 字段一致；
// 不识别的返回 Uncategorized 兜底。
func (b *EnvBlockLayers) mapForFamily(family string) *map[string]string {
	switch family {
	case "ExtensionHost":
		return &b.ExtensionHost
	case "IDE":
		return &b.IDE
	case "Session":
		return &b.Session
	case "CI":
		return &b.CI
	case "Remote":
		return &b.Remote
	case "Terminal":
		return &b.Terminal
	case "ShellVersion":
		return &b.ShellVersion
	case "Editor":
		return &b.Editor
	case "Locale":
		return &b.Locale
	case "User":
		return &b.User
	case "ShellPref":
		return &b.ShellPref
	default:
		return &b.Uncategorized
	}
}

// dumpEnvBlock 采集 env 块并写入 stderr。
//
// 策略：
//   - Raw 保留 os.Environ() 原样（顺序 + 重复），不丢信息的第一道防线
//   - 每个 env 走 classifyEnvKey 单归属到 family，归不到进 Uncategorized
//   - stderr 按 family 顺序打（空 family 跳过 header）
func dumpEnvBlock(logLevel LogLevel) EnvBlockLayers {
	b := EnvBlockLayers{
		Raw: os.Environ(),
	}
	b.Platform.Path = os.Getenv("PATH")
	b.Platform.HOME = os.Getenv("HOME")
	b.Platform.USERPROFILE = os.Getenv("USERPROFILE")
	b.Platform.TMPDIR = os.Getenv("TMPDIR")
	b.Platform.TEMP = os.Getenv("TEMP")
	b.Platform.TMP = os.Getenv("TMP")
	b.Platform.APPDATA = os.Getenv("APPDATA")
	b.Platform.LOCALAPPDATA = os.Getenv("LOCALAPPDATA")
	b.Shell.SHELL = os.Getenv("SHELL")
	b.Shell.MSYS_NO_PATHCONV = os.Getenv("MSYS_NO_PATHCONV")
	b.Shell.PSModulePath = os.Getenv("PSModulePath")
	b.KiloCode.KILO_PLATFORM = os.Getenv("KILO_PLATFORM")
	b.KiloCode.KILO_APP_NAME = os.Getenv("KILO_APP_NAME")
	b.KiloCode.KILO_APP_VERSION = os.Getenv("KILO_APP_VERSION")
	b.KiloCode.KILO_DISABLE_AUTOUPDATE = os.Getenv("KILO_DISABLE_AUTOUPDATE")
	b.KiloCode.KILO_DISABLE_CLAUDE_CODE = os.Getenv("KILO_DISABLE_CLAUDE_CODE")
	b.KiloCode.KILO_ENABLE_QUESTION_TOOL = os.Getenv("KILO_ENABLE_QUESTION_TOOL")
	b.KiloCode.KILO_EDITOR_NAME = os.Getenv("KILO_EDITOR_NAME")
	b.KiloCode.KILOCODE_EDITOR_NAME = os.Getenv("KILOCODE_EDITOR_NAME")
	b.KiloCode.KILOCODE_VERSION = os.Getenv("KILOCODE_VERSION")
	b.KiloCode.KILOCODE_FEATURE = os.Getenv("KILOCODE_FEATURE")
	b.KiloCode.KILO_CLIENT = os.Getenv("KILO_CLIENT")
	b.KiloCode.KILO_VSCODE_VERSION = os.Getenv("KILO_VSCODE_VERSION")
	b.KiloCode.KILO_TELEMETRY_LEVEL = os.Getenv("KILO_TELEMETRY_LEVEL")
	b.KiloCode.KILO_MACHINE_ID = os.Getenv("KILO_MACHINE_ID")
	b.KiloCode.KILO_PID = os.Getenv("KILO_PID")
	b.KiloCode.KILO_SERVER_PASSWORD = os.Getenv("KILO_SERVER_PASSWORD")
	b.Flutter.PUB_HOSTED_URL = os.Getenv("PUB_HOSTED_URL")
	b.Flutter.FLUTTER_STORAGE_BASE_URL = os.Getenv("FLUTTER_STORAGE_BASE_URL")

	// skipKeys 内的所有 key 均已大写化，直接与 upperKey 匹配
	skipKeys := map[string]struct{}{
		// Platform
		"PATH": {}, "HOME": {}, "USERPROFILE": {}, "TMPDIR": {}, "TEMP": {}, "TMP": {}, "APPDATA": {}, "LOCALAPPDATA": {},
		// Shell
		"SHELL": {}, "MSYS_NO_PATHCONV": {}, "PSMODULEPATH": {},
		// KiloCode
		"KILO_PLATFORM": {}, "KILO_APP_NAME": {}, "KILO_APP_VERSION": {}, "KILO_DISABLE_AUTOUPDATE": {},
		"KILO_DISABLE_CLAUDE_CODE": {}, "KILO_ENABLE_QUESTION_TOOL": {},
		"KILO_EDITOR_NAME": {}, "KILOCODE_EDITOR_NAME": {}, "KILOCODE_VERSION": {}, "KILOCODE_FEATURE": {},
		"KILO_CLIENT": {}, "KILO_VSCODE_VERSION": {}, "KILO_TELEMETRY_LEVEL": {},
		"KILO_MACHINE_ID": {}, "KILO_PID": {}, "KILO_SERVER_PASSWORD": {},
		// Flutter
		"PUB_HOSTED_URL": {}, "FLUTTER_STORAGE_BASE_URL": {},
	}

	for _, kv := range b.Raw {
		var key, val string
		if eq := strings.IndexByte(kv, '='); eq < 0 {
			key = kv
			val = ""
		} else {
			key = kv[:eq]
			val = kv[eq+1:]
		}

		// Platform / Shell / KiloCode / Flutter 字段已单独处理，不归入其他 family
		upperKey := strings.ToUpper(key)
		if _, exists := skipKeys[upperKey]; exists {
			continue
		}

		target := b.mapForFamily(classifyEnvKey(key))
		if *target == nil {
			*target = make(map[string]string)
		}
		(*target)[key] = val
	}

	if logLevel >= LogLevelDebug {
		printBlockByReflect(b.Platform, logLevel, "Env - Platform")
		printBlockByReflect(b.Shell, logLevel, "Env - Shell")
		printBlockByReflect(b.KiloCode, logLevel, "Env - KiloCode")
		printBlockByReflect(b.Flutter, logLevel, "Env - Flutter")
	}

	if logLevel >= LogLevelDebug {
		fmt.Fprintf(os.Stderr, "---- Env ----\n")
		printEnvFamily("extension_host", b.ExtensionHost)
		printEnvFamily("ide", b.IDE)
		printEnvFamily("session", b.Session)
		printEnvFamily("ci", b.CI)
		printEnvFamily("remote", b.Remote)
		printEnvFamily("terminal", b.Terminal)
		printEnvFamily("shell_version", b.ShellVersion)
		printEnvFamily("editor", b.Editor)
		printEnvFamily("locale", b.Locale)
		printEnvFamily("user", b.User)
		printEnvFamily("shell_pref", b.ShellPref)
		printEnvFamily("uncategorized", b.Uncategorized)
	}
	return b
}

func printEnvFamily(label string, m map[string]string) {
	if len(m) == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "  [%s]\n", label)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(os.Stderr, "    %s\n", k)
	}
}
