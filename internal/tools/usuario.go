package tools

import (
	"context"

	"github.com/wellingtondevpro/artemis-mcp-go/internal/services"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterUsuarioTools(s *server.MCPServer) {
	cadastrarTool := mcp.NewTool("rh_cadastrar_funcionario",
		mcp.WithDescription("Use esta ferramenta para CADASTRAR/CRIAR um novo funcionário no módulo de Recursos Humanos (RH). Requer nome, cargo e salário."),
		mcp.WithString("nome", mcp.Required(), mcp.Description("Nome do funcionário")),
		mcp.WithString("cargo", mcp.Required(), mcp.Description("Cargo do funcionário (ex: Desenvolvedor, Gerente)")),
		mcp.WithNumber("salario", mcp.Required(), mcp.Description("Salário bruto mensal em formato numérico")),
	)
	s.AddTool(cadastrarTool, handleCadastrar)

	listarTool := mcp.NewTool("rh_listar_funcionarios",
		mcp.WithDescription("Use esta ferramenta para LISTAR TODOS os funcionários do RH. Utilize isso se o usuário pedir para listar todos em geral."),
	)
	s.AddTool(listarTool, handleListarUsuarios)

	buscarPorCargoTool := mcp.NewTool("rh_buscar_por_cargo",
		mcp.WithDescription("Use esta ferramenta para BUSCAR FUNCIONÁRIOS POR CARGO ou departamento (ex: desenvolvedor, gerente). Retorna uma lista de funcionários que têm aquele cargo específico."),
		mcp.WithString("cargo", mcp.Required(), mcp.Description("O cargo a ser pesquisado (ex: Desenvolvedor, Analista, Gerente)")),
	)
	s.AddTool(buscarPorCargoTool, handleBuscarPorCargo)

	resolverIdTool := mcp.NewTool("rh_resolver_id_por_nome",
		mcp.WithDescription("Uso interno: Resolve o nome do funcionário para o seu ID numérico real consultando o banco de dados."),
		mcp.WithString("nome", mcp.Required(), mcp.Description("Nome do funcionário (ex: Pedro)")),
	)
	s.AddTool(resolverIdTool, handleBuscarIdPorNome)

	consultarTool := mcp.NewTool("rh_consultar_funcionario",
		mcp.WithDescription("Use esta ferramenta para CONSULTAR os dados de um único funcionário específico no RH. EXIGE O ID numérico do funcionário. Não use se tiver apenas o nome."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("ID numérico obrigatório do funcionário")),
	)
	s.AddTool(consultarTool, handleConsultarUsuario)

	deletarTool := mcp.NewTool("rh_deletar_funcionario",
		mcp.WithDescription("AÇÃO DESTRUTIVA. Use esta ferramenta APENAS para DELETAR/EXCLUIR um funcionário do RH. Exige o ID numérico."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("ID numérico do funcionário a ser deletado")),
	)
	s.AddTool(deletarTool, handleDeletarUsuario)

	atualizarTool := mcp.NewTool("rh_atualizar_funcionario",
		mcp.WithDescription("Use esta ferramenta para ATUALIZAR/EDITAR os dados de um funcionário no RH (como mudar cargo, salário ou nome). Exige o ID numérico."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("ID numérico do funcionário")),
		mcp.WithString("nome", mcp.Required(), mcp.Description("Nome atualizado do funcionário")),
		mcp.WithString("cargo", mcp.Required(), mcp.Description("Cargo atualizado do funcionário")),
		mcp.WithNumber("salario", mcp.Required(), mcp.Description("Salário atualizado em formato numérico")),
	)
	s.AddTool(atualizarTool, handleAtualizarUsuario)
}

func handleCadastrar(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Argumentos inválidos"), nil
	}

	nome, _ := args["nome"].(string)
	cargo, _ := args["cargo"].(string)
	salario, _ := args["salario"].(float64)

	res, err := services.CadastrarUsuario(nome, cargo, salario)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}

func handleListarUsuarios(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	res, err := services.ListarUsuarios()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}

func handleBuscarPorCargo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Argumentos inválidos"), nil
	}

	cargo, _ := args["cargo"].(string)
	res, err := services.BuscarPorCargo(cargo)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}

func handleBuscarIdPorNome(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Argumentos inválidos"), nil
	}

	nome, _ := args["nome"].(string)
	res, err := services.BuscarIdPorNome(nome)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}

func handleConsultarUsuario(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Argumentos inválidos"), nil
	}

	id, _ := args["id"].(float64)
	res, err := services.ConsultarUsuario(int(id))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}

func handleDeletarUsuario(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Argumentos inválidos"), nil
	}

	id, _ := args["id"].(float64)
	res, err := services.DeletarUsuario(int(id))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}

func handleAtualizarUsuario(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Argumentos inválidos"), nil
	}

	id, _ := args["id"].(float64)
	nome, _ := args["nome"].(string)
	cargo, _ := args["cargo"].(string)
	salario, _ := args["salario"].(float64)

	res, err := services.AtualizarUsuario(int(id), nome, cargo, salario)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}
