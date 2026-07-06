package blame

import (
	"fmt"

	internalGitea "com.lisitede.ai/internal/gitea"
	"com.lisitede.ai/internal/plan"
)

func CreateBlame(appId, planId, title, context string) (map[string]interface{}, error) {
	const aId = "backstage"
	const pId = "blames"

	blameId, err := plan.GenNextBlameId(aId, pId)
	if err != nil {
		return nil, err
	}

	blameTitle := fmt.Sprintf("%s: %s", blameId, title)

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	appLabel, err := adapter.GetLabelByName(aId, pId, "app/"+appId)
	if err != nil {
		appLabel, err = adapter.GetLabelByName(aId, pId, "app/unknown")
		if err != nil {
			return nil, err
		}
	}

	planLabel, err := adapter.GetLabelByName(aId, pId, "plan/"+planId)
	if err != nil {
		planLabel, err = adapter.GetLabelByName(aId, pId, "plan/unknown")
		if err != nil {
			return nil, err
		}
	}

	issue, err := adapter.CreateIssue(aId, pId, blameTitle, context, "")
	if err != nil {
		return nil, err
	}

	issueNo := fmt.Sprintf("%d", issue.Index)
	if _, err := adapter.AddLabelToIssue(aId, pId, issueNo, fmt.Sprintf("%d", appLabel.ID)); err != nil {
		return nil, err
	}
	if _, err := adapter.AddLabelToIssue(aId, pId, issueNo, fmt.Sprintf("%d", planLabel.ID)); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":      blameId,
		"title":   issue.Title,
		"context": context,
	}, nil
}
