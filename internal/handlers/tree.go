// Package handlers 包含所有 API 接口的处理函数
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sower-proxy/deferlog/v2"
)

// ApiInfo API信息
// 对应API: GET https://api.apifox.com/api/v1/api-details
type ApiInfo struct {
	Id              int64    `json:"id"`              // API ID
	Name            string   `json:"name"`            // API名称
	ModuleId        int64    `json:"moduleId"`        // 模块ID
	Type            string   `json:"type"`            // API类型
	Method          string   `json:"method"`          // HTTP方法
	Path            string   `json:"path"`            // 路径
	FolderId        int64    `json:"folderId"`        // 文件夹ID
	Tags            []string `json:"tags"`            // 标签列表
	Status          int64    `json:"status"`          // 状态
	ResponsibleId   int64    `json:"responsibleId"`   // 负责人ID
	CustomApiFields any      `json:"customApiFields"` // 自定义API字段
	EditorId        int64    `json:"editorId"`        // 编辑者ID
}

// ApiDetailFolder API文档目录
// 对应API: GET https://api.apifox.com/api/v1/projects/{projectId}/api-tree-list
type ApiDetailFolder struct {
	Key      string            `json:"key"`           // 目录关键字
	Type     string            `json:"type"`          // 类型（apiDetailFolder或apiDetail）
	Name     string            `json:"name"`          // 名称
	ModuleId int64             `json:"moduleId"`      // 模块ID
	Children []ApiDetailFolder `json:"children"`      // 子目录
	Api      *ApiInfo          `json:"api,omitempty"` // API信息（如果是API）
	Folder   *struct {
		Id              int64  `json:"id"`              // 文件夹ID
		Name            string `json:"name"`            // 文件夹名称
		ModuleId        int64  `json:"moduleId"`        // 模块ID
		DocId           int64  `json:"docId"`           // 文档ID
		ParentId        int64  `json:"parentId"`        // 父文件夹ID
		ProjectBranchId int64  `json:"projectBranchId"` // 项目分支ID
		ShareSettings   any    `json:"shareSettings"`   // 共享设置
		EditorId        int64  `json:"editorId"`        // 编辑者ID
		Type            string `json:"type"`            // 文件夹类型
	} `json:"folder,omitempty"` // 文件夹信息（如果是文件夹）
}

// ApiTreeListResponse 项目文档目录树响应
// 对应API: GET https://api.apifox.com/api/v1/projects/{projectId}/api-tree-list
type ApiTreeListResponse struct {
	Data []ApiDetailFolder `json:"data"` // 目录树数据
}

// HandleGetAPITreeList 处理获取项目文档目录树的请求
func HandleGetAPITreeList(ctx context.Context, request mcp.CallToolRequest) (_ *mcp.CallToolResult, err error) {
	defer func() { deferlog.DebugErrorContext(ctx, err, "HandleGetAPITreeList completed") }()

	projectId := int(mcp.ParseInt(request, "projectId", 0))
	locale := mcp.ParseString(request, "locale", "zh-CN")

	url := fmt.Sprintf("https://api.apifox.com/api/v1/projects/%d/api-tree-list?locale=%s",
		projectId, locale)

	return handleApiRequest(ctx, request, func() (*ApiResponse, error) {
		resp, err := GetApiResponse[[]ApiDetailFolder](url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get API tree list (projectId=%d): %v", projectId, err)
		}
		return resp, nil
	})
}
