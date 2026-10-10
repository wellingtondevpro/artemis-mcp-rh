package tools

import (
	"context"

	"github.com/wellingtondevpro/artemis-mcp-go/internal/services"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterBeneficiosTools(s *server.MCPServer) {
	listarTool := mcp.NewTool("rh_listar_beneficios",
		mcp.WithDescription("Use esta ferramenta para LISTAR TODOS os benefícios disponíveis no catálogo da empresa. Ideal para consultar o que pode ser oferecido."),
	)
	s.AddTool(listarTool, handleListarBeneficios)

	atribuirTool := mcp.NewTool("rh_atribuir_beneficio",
		mcp.WithDescription("Use esta ferramenta para ATRIBUIR/VINCULAR um benefício (ex: Plano de Saúde, Vale Alimentação) a um funcionário. Exige o ID numérico do funcionário e o nome do benefício."),
		mcp.WithNumber("usuario_id", mcp.Required(), mcp.Description("ID numérico do funcionário")),
		mcp.WithString("beneficio", mcp.Required(), mcp.Description("Nome do benefício (ex: Vale Transporte, Gympass)")),
	)
	s.AddTool(atribuirTool, handleAtribuirBeneficio)
}

func handleListarBeneficios(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	res, err := services.ListarBeneficios()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}

func handleAtribuirBeneficio(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Argumentos inválidos"), nil
	}

	usuarioId, _ := args["usuario_id"].(float64)
	beneficio, _ := args["beneficio"].(string)

	res, err := services.AtribuirBeneficio(int(usuarioId), beneficio)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}
