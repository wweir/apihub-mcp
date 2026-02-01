// Package handlers 包含所有 API 接口的处理函数
package handlers

import (
	"context"
	"fmt"

	"log/slog"

	"git.tophant.com/scc/cyber-guard/tools/tool-call/apifox/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sower-proxy/deferlog/v2"
	"github.com/toon-format/toon-go"
)

// ApiResponse 通用API响应结构
type ApiResponse struct {
	Success bool        `json:"success"`        // 是否成功
	Data    interface{} `json:"data"`           // 响应数据
	Meta    interface{} `json:"meta,omitempty"` // 元数据
}

// RequestHandler 通用请求处理函数类型
type RequestHandler func() (*ApiResponse, error)

// successResult 创建成功响应
func successResult(data any) *mcp.CallToolResult {
	text, err := toon.MarshalString(data)
	if err != nil {
		slog.Error("Failed to marshal response data", "error", err)
		return errorResult(fmt.Sprintf("Failed to marshal response: %v", err))
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}
}

// errorResult 创建错误响应
func errorResult(text string) *mcp.CallToolResult {
	slog.Error("API request failed", "error", text)
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}
}

// handleApiRequest 处理API请求的通用函数
func handleApiRequest(ctx context.Context, request mcp.CallToolRequest, handler RequestHandler) (_ *mcp.CallToolResult, err error) {
	defer func() { deferlog.DebugErrorContext(ctx, err, "handleApiRequest completed") }()

	apiResp, err := handler()
	if err != nil {
		err = fmt.Errorf("HTTP request failed: %v", err)
		return errorResult(err.Error()), nil
	}

	if !apiResp.Success {
		err = fmt.Errorf("API request failed")
		return errorResult(err.Error()), nil
	}

	return successResult(apiResp.Data), nil
}

// GetApiResponse 发送GET请求并解析响应
func GetApiResponse[T any](url string, headers map[string]string) (*ApiResponse, error) {
	type Response struct {
		Success bool `json:"success"`
		Data    T    `json:"data"`
		Meta    any  `json:"meta,omitempty"`
	}

	resp, err := client.GetJSON[Response](url, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %v", url, err)
	}

	return &ApiResponse{
		Success: resp.Success,
		Data:    resp.Data,
		Meta:    resp.Meta,
	}, nil
}
