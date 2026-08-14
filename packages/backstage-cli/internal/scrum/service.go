package scrum

import (
	"cmp"
	"strconv"

	"com.lisitede.backstage/framework"
	"com.lisitede.backstage/internal/vikunja"
)

// PullTask 拉取所有工作区的任务，找出优先级最高的一个
func PullTask(filter, status string) (*Task, error) {
	workspaces, err := vikunja.ListWorkspaces()
	if err != nil {
		return nil, err
	}

	// workspaceId → workspace 映射
	wsMap := make(map[int64]vikunja.Workspace, len(workspaces))
	for _, ws := range workspaces {
		wsMap[ws.Id] = ws
	}

	var topTask *vikunja.VikunjaTask
	var topPriority int64
	var topWsId int64
	for _, ws := range workspaces {
		tasks, err := vikunja.ListTasksOfWorkspace(strconv.FormatInt(ws.Id, 10), filter)
		if err != nil {
			return nil, err
		}
		for i := range tasks {
			if topTask == nil || cmp.Or(
				cmp.Compare(tasks[i].Priority, topPriority),
				cmp.Compare(wsMap[tasks[i].WorkspaceId].Priority, wsMap[topWsId].Priority),
			) > 0 {
				topTask = &tasks[i]
				topPriority = tasks[i].Priority
				topWsId = tasks[i].WorkspaceId
			}
		}
	}

	if topTask == nil {
		return nil, nil
	}
	if err := vikunja.UpdateTaskStatus(
		strconv.FormatInt(topTask.WorkspaceId, 10),
		strconv.FormatInt(topTask.Id, 10),
		status,
	); err != nil {
		return nil, err
	}

	return &Task{
		Id:        framework.EncodeId(topTask.Id),
		Name:      topTask.Name,
		Context:   topTask.Context,
		Workspace: wsMap[topTask.WorkspaceId].Title,
	}, nil
}
