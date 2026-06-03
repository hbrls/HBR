//go:build windows

package agent

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

func queryChainLayerOS(pid int) (ChainBlock, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return ChainBlock{
			PID:       pid,
			PPID:      -1,
			PGID:      -1,
			SID:       -1,
			UID:       -1,
			RawPsLine: fmt.Sprintf("<gopsutil err: %v; pid=%d>", err, pid),
		}, err
	}

	ppid, err := p.Ppid()
	if err != nil {
		ppid = -1
	}

	// Exe() 走 QueryFullProcessImageNameW，返回完整绝对路径（含空格，不截断），
	// 等价于 darwin 上 KERN_PROCARGS2 的 exec_path 语义。
	// 权限受限（System / csrss.exe 等受保护进程）时会失败，此时回退到 Name()
	// （仅 basename），供下游启发式使用。
	name, err := p.Exe()
	if err != nil || name == "" {
		if n, e := p.Name(); e == nil {
			name = n
		} else {
			name = ""
		}
	}

	cmdline, err := p.Cmdline()
	if err != nil {
		cmdline = ""
	}

	username, err := p.Username()
	if err != nil {
		username = ""
	}

	createTime, err := p.CreateTime()
	if err != nil {
		createTime = 0
	}

	raw := fmt.Sprintf("pid=%d ppid=%d name=%q cmdline=%q username=%q create_time=%d",
		pid, ppid, name, cmdline, username, createTime)

	layer := ChainBlock{
		PID:       pid,
		PPID:      int(ppid),
		PGID:      -1,
		SID:       -1,
		UID:       -1,
		Comm:      name,
		Command:   cmdline,
		Username:  username,
		Lstart:    "",
		RawPsLine: raw,
	}

	return layer, nil
}
