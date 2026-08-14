package vikunja

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Call 调用 backstage-vikunja CLI，同步阻塞，返回 stdout
func Call(args ...string) (string, error) {
	output, err := exec.Command("backstage-vikunja", args...).CombinedOutput()
	if err != nil {
		if details := strings.TrimSpace(string(output)); details != "" {
			return "", fmt.Errorf("%s: %w", details, err)
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// UpdateTaskStatus 将任务移动到指定的看板状态。
func UpdateTaskStatus(workspaceId, taskId, status string) error {
	_, err := Call(
		"scrum",
		"PATCH",
		"/workspaces/"+workspaceId+"/tasks/"+taskId,
		"--status",
		status,
	)
	return err
}

// ListWorkspaces 列出所有工作区
func ListWorkspaces() ([]Workspace, error) {
	raw, err := Call("scrum", "GET", "/workspaces")
	if err != nil {
		return nil, err
	}

	var workspaces []Workspace
	if err := json.Unmarshal([]byte(raw), &workspaces); err != nil {
		return nil, err
	}
	return workspaces, nil
}

// ListTasksOfWorkspace 列出指定工作区下的任务
func ListTasksOfWorkspace(workspaceId string, filter string) ([]VikunjaTask, error) {
	raw, err := Call("scrum", "GET", "/workspaces/"+workspaceId+"/tasks", "--filter", filter)
	if err != nil {
		return nil, err
	}

	if raw == "" {
		return []VikunjaTask{}, nil
	}

	var tasks []VikunjaTask
	if err := json.Unmarshal([]byte(raw), &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
