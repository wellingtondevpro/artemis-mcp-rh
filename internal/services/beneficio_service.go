package services

import (
	"fmt"
)

// Simulando um serviço de benefícios

func ListarBeneficios() (string, error) {
	return "Benefícios disponíveis: Vale Refeição, Plano de Saúde, Gympass, Auxílio Home Office.", nil
}

func AtribuirBeneficio(usuarioId int, beneficio string) (string, error) {
	return fmt.Sprintf("Benefício '%s' atribuído ao usuário ID %d com sucesso.", beneficio, usuarioId), nil
}
