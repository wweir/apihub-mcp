// Package handlers 包含所有 API 接口的处理函数
package handlers

import (
	"context"
	"fmt"

	"git.tophant.com/scc/cyber-guard/tools/tool-call/apifox/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sower-proxy/deferlog/v2"
)

// SprintBranch 项目迭代分支
// 对应API: GET https://api.apifox.com/api/v1/projects/{projectId}/sprint-branches
type SprintBranch struct {
	Id                                 int64   `json:"id"`                                 // 分支ID
	CreatedAt                          string  `json:"createdAt"`                          // 创建时间
	CreatedById                        int64   `json:"createdById"`                        // 创建者ID
	IsArchived                         bool    `json:"isArchived"`                         // 是否归档
	IsMain                             bool    `json:"isMain"`                             // 是否为主分支
	IsProtected                        bool    `json:"isProtected"`                        // 是否受保护
	Name                               string  `json:"name"`                               // 分支名称
	ProjectId                          int64   `json:"projectId"`                          // 项目ID
	Type                               string  `json:"type"`                               // 分支类型
	ForkFromBranchId                   int64   `json:"forkFromBranchId"`                   // 从哪个分支 fork
	AdminUserIds                       []int64 `json:"adminUserIds"`                       // 管理员用户ID列表
	EnableAdminUpdateToProtectedBranch bool    `json:"enableAdminUpdateToProtectedBranch"` // 是否允许管理员更新到受保护分支
	ProjectBranchState                 struct {
		ApiCount               int64 `json:"apiCount"`               // API数量
		DataSchemaCount        int64 `json:"dataSchemaCount"`        // 数据模式数量
		ResponseComponentCount int64 `json:"responseComponentCount"` // 响应组件数量
		DocCount               int64 `json:"docCount"`               // 文档数量
		ApiTestCaseCount       int64 `json:"apiTestCaseCount"`       // API测试用例数量
	} `json:"projectBranchState"` // 项目分支状态
}

// SprintBranchesResponse 项目迭代分支列表响应
// 对应API: GET https://api.apifox.com/api/v1/projects/{projectId}/sprint-branches
type SprintBranchesResponse struct {
	Data []SprintBranch `json:"data"` // 分支列表
	Meta struct {
		BranchQuantityLimitation int64 `json:"branchQuantityLimitation"` // 分支数量限制
		BranchQuantityExceeded   bool  `json:"branchQuantityExceeded"`   // 是否超出分支数量限制
		BranchSettings           struct {
			MainBranch struct {
				EnableRefFallbackToMain  bool `json:"enableRefFallbackToMain"`  // 是否启用引用回退到主分支
				EnableResourceFork       bool `json:"enableResourceFork"`       // 是否启用资源 fork
				EnableResourceMerge      bool `json:"enableResourceMerge"`      // 是否启用资源合并
				EnableEnvironmentSetting bool `json:"enableEnvironmentSetting"` // 是否启用环境设置
				EnableResourceShareDoc   bool `json:"enableResourceShareDoc"`   // 是否启用资源共享文档
				EnableImportData         bool `json:"enableImportData"`         // 是否启用导入数据
				EnableExportData         bool `json:"enableExportData"`         // 是否启用导出数据
				EnableRequestHistory     bool `json:"enableRequestHistory"`     // 是否启用请求历史
				EnableApiTest            bool `json:"enableApiTest"`            // 是否启用API测试
			} `json:"mainBranch"` // 主分支设置
			SprintBranch struct {
				EnableRefFallbackToMain  bool `json:"enableRefFallbackToMain"`  // 是否启用引用回退到主分支
				EnableResourceFork       bool `json:"enableResourceFork"`       // 是否启用资源 fork
				EnableResourceMerge      bool `json:"enableResourceMerge"`      // 是否启用资源合并
				EnableEnvironmentSetting bool `json:"enableEnvironmentSetting"` // 是否启用环境设置
				EnableResourceShareDoc   bool `json:"enableResourceShareDoc"`   // 是否启用资源共享文档
				EnableImportData         bool `json:"enableImportData"`         // 是否启用导入数据
				EnableExportData         bool `json:"enableExportData"`         // 是否启用导出数据
				EnableRequestHistory     bool `json:"enableRequestHistory"`     // 是否启用请求历史
				EnableApiTest            bool `json:"enableApiTest"`            // 是否启用API测试
			} `json:"sprintBranch"` // 迭代分支设置
		} `json:"branchSettings"` // 分支设置
	} `json:"meta"` // 元数据
}

// HandleGetSprintBranches 处理获取项目迭代分支列表的请求
func HandleGetSprintBranches(ctx context.Context, request mcp.CallToolRequest) (_ *mcp.CallToolResult, err error) {
	defer func() { deferlog.DebugErrorContext(ctx, err, "HandleGetSprintBranches completed") }()

	projectId := int(mcp.ParseInt(request, "projectId", 0))
	locale := mcp.ParseString(request, "locale", "zh-CN")

	url := fmt.Sprintf("https://api.apifox.com/api/v1/projects/%d/sprint-branches?locale=%s",
		projectId, locale)

	type Response struct {
		Success bool                   `json:"success"`
		Data    []SprintBranch         `json:"data"`
		Meta    SprintBranchesResponse `json:"meta"`
	}

	return handleApiRequest(ctx, request, func() (*ApiResponse, error) {
		resp, err := client.GetJSON[Response](url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get sprint branches (projectId=%d): %v", projectId, err)
		}
		return &ApiResponse{
			Success: resp.Success,
			Data:    resp.Data,
			Meta:    resp.Meta,
		}, nil
	})
}
