//go:build darwin

package agent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func queryChainLayerOS(pid int) (ChainBlock, error) {
	layer := ChainBlock{
		PID:  pid,
		PPID: -1,
		PGID: -1,
		SID:  -1,
		UID:  -1,
	}
	pidStr := strconv.Itoa(pid)

	// 第 1 次：仅取 token 可分的定长字段（lstart 恒为 5 token）。
	// 不与 comm/command 同行采集，避免 comm 内嵌空格（如 "Application Support"）破坏切分。
	out, err := exec.Command("ps",
		"-o", "pid=,ppid=,pgid=,sess=,uid=,tty=,user=,lstart=",
		"-p", pidStr).CombinedOutput()
	if err != nil {
		layer.RawPsLine = fmt.Sprintf("<ps err: %v; output=%q>", err, strings.TrimSpace(string(out)))
		return layer, err
	}
	raw := strings.TrimRight(string(out), "\r\n")
	layer.RawPsLine = raw

	fields := strings.Fields(raw)
	if len(fields) >= 12 {
		if n, e := strconv.Atoi(fields[0]); e == nil {
			layer.PID = n
		}
		if n, e := strconv.Atoi(fields[1]); e == nil {
			layer.PPID = n
		}
		if n, e := strconv.Atoi(fields[2]); e == nil {
			layer.PGID = n
		}
		if n, e := strconv.Atoi(fields[3]); e == nil {
			layer.SID = n
		}
		if n, e := strconv.Atoi(fields[4]); e == nil {
			layer.UID = n
		}
		layer.TTY = fields[5]
		layer.Username = fields[6]
		layer.Lstart = strings.Join(fields[7:12], " ")
	}

	// comm / command 改走 kernel sysctl KERN_PROCARGS2：
	// 原因：ps -o comm= 受 MAXCOMLEN=16 字节截断（内核层 p_comm 限制），
	// 且 ps 文本输出的 comm 与 command 之间无分隔，在含空格路径
	// （如 "Application Support"）上无法可靠切分。
	// KERN_PROCARGS2 直接返回 kernel 为 exec 保存的原始缓冲区：
	//   [argc:int32] [exec_path\0] [align \0...] [argv[0]\0 argv[1]\0 ...] [envp...]
	// 完整路径、完整 argv、不截断、空格安全，是 ps(1)/Activity Monitor 使用的同一条路径。
	exe, argv, perr := procArgs2(pid)
	if perr == nil {
		layer.Comm = exe
		layer.Command = strings.Join(argv, " ")
	} else {
		// Fallback：KERN_PROCARGS2 对僵尸进程 / 权限不足时可能失败，
		// 此时退回到 ps -o command=（至少保住可读命令行），comm 留空。
		if out, err := exec.Command("ps", "-o", "command=", "-p", pidStr).CombinedOutput(); err == nil {
			layer.Command = strings.TrimRight(string(out), "\r\n")
		}
		layer.RawPsLine += fmt.Sprintf(" | <procargs2 err: %v>", perr)
	}
	return layer, nil
}

// procArgs2 通过 sysctl KERN_PROCARGS2 解析目标进程的完整 exec path 与 argv。
// 仅在 darwin 可用；返回值：exe = 完整可执行路径（未截断），argv = 原始命令行
// 参数数组（含 argv[0]，空格安全）。缓冲区格式参见 Apple xnu sysctl.c /
// ps(1) 源码对 "KERN_PROCARGS2" 的解析约定。
func procArgs2(pid int) (exe string, argv []string, err error) {
	buf, err := unix.SysctlRaw("kern.procargs2", pid)
	if err != nil {
		return "", nil, err
	}
	if len(buf) < 4 {
		return "", nil, fmt.Errorf("procargs2 buffer too short: %d", len(buf))
	}
	argc := int(binary.LittleEndian.Uint32(buf[:4]))
	rest := buf[4:]

	// exec_path：到第一个 NUL 为止。
	end := bytes.IndexByte(rest, 0)
	if end < 0 {
		return "", nil, fmt.Errorf("procargs2: no NUL after argc")
	}
	exe = string(rest[:end])
	rest = rest[end:]

	// 跳过对齐填充的连续 NUL。
	for len(rest) > 0 && rest[0] == 0 {
		rest = rest[1:]
	}

	// 依次读取 argc 个以 NUL 分隔的 argv 项。
	argv = make([]string, 0, argc)
	for i := 0; i < argc && len(rest) > 0; i++ {
		end := bytes.IndexByte(rest, 0)
		if end < 0 {
			argv = append(argv, string(rest))
			break
		}
		argv = append(argv, string(rest[:end]))
		rest = rest[end+1:]
	}
	return exe, argv, nil
}
