package gitea

import (
	"fmt"
	"os"

	"code.gitea.io/sdk/gitea"
)

// Adapter Gitea API 适配器
type Adapter struct {
	client *gitea.Client
}

// NewAdapter 创建新的 Gitea 适配器
func NewAdapter() (*Adapter, error) {
	url := os.Getenv("BACKSTAGE_GITEA_URL")
	token := os.Getenv("BACKSTAGE_GITEA_TOKEN")

	if url == "" {
		return nil, fmt.Errorf("BACKSTAGE_GITEA_URL is required")
	}

	client, err := gitea.NewClient(url, gitea.SetToken(token))
	if err != nil {
		return nil, fmt.Errorf("Gitea.NewClient failed: %w", err)
	}

	return &Adapter{client: client}, nil
}

// ListIssueOfRepo 列出仓库 Issue 列表
func (a *Adapter) ListIssueOfRepo(owner, repoName string) ([]*gitea.Issue, error) {
	// 当前项目体量较小，直接使用较大的 PageSize 一次性获取
	// 后续若需跨页分页，可改为 Page 遍历
	opts := gitea.ListIssueOption{}
	opts.PageSize = 100
	issues, _, err := a.client.ListRepoIssues(owner, repoName, opts)
	return issues, err
}

// CreateIssue 创建 Issue
// milestoneId: Milestone 的数字 ID（可选，传入空字符串表示不关联）
func (a *Adapter) CreateIssue(owner, repo, title string, context string, milestoneId string) (*gitea.Issue, error) {
	opts := gitea.CreateIssueOption{
		Title: title,
		Body:  context,
	}
	if milestoneId != "" {
		var id int64
		fmt.Sscanf(milestoneId, "%d", &id)
		if id > 0 {
			opts.Milestone = id
		}
	}
	issue, _, err := a.client.CreateIssue(owner, repo, opts)
	return issue, err
}
