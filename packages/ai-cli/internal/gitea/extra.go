package gitea

import (
	"strings"

	"code.gitea.io/sdk/gitea"
)

// extra.go 提供 Gitea SDK 的功能扩展
// 通过组合 adapter 的原子接口，构建项目所需的复合功能
// 保持纯工具层属性，不承载业务逻辑

// SearchIssueByPrefix 根据 Issue Title 前缀搜索指定仓库的所有 Issues
func (a *Adapter) SearchIssueByPrefix(owner, repoName, prefix string) ([]*gitea.Issue, error) {
	issues, err := a.ListIssueOfRepo(owner, repoName)

	if err != nil {
		return nil, err
	}

	var matchedIssues []*gitea.Issue
	for _, issue := range issues {
		if strings.HasPrefix(issue.Title, prefix) {
			matchedIssues = append(matchedIssues, issue)
		}
	}

	return matchedIssues, nil
}
