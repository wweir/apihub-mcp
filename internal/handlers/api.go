// Package handlers 包含所有 API 接口的处理函数
package handlers

import (
	"context"
	"fmt"

	"git.tophant.com/scc/cyber-guard/tools/tool-call/apifox/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sower-proxy/deferlog/v2"
)

// HandleGetAPIDetails 处理获取API详情的请求
func HandleGetAPIDetails(ctx context.Context, request mcp.CallToolRequest) (_ *mcp.CallToolResult, err error) {
	defer func() { deferlog.DebugErrorContext(ctx, err, "HandleGetAPIDetails completed") }()

	projectId := int(mcp.ParseInt(request, "projectId", 0))
	branchId := int(mcp.ParseInt(request, "branchId", 0))
	locale := mcp.ParseString(request, "locale", "zh-CN")

	url := fmt.Sprintf("https://api.apifox.com/api/v1/api-details?locale=%s", locale)

	headers := map[string]string{
		"x-branch-id":  fmt.Sprintf("%d", branchId),
		"x-project-id": fmt.Sprintf("%d", projectId),
	}

	return handleApiRequest(ctx, request, func() (*ApiResponse, error) {
		type Response struct {
			Success bool        `json:"success"`
			Data    interface{} `json:"data"`
		}
		resp, err := client.GetJSON[Response](url, headers)
		if err != nil {
			return nil, fmt.Errorf("failed to get API details (projectId=%d, branchId=%d): %v", projectId, branchId, err)
		}
		return &ApiResponse{
			Success: resp.Success,
			Data:    resp.Data,
		}, nil
	})
}
