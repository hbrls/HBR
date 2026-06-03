//go:build darwin

package agent

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func collectPlatformProcess(b *ProcessBlock) {
	if pgid, err := syscall.Getpgid(os.Getpid()); err == nil {
		b.PGID = pgid
	} else {
		b.PGID = -1
	}
	if sid, err := syscall.Getsid(os.Getpid()); err == nil {
		b.SID = sid
	} else {
		b.SID = -1
	}
	b.GID = syscall.Getgid()
	b.EGID = syscall.Getegid()

	if g, err := user.LookupGroupId(strconv.Itoa(b.GID)); err == nil {
		b.Groupname = g.Name
	}
}
