package plan

import (
	"fmt"
	"sort"
	"strconv"

	internalGitea "com.lisitede.ai/internal/gitea"
)

func GenNextBlameId(appName, planId string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", err
	}

	issues, err := adapter.SearchIssueByPrefix(appName, planId, "BLAME-")
	if err != nil {
		return "", err
	}

	var ids []string
	for _, issue := range issues {
		_, _, id, err := ExtractBlameId(issue.Title)
		if err != nil {
			return "", err
		}
		ids = append(ids, id)
	}

	sort.Strings(ids)
	if len(ids) == 0 {
		ids = append(ids, "000")
	}

	nextId, err := genNextId(ids)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("BLAME-%s", nextId), nil
}

func GenNextUseId(appName, planId string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", err
	}

	issues, err := adapter.SearchIssueByPrefix(appName, planId, "USE-")
	if err != nil {
		return "", err
	}

	var ids []string
	for _, issue := range issues {
		_, _, id, err := ExtractUseId(issue.Title)
		if err != nil {
			return "", err
		}
		ids = append(ids, id)
	}

	sort.Strings(ids)
	if len(ids) == 0 {
		ids = append(ids, "000")
	}

	nextId, err := genNextId(ids)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("USE-%s", nextId), nil
}

// genNextId 私有方法
// 规则：连续递增（如 100 → 101，225 → 226）
func genNextId(existIds []string) (string, error) {
	if len(existIds) == 0 {
		return "100", nil
	}

	lastNum := existIds[len(existIds)-1]
	lastInt, err := strconv.Atoi(lastNum)
	if err != nil {
		return "", fmt.Errorf("invalid task id: %s", lastNum)
	}

	next := lastInt + 1
	if next >= 900 {
		return "", fmt.Errorf("no available task id: reached 900")
	}

	return fmt.Sprintf("%03d", next), nil
}
