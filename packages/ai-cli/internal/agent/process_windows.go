//go:build windows

package agent

import (
	"os/user"
	"strings"
	"syscall"
	"unsafe"
)

func collectPlatformProcess(b *ProcessBlock) {
	u, err := user.Current()
	if err != nil {
		return
	}

	b.WindowsUserSID = u.Uid
	b.WindowsGroupSID = u.Gid

	if strings.Contains(b.Username, "\\") {
		parts := strings.Split(b.Username, "\\")
		b.WindowsDomain = parts[0]
	}

	collectWindowsTokenInfo(b)
}

func collectWindowsTokenInfo(b *ProcessBlock) {
	var token syscall.Token
	currentProcess, _ := syscall.GetCurrentProcess()
	err := syscall.OpenProcessToken(currentProcess, syscall.TOKEN_QUERY, &token)
	if err != nil {
		return
	}
	defer token.Close()

	b.WindowsElevated = isTokenElevated(token)
	b.WindowsSessionID = getTokenSessionID(token)
	b.WindowsIntegrityLevel = getTokenIntegrityLevel(token)
}

func isTokenElevated(token syscall.Token) bool {
	var elevated uint32
	var returned uint32
	err := syscall.GetTokenInformation(token, syscall.TokenElevation, (*byte)(unsafe.Pointer(&elevated)), 4, &returned)
	if err != nil {
		return false
	}
	return elevated != 0
}

func getTokenSessionID(token syscall.Token) int {
	var sessionID uint32
	var returned uint32
	err := syscall.GetTokenInformation(token, syscall.TokenSessionId, (*byte)(unsafe.Pointer(&sessionID)), 4, &returned)
	if err != nil {
		return -1
	}
	return int(sessionID)
}

func getTokenIntegrityLevel(token syscall.Token) string {
	var returned uint32
	syscall.GetTokenInformation(token, syscall.TokenIntegrityLevel, nil, 0, &returned)
	if returned == 0 {
		return ""
	}

	buf := make([]byte, returned)
	err := syscall.GetTokenInformation(token, syscall.TokenIntegrityLevel, &buf[0], returned, &returned)
	if err != nil {
		return ""
	}

	type TOKEN_MANDATORY_LABEL struct {
		Label syscall.SIDAndAttributes
	}

	tml := (*TOKEN_MANDATORY_LABEL)(unsafe.Pointer(&buf[0]))
	rid := getSIDRID(tml.Label.Sid)

	switch rid {
	case 0x0000:
		return "Untrusted"
	case 0x1000:
		return "Low"
	case 0x2000:
		return "Medium"
	case 0x3000:
		return "High"
	case 0x4000:
		return "System"
	default:
		return "Unknown"
	}
}

func getSIDRID(sid *syscall.SID) uint32 {
	type SIDInternal struct {
		Revision            byte
		SubAuthorityCount   byte
		IdentifierAuthority [6]byte
		SubAuthority        [1]uint32
	}

	internal := (*SIDInternal)(unsafe.Pointer(sid))
	count := internal.SubAuthorityCount
	if count == 0 {
		return 0
	}

	subAuthPtr := unsafe.Pointer(&internal.SubAuthority[0])
	subAuthSlice := (*[256]uint32)(subAuthPtr)[:count:count]
	return subAuthSlice[count-1]
}
