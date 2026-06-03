//go:build !darwin && !windows

package agent

func collectPlatformProcess(b *ProcessBlock) {}
