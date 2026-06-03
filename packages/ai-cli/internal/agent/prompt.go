package agent

import (
	"bytes"
	"strings"
	"text/template"
)

// hasshinPromptTemplate 是 Hasshin 向 agent 交付的环境上下文提示模板。
//
// 用 text/template 而非字符串拼接，便于后续扩展更多字段（例如 shell / terminal
// 指纹、CI 标志）而不污染调用点；数据入参 hasshinPromptData 只保留模板需要的
// 最小视图，解耦 RuntimeBlock 等底层结构。
//
// 输出面向 LLM，英文、短句、祈使语气；避免情绪化 / 冗余修辞。
const hasshinPromptTemplate = `⇒⇒⇒⇒ Agent Hasshin ⇒⇒⇒⇒
The agent is running on a {{.OS}}/{{.Arch}} host.
Its invocation chain (nearest parent first) is: {{.Chain}}.
Take this context into account when choosing which tools to use.
`

// hasshinPromptData 是 hasshinPromptTemplate 的数据视图。
type hasshinPromptData struct {
	OS    string // RuntimeBlock.OS
	Arch  string // RuntimeBlock.Arch
	Chain string // 归一化后的链路，形如 "zsh <- Windsurf"；空链时为 "<unknown>"
}

// hasshinPromptTmpl 在包加载时一次性编译，避免每次调用重解析。
// template.Must 会在模板语法错误时 panic，保证部署前暴露问题。
var hasshinPromptTmpl = template.Must(template.New("hasshin").Parse(hasshinPromptTemplate))

// renderHasshinPrompt 把 runtime 信息 + 归一链路拼成最终的 agent 提示文本。
// chain 顺序：index 0 是直接父进程，向上递增；渲染时用 " <- " 连接，语义为
// "caller <- caller-of-caller"，贴合 shell 链条的阅读习惯。
func renderHasshinPrompt(rt RuntimeBlock, chain []string) string {
	data := hasshinPromptData{
		OS:    rt.OS,
		Arch:  rt.Arch,
		Chain: formatChain(chain),
	}
	var buf bytes.Buffer
	// Execute 仅在模板字段引用不存在 / IO 失败时报错；此处 buf 不会失败，
	// 字段已与模板对齐，实践中不会触发错误路径，但仍保留 fallback 防御。
	if err := hasshinPromptTmpl.Execute(&buf, data); err != nil {
		return "<prompt render failed: " + err.Error() + ">"
	}
	return buf.String()
}

// formatChain 把归一链路 slice 转成模板可用字符串；空链回退到 "<unknown>"，
// 避免模板渲染出 "... is: ." 这种语义悬空的句子。
func formatChain(chain []string) string {
	if len(chain) == 0 {
		return "<unknown>"
	}
	return strings.Join(chain, " <- ")
}
