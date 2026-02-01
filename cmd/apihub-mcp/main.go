package main

import (
	"flag"
	"log/slog"
	"os"

	"git.tophant.com/scc/cyber-guard/tools/tool-call/apifox/internal/handlers"
	"github.com/lmittmann/tint"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sower-proxy/deferlog/v2"
)

// NewMCPServer 创建并配置 MCP 服务器
func NewMCPServer() *server.MCPServer {
	mcpServer := server.NewMCPServer(
		"apifox-server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register search_public_projects tool
	mcpServer.AddTool(
		mcp.NewTool("search_public_projects",
			mcp.WithDescription("Search public projects on Apifox"),
			mcp.WithString("projectName",
				mcp.Description("Project name to search"),
				mcp.Required(),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number"),
				mcp.DefaultNumber(1),
			),
			mcp.WithNumber("pageSize",
				mcp.Description("Page size"),
				mcp.DefaultNumber(20),
			),
			mcp.WithString("order",
				mcp.Description("Sort order"),
				mcp.DefaultString("default"),
			),
		),
		handlers.HandleSearchPublicProjects,
	)

	// Register get_project_info tool
	mcpServer.AddTool(
		mcp.NewTool("get_project_info",
			mcp.WithDescription("Get project information"),
			mcp.WithNumber("projectId",
				mcp.Description("Project ID"),
				mcp.Required(),
			),
			mcp.WithString("locale",
				mcp.Description("Locale, e.g., zh-CN, en-US"),
				mcp.DefaultString("zh-CN"),
			),
		),
		handlers.HandleGetProjectInfo,
	)

	// Register get_sprint_branches tool
	mcpServer.AddTool(
		mcp.NewTool("get_sprint_branches",
			mcp.WithDescription("Get project sprint branches"),
			mcp.WithNumber("projectId",
				mcp.Description("Project ID"),
				mcp.Required(),
			),
			mcp.WithString("locale",
				mcp.Description("Locale, e.g., zh-CN, en-US"),
				mcp.DefaultString("zh-CN"),
			),
		),
		handlers.HandleGetSprintBranches,
	)

	// Register get_api_tree_list tool
	mcpServer.AddTool(
		mcp.NewTool("get_api_tree_list",
			mcp.WithDescription("Get project API tree list"),
			mcp.WithNumber("projectId",
				mcp.Description("Project ID"),
				mcp.Required(),
			),
			mcp.WithString("locale",
				mcp.Description("Locale, e.g., zh-CN, en-US"),
				mcp.DefaultString("zh-CN"),
			),
		),
		handlers.HandleGetAPITreeList,
	)

	// Register get_api_details tool
	mcpServer.AddTool(
		mcp.NewTool("get_api_details",
			mcp.WithDescription("Get API details"),
			mcp.WithNumber("projectId",
				mcp.Description("Project ID"),
				mcp.Required(),
			),
			mcp.WithNumber("branchId",
				mcp.Description("Branch ID"),
				mcp.Required(),
			),
			mcp.WithString("locale",
				mcp.Description("Locale, e.g., zh-CN, en-US"),
				mcp.DefaultString("zh-CN"),
			),
		),
		handlers.HandleGetAPIDetails,
	)

	return mcpServer
}

// StartServer 启动 MCP 服务器
func StartServer(mcpServer *server.MCPServer, addr string) error {
	if addr != "" {
		httpServer := server.NewStreamableHTTPServer(mcpServer)
		slog.Info("MCP HTTP server listening", "addr", addr+"/mcp")
		return httpServer.Start(addr)
	}

	slog.Info("MCP server starting with stdio transport")
	return server.ServeStdio(mcpServer)
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "", "HTTP server address (if not set, use stdio transport)")
	flag.Parse()

	fi, _ := os.Stdout.Stat()
	noColor := (fi.Mode() & os.ModeCharDevice) == 0
	deferlog.SetDefault(slog.New(tint.NewHandler(os.Stdout,
		&tint.Options{AddSource: true, NoColor: noColor})))

	// 创建并配置 MCP 服务器
	mcpServer := NewMCPServer()

	// 启动服务器
	if err := StartServer(mcpServer, addr); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}
