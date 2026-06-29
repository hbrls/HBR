package use

import (
	"fmt"

	internalGitea "com.lisitede.ai/internal/gitea"
	"com.lisitede.ai/internal/plan"
)

func CreateUse(appId, planId, title, context string) (map[string]interface{}, error) {
	useId, err := plan.GenNextUseId(appId, planId)
	if err != nil {
		return nil, err
	}

	useTitle := fmt.Sprintf("%s: %s", useId, title)

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	issue, err := adapter.CreateIssue(appId, planId, useTitle, context, "")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":      useId,
		"title":   issue.Title,
		"context": context,
	}, nil
}
