//go:build !darwin && !windows

package agent

import (
	"fmt"
	"runtime"

	"com.lisitede.ai/framework"
)

func queryChainLayerOS(pid int) (ChainBlock, error) {
	msg := fmt.Sprintf("chain query on %s", runtime.GOOS)
	return ChainBlock{}, framework.NotImplementedException(msg)
}
