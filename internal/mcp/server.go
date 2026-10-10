package mcp_internal

import (
	"context"
	"fmt"
	"os"

	"github.com/wellingtondevpro/artemis-mcp-go/internal/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func StartServer(transport string, port int) {
	mcpServer := server.NewMCPServer(
		"artemis-mcp-rh",
		"1.0.1",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)

	// Registra tools usando os pacotes internos
	tools.RegisterUsuarioTools(mcpServer)
	tools.RegisterBeneficiosTools(mcpServer)

	// Registra prompts
	registerPrompts(mcpServer)

	// Registra resources
	registerResources(mcpServer)

	if transport == "sse" {
		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("Iniciando MCP Server de RH em Go (via SSE) na porta %d...\n", port)
		sseServer := server.NewSSEServer(mcpServer)
		if err := sseServer.Start(addr); err != nil {
			fmt.Fprintf(os.Stderr, "Erro no MCP Server (SSE): %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Iniciando MCP Server de RH em Go (via STDIO)...")
		if err := server.ServeStdio(mcpServer); err != nil {
			fmt.Fprintf(os.Stderr, "Erro no MCP Server (STDIO): %v\n", err)
			os.Exit(1)
		}
	}
}

func registerPrompts(s *server.MCPServer) {
	prompt := mcp.NewPrompt("rh_entrevista",
		mcp.WithPromptDescription("Um prompt para ajudar na condução de entrevistas de RH"),
		mcp.WithArgument("cargo", mcp.ArgumentDescription("O cargo sendo entrevistado"), mcp.RequiredArgument()),
	)
	s.AddPrompt(prompt, func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		cargo := ""
		if req.Params.Arguments != nil {
			if c, ok := req.Params.Arguments["cargo"]; ok {
				cargo = c
			}
		}

		return mcp.NewGetPromptResult("Uma entrevista padrão de RH", []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(fmt.Sprintf("Você é um especialista em RH entrevistando um candidato para a vaga de %s. Faça 3 perguntas essenciais.", cargo)),
			),
		}), nil
	})
}

func registerResources(s *server.MCPServer) {
	resource := mcp.NewResource("rh://politicas/codigo-conduta",
		"Código de Conduta",
		mcp.WithResourceDescription("Código de conduta e ética da empresa"),
		mcp.WithMIMEType("text/plain"),
	)

	s.AddResource(resource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content := "1. Respeito mútuo.\n2. Ética no trabalho.\n3. Pontualidade."
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "text/plain",
				Text:     content,
			},
		}, nil
	})
}
