package mcp

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
)

// WithRecovery wraps an MCP tool handler and intercepts any panics that occur during execution.
// It converts them into safe, formatted mcp.ToolResultError structs to ensure the stdio connection
// isn't broken violently upon nil-pointer dereferences or logical panics inside Audio Engine/Providers.
func WithRecovery(handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in tool execution: %v\n", r)

				// Ensure the JSON-RPC response conveys the error without killing the underlying IO pipe
				errMsg := fmt.Sprintf("Internal Server Error: Panic recovered during tool execution: %v", r)
				res = mcp.NewToolResultError(errMsg)
				err = nil // We explicitly swallow the error so it serializes over JSON successfully
			}
		}()
		return handler(ctx, request)
	}
}
