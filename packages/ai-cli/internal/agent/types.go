package agent

import "time"

type LogLevel int

const (
	LogLevelInfo LogLevel = iota
	LogLevelDebug
	LogLevelDangerous
)

// RuntimeBlock 是 Schema v1 中 runtime 分组的 Go 映射：
// 描述本 CLI 被执行时的宿主机器 / OS 上下文快照。
//
// 采集失败的字符串字段沿用 "<err: ...>" 字串风格（见 dumpRuntimeBlock）；
// 整数字段采集失败保持零值（0 在本 block 的整数字段上天然非合法值，
// 等价于哨位）。错误信息保留到 stderr 诊断，便于样本分析时定位原因。
//
// 平台策略：采用 macOS ∪ Windows 并集 schema（同 ProcessBlock）；字段
// 只在一边有语义时，另一边留零值或哨位（字符串空串、int 零）。Linux
// 不作为一等目标——合辙则顺带可用，不合辙留空。
//
// Scope：仅记录"被探测环境"；不含本 CLI binary 自身的构建身份
// （Go 版本、build target、cliVersion 等），那些属于"探测器自己"，
// 若未来需要会走独立 block 承载。
type RuntimeBlock struct {
	// Platform identity. (both)
	OS   string `loglevel:"debug"` // runtime.GOOS：如 "darwin" / "windows" / "linux"
	Arch string `loglevel:"debug"` // runtime.GOARCH：如 "arm64" / "amd64"

	// Machine identity. (both)
	Hostname string `loglevel:"dangerous"` // os.Hostname()

	// CPU. (both)
	CPUCount         int    `loglevel:"dangerous"` // runtime.NumCPU()：进程可用的 logical CPU 数
	CPUCoresPhysical int    `loglevel:"dangerous"` // 物理核心数；mac: sysctl hw.physicalcpu
	CPUModel         string `loglevel:"dangerous"` // 型号字串；mac: sysctl machdep.cpu.brand_string

	// Memory. (both)
	MemoryBytes int64 `loglevel:"dangerous"` // 物理总内存字节数；mac: sysctl hw.memsize

	// OS version / build. (both)
	OSVersion     string `loglevel:"dangerous"` // 用户可读版本号；mac: "14.x"；win: "11 x"
	OSBuild       string `loglevel:"dangerous"` // 细粒度 build 号；mac: "23xxx"；win: "22xxx.xxx"
	KernelVersion string `loglevel:"dangerous"` // 内核版本；mac: uname -r；win: NT 版本号

	// Linux-only hints. (linux; mac/win 留空)
	Distro      string `loglevel:"dangerous"` // Linux 发行版标识
	SessionType string `loglevel:"dangerous"` // XDG_SESSION_TYPE：wayland / x11 / tty

	// Virtualization / emulation hint. (all; TASK-105 / Phase 2 实现)
	Virtualization string `loglevel:"dangerous"` // 可能值："rosetta" / "wsl" / "docker" / "vm" / ""
}

// ProcessBlock 是 Schema v1 中 process 分组的 Go 映射：
// 本进程的运行身份快照，字段均可从 kernel 直接读取。
//
// 采集失败的字符串字段沿用 "<err: ...>" 字串风格（见 dumpProcessBlock）；
// 整数字段采集失败或平台不适用时留哨位 -1（0 在 uid/gid/pid 等字段上
// 是合法值，不能当哨位）。错误信息保留到 stderr 诊断，便于样本分析
// 时定位原因。
//
// 平台策略：采用 macOS ∪ Windows 并集 schema；字段只在一边有语义时，
// 另一边留零值或哨位（数值 -1、字符串空、bool false、slice nil、
// Time 零值）。Linux 不作为一等目标——合辙则顺带可用，不合辙留空。
type ProcessBlock struct {
	// Invocation shape. (both)
	Argv []string `loglevel:"debug"` // 完整命令行参数（含 argv[0]）；nil 表示不可用
	CWD  string   `loglevel:"debug"` // current working directory；"<err: ...>" 表示 os.Getwd() 失败

	// 跨平台的可读 user / group 名. (both)
	Username  string `loglevel:"dangerous"` // 登录名；mac: 例如 "alice"；win: "DOMAIN\\user" 或 "user"
	Groupname string `loglevel:"dangerous"` // 主组名；mac: 例如 "staff"；win: 若可解析为可读名，否则为空

	// Kernel process IDs.
	PID  int `loglevel:"dangerous"` // (both) Process ID：本进程在 kernel 里的唯一编号
	PPID int `loglevel:"dangerous"` // (both) Parent PID：fork/exec（mac）/ CreateProcess（win）本进程的父进程
	PGID int `loglevel:"dangerous"` // (mac)  Process Group ID：进程组号，用于信号投递与前后台作业控制；win 留 -1
	SID  int `loglevel:"dangerous"` // (mac)  Unix Session ID：进程组的集合，绑定控制终端；win 留 -1（Windows 登录会话见 WindowsSessionID）

	// Unix numeric user / group identity. (mac; win 留 -1)
	UID  int `loglevel:"dangerous"` // real user ID：启动本进程的用户
	EUID int `loglevel:"dangerous"` // effective user ID：权限判定实际使用的 UID（setuid 后可与 UID 不同）
	GID  int `loglevel:"dangerous"` // real group ID：用户主组
	EGID int `loglevel:"dangerous"` // effective group ID：权限判定实际使用的 GID

	// Windows-specific identity. (win; mac 留空)
	WindowsUserSID        string `loglevel:"dangerous"` // 用户 SID，如 "S-1-5-21-..."
	WindowsGroupSID       string `loglevel:"dangerous"` // 主组 SID
	WindowsSessionID      int    `loglevel:"dangerous"` // 登录会话号；0=services session；1+ 为交互登录；mac 留 -1
	WindowsDomain         string `loglevel:"dangerous"` // 登录域；即 "DOMAIN\\user" 中 DOMAIN 部分
	WindowsIntegrityLevel string `loglevel:"dangerous"` // 进程完整性等级：Low / Medium / High / System
	WindowsElevated       bool   `loglevel:"dangerous"` // 是否以管理员运行；对应 token 的 Elevation 标志

	// Time. (both；两边都需平台代码实现，stdlib 不直接给)
	StartTime time.Time `loglevel:"dangerous"` // 本进程启动时刻（wall-clock）
}

// ChainBlockLayers 是 Schema v1 中 chain 分组的 Go 映射：
// 从直接父进程（os.Getppid()）沿 ppid 向上回溯，采集每一层祖先的
// 身份快照，直至命中语义停止条件。
//
// 遍历终止条件（写入 Termination）：
//   - "init":         当前 pid ≤ 1（爬到 launchd / init），自然终止
//   - "cycle":        已访问 pid 再次出现（数据异常），立即停
//   - "query_failed": ps 查询失败，记录原因、停
//   - "orphan":       ppid = 0（父死了但未挂到 init），极罕见
//   - "depth_cap":    累计访问层数达到 maxChainDepth = 64（死循环
//     保护，**非业务约束**；正常链通常 3~10 层）
//
// "不丢信息" 原则：每层除结构化字段外同时保留 RawPsLine（ps 输出
// 整行原文），即便解析出错也能回溯原始观测。
//
// **本 block 只记录 ps 直接返回的信息**，不做角色分类 / 指纹
// 匹配等二次判定——那些逻辑归特征指纹章节。
type ChainBlockLayers struct {
	Depth       int          `loglevel:"dangerous"` // len(Layers) 冗余存，便于快速筛选 / 诊断
	Termination string       `loglevel:"dangerous"` // 见上方枚举：init / cycle / query_failed / orphan / depth_cap
	Layers      []ChainBlock // 从 index 0（直接父进程）开始向上；不打印，由 for loop 单独打印每个 layer
}

// ChainBlock 祖先链单层快照：仅保存 ps 直接返回的原始字段。
//
// 原始字段来自 `ps -o pid=,ppid=,pgid=,sess=,uid=,tty=,user=,lstart=,comm=,command=`；
// 整数字段解析失败或平台不适用留哨位 -1（0 在 uid/pid/sid 等字段
// 上是合法值，不能当哨位）。字符串字段解析失败留空串。
//
// RawPsLine 保留 ps 整行输出原样，防止解析出错时丢信息；消费者
// 对结构化字段存疑时可回退看原文。
//
// 角色分类 / 规则匹配 等二次判定不在本结构里做——留给特征指纹章
// 节处理（见 .context/current-task.md 第八节）。
type ChainBlock struct {
	// Invocation shape.
	Comm    string `loglevel:"debug"` // 可执行 basename（ps -o comm=）；登录 shell 会带 "-" 前缀
	Command string `loglevel:"debug"` // 完整 argv 字串（ps -o command=）

	// Kernel process IDs. (mac)
	PID  int `loglevel:"dangerous"` // 本层 pid
	PPID int `loglevel:"dangerous"` // 本层父 pid（决定下一跳）
	PGID int `loglevel:"dangerous"` // (mac) 进程组号；win 留 -1
	SID  int `loglevel:"dangerous"` // (mac) Unix Session ID；win 留 -1

	// Identity. (mac)
	UID      int    `loglevel:"dangerous"` // 本层实际运行用户 uid；win 留 -1
	Username string `loglevel:"dangerous"` // ps -o user= 输出

	// TTY.
	TTY string `loglevel:"dangerous"` // 控制终端，如 "ttys003"；"??" 表示无 tty

	// Time (ps 原样字符串；time.Time 解析留给消费者).
	Lstart string `loglevel:"dangerous"` // 如 "Sun Apr 19 12:00:00 2026"

	// Raw fallback: 防解析失败丢信息.
	RawPsLine string `loglevel:"dangerous"` // ps 整行输出原样
}

// EnvBlockLayers 环境变量快照：按 family 分组 + 兜底 Uncategorized + 全量 Raw。
//
// 设计原则：
//   - **扁平**：顶层 14 个字段（12 family + Uncategorized + Raw），不嵌套
//   - **单归属**：每个 env 最多落一个 family（按 envFamilyOrder 优先级匹配，
//     更窄的 family 在前）；归不到进 Uncategorized
//   - **不丢信息**：Raw 保留 os.Environ() 原样（顺序 + 可能的重复），
//     即便家族规则错、将来规则调整，都能从 Raw 回溯
//
// 每个 family 用 map[string]string 而非固定 struct 字段，因为 IDE / Terminal
// 这类家族包含**前缀匹配**（VSCODE_* / WINDSURF_*），数量不定；用 map
// 保证 family 内全量命中不丢。
//
// 归类目前按 Unix 大小写敏感匹配；Windows env key 大小写不敏感，
// Windows 轮次会补 case-fold 逻辑（见 TASK-105 后续）。
//
// **本 block 只记录 os.Environ() 直接返回的信息**，不做语义判定 / 值解析
// （例如不解析 PATH 为 []string、不把 CI 值转 bool） —— 消费者自行处理。
type EnvBlockLayers struct {
	// Platform 组（专用 struct，用 printBlockByReflect 打印）
	Platform EnvPlatformBlock // PATH / HOME / USERPROFILE / TMPDIR / TEMP / TMP / APPDATA / LOCALAPPDATA

	// Shell 组（专用 struct，用 printBlockByReflect 打印）
	Shell EnvShellBlock // PSModulePath

	// KiloCode 组（专用 struct，用 printBlockByReflect 打印）
	KiloCode EnvKiloCodeBlock

	// Flutter 组（专用 struct，用 printBlockByReflect 打印）
	Flutter EnvFlutterBlock

	// 11 个 map family，按特异性从窄到宽排列（和匹配优先级一致）
	ExtensionHost map[string]string // VSCODE_CRASH_REPORTER_PROCESS_TYPE / VSCODE_PID / VSCODE_IPC_HOOK / VSCODE_ESM_ENTRYPOINT / VSCODE_NLS_CONFIG
	IDE           map[string]string // 前缀：VSCODE_ / WINDSURF_ / CURSOR_ / CASCADE_ / JETBRAINS_ / IDEA_ / PYCHARM_
	Session       map[string]string // SSH_CONNECTION / SSH_CLIENT / SSH_TTY / TMUX / STY
	CI            map[string]string // CI / GITHUB_ACTIONS / GITLAB_CI / CIRCLECI / BUILDKITE / JENKINS_URL / TRAVIS
	Remote        map[string]string // KUBERNETES_SERVICE_HOST / IN_NIX_SHELL
	Terminal      map[string]string // TERM / TERM_PROGRAM / TERM_PROGRAM_VERSION / COLORTERM / LC_TERMINAL / KITTY_WINDOW_ID / WEZTERM_EXECUTABLE / ALACRITTY_SOCKET / TERMINAL_EMULATOR
	ShellVersion  map[string]string // BASH_VERSION / ZSH_VERSION / FISH_VERSION
	Editor        map[string]string // EDITOR / VISUAL / PAGER
	Locale        map[string]string // LANG / LC_* (LC_TERMINAL 除外，已归 Terminal)
	User          map[string]string // Unix: USER / LOGNAME；Windows: USERDOMAIN
	ShellPref     map[string]string // SHELL

	// 兜底两道防线
	Uncategorized map[string]string // 未命中任何 family 的 env
	Raw           []string          // os.Environ() 全量原样，"KEY=VALUE" 格式，保留顺序和重复
}

// EnvPlatformBlock Platform 组的专用输出结构，用 printBlockByReflect 打印。
type EnvPlatformBlock struct {
	Path         string `loglevel:"debug"`
	HOME         string `loglevel:"debug"` // Mac
	USERPROFILE  string `loglevel:"debug"` // Windows
	TMPDIR       string `loglevel:"debug"` // Mac
	TEMP         string `loglevel:"debug"` // Windows
	TMP          string `loglevel:"debug"` // Windows
	APPDATA      string `loglevel:"debug"` // Windows
	LOCALAPPDATA string `loglevel:"debug"` // Windows
}

// EnvShellBlock Shell 组的专用输出结构，用 printBlockByReflect 打印。
type EnvShellBlock struct {
	SHELL            string `loglevel:"debug"`
	MSYS_NO_PATHCONV string `loglevel:"debug"`
	PSModulePath     string `loglevel:"debug"` // Windows PowerShell
}

// EnvKiloCodeBlock KiloCode 组的专用输出结构，用 printBlockByReflect 打印。
type EnvKiloCodeBlock struct {
	KILO_PLATFORM             string `loglevel:"debug"`
	KILO_APP_NAME             string `loglevel:"debug"`
	KILO_APP_VERSION          string `loglevel:"debug"`
	KILO_DISABLE_AUTOUPDATE   string `loglevel:"debug"`
	KILO_DISABLE_CLAUDE_CODE  string `loglevel:"debug"`
	KILO_ENABLE_QUESTION_TOOL string `loglevel:"debug"`

	KILO_EDITOR_NAME     string `loglevel:"debug"`
	KILOCODE_EDITOR_NAME string `loglevel:"debug"`
	KILOCODE_VERSION     string `loglevel:"debug"`
	KILOCODE_FEATURE     string `loglevel:"debug"`
	KILO_CLIENT          string `loglevel:"debug"`
	KILO_VSCODE_VERSION  string `loglevel:"debug"`
	KILO_TELEMETRY_LEVEL string `loglevel:"debug"`

	KILO_MACHINE_ID      string `loglevel:"dangerous"`
	KILO_PID             string `loglevel:"dangerous"`
	KILO_SERVER_PASSWORD string `loglevel:"dangerous"`
}

type EnvFlutterBlock struct {
	PUB_HOSTED_URL           string `loglevel:"debug"`
	FLUTTER_STORAGE_BASE_URL string `loglevel:"debug"`
}
