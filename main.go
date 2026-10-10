package main

import (
	"flag"

	mcp_internal "github.com/wellingtondevpro/artemis-mcp-go/internal/mcp"
)

func main() {
	transport := flag.String("transport", "stdio", "Transport to use: stdio or sse")
	port := flag.Int("port", 8080, "Port for SSE transport")
	flag.Parse()

	mcp_internal.StartServer(*transport, *port)
}
