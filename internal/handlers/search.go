// Package handlers 包含所有 API 接口的处理函数
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sower-proxy/deferlog/v2"
)

// PublicProject 公共项目信息
// 对应API: GET https://api.apifox.com/api/v1/public-projects
type PublicProject struct {
	Id          int64   `json:"id"`          // 项目ID
	Name        string  `json:"name"`        // 项目名称
	Description string  `json:"description"` // 项目描述
	Icon        string  `json:"icon"`        // 项目图标
	CategoryIds []int64 `json:"categoryIds"` // 分类ID列表
	Views       int64   `json:"views"`       // 查看次数
	Collections int64   `json:"collections"` // 收藏次数
	Grade       int64   `json:"grade"`       // 评分
	RoleType    int64   `json:"roleType"`    // 角色类型
	SysDomain   string  `json:"sysDomain"`   // 系统域名
	DomainName  string  `json:"domainName"`  // 域名名称
	DocsSiteId  int64   `json:"docsSiteId"`  // 文档站点ID
}

// SearchProjectsResponse 搜索公共项目响应
// 对应API: GET https://api.apifox.com/api/v1/public-projects
type SearchProjectsResponse struct {
	Page      int64           `json:"page"`      // 当前页码
	PageSize  int64           `json:"pageSize"`  // 每页大小
	TotalPage int64           `json:"totalPage"` // 总页数
	Total     int64           `json:"total"`     // 总项目数
	Data      []PublicProject `json:"data"`      // 项目列表
}

// HandleSearchPublicProjects 处理搜索公共项目的请求
func HandleSearchPublicProjects(ctx context.Context, request mcp.CallToolRequest) (_ *mcp.CallToolResult, err error) {
	defer func() { deferlog.DebugErrorContext(ctx, err, "HandleSearchPublicProjects completed") }()

	projectName := mcp.ParseString(request, "projectName", "")
	page := int(mcp.ParseInt(request, "page", 1))
	pageSize := int(mcp.ParseInt(request, "pageSize", 20))
	order := mcp.ParseString(request, "order", "default")

	url := fmt.Sprintf("https://api.apifox.com/api/v1/public-projects?order=%s&pageSize=%d&projectName=%s&page=%d",
		order, pageSize, projectName, page)

	return handleApiRequest(ctx, request, func() (*ApiResponse, error) {
		resp, err := GetApiResponse[SearchProjectsResponse](url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to search public projects (projectName=%s, page=%d, pageSize=%d): %v", projectName, page, pageSize, err)
		}
		return resp, nil
	})
}
