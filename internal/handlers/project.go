// Package handlers 包含所有 API 接口的处理函数
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sower-proxy/deferlog/v2"
)

// ProjectInfo 项目信息
// 对应API: GET https://api.apifox.com/api/v1/projects/{projectId}
type ProjectInfo struct {
	Id          int64  `json:"id"`          // 项目ID
	Name        string `json:"name"`        // 项目名称
	Visibility  string `json:"visibility"`  // 可见性（public/private）
	Description string `json:"description"` // 项目描述
	Icon        string `json:"icon"`        // 项目图标
	MockRule    any    `json:"mockRule"`    // 模拟规则
	RoleType    int64  `json:"roleType"`    // 角色类型
	Type        string `json:"type"`        // 项目类型
	TeamId      int64  `json:"teamId"`      // 团队ID
}

// HandleGetProjectInfo 处理获取项目信息的请求
func HandleGetProjectInfo(ctx context.Context, request mcp.CallToolRequest) (_ *mcp.CallToolResult, err error) {
	defer func() { deferlog.DebugErrorContext(ctx, err, "HandleGetProjectInfo completed") }()

	projectId := int(mcp.ParseInt(request, "projectId", 0))
	locale := mcp.ParseString(request, "locale", "zh-CN")

	url := fmt.Sprintf("https://api.apifox.com/api/v1/projects/%d?locale=%s",
		projectId, locale)

	return handleApiRequest(ctx, request, func() (*ApiResponse, error) {
		resp, err := GetApiResponse[ProjectInfo](url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get project info (projectId=%d): %v", projectId, err)
		}
		return resp, nil
	})
}
