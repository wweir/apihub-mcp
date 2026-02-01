// Package client 提供 HTTP 请求的封装
package client

import (
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	// restyClient Resty 客户端实例
	restyClient *resty.Client
)

func init() {
	// 初始化 Resty 客户端
	restyClient = resty.New()
	restyClient.SetHeader("User-Agent", "Apifox-MCP-Service/1.0")
	restyClient.SetHeader("Content-Type", "application/json")
	restyClient.SetTimeout(30 * time.Second) // 设置超时
}

// Get 发送 GET 请求
func Get(url string, headers map[string]string) (*resty.Response, error) {
	req := restyClient.R()
	for key, value := range headers {
		req.SetHeader(key, value)
	}

	resp, err := req.Get(url)
	if err != nil {
		slog.Error("GET request failed", "url", url, "error", err)
		return nil, err
	}

	return resp, nil
}

// Post 发送 POST 请求
func Post(url string, body any, headers map[string]string) (*resty.Response, error) {
	req := restyClient.R()
	for key, value := range headers {
		req.SetHeader(key, value)
	}

	resp, err := req.SetBody(body).Post(url)
	if err != nil {
		slog.Error("POST request failed", "url", url, "error", err)
		return nil, err
	}

	return resp, nil
}

// GetJSON 发送 GET 请求并直接解析响应为指定类型（泛型版本）
func GetJSON[T any](url string, headers map[string]string) (T, error) {
	var result T
	req := restyClient.R()
	for key, value := range headers {
		req.SetHeader(key, value)
	}

	_, err := req.SetResult(&result).Get(url)
	if err != nil {
		slog.Error("GET request failed", "url", url, "error", err)
		return result, err
	}

	return result, nil
}

// PostJSON 发送 POST 请求并直接解析响应为指定类型（泛型版本）
func PostJSON[T any](url string, body any, headers map[string]string) (T, error) {
	var result T
	req := restyClient.R()
	for key, value := range headers {
		req.SetHeader(key, value)
	}

	_, err := req.SetBody(body).SetResult(&result).Post(url)
	if err != nil {
		slog.Error("POST request failed", "url", url, "error", err)
		return result, err
	}

	return result, nil
}

// GetClient 获取 Resty 客户端实例
func GetClient() *resty.Client {
	return restyClient
}
