package blame

import (
	"fmt"

	internalGitea "com.lisitede.ai/internal/gitea"
	"com.lisitede.ai/internal/plan"
)

func CreateBlame(appId, planId, title, context string) (map[string]interface{}, error) {
	blameId, err := plan.GenNextBlameId(appId, planId)
	if err != nil {
		return nil, err
	}

	blameTitle := fmt.Sprintf("%s: %s", blameId, title)

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	issue, err := adapter.CreateIssue(appId, planId, blameTitle, context, "")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":      blameId,
		"title":   issue.Title,
		"context": context,
	}, nil
}
