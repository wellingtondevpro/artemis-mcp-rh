package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getUsuarioApiUrl() string {
	url := os.Getenv("RH_API_URL")
	if url == "" {
		return "http://localhost:8081/api/funcionarios"
	}
	return url
}

func CadastrarUsuario(nome, cargo string, salario float64) (string, error) {
	fmt.Printf("[MCP LOG] Recebida tool call para CadastrarUsuario: %s, %s, %.2f\n", nome, cargo, salario)
	payload := map[string]interface{}{
		"nome":    nome,
		"cargo":   cargo,
		"salario": salario,
	}
	bodyBytes, _ := json.Marshal(payload)

	url := getUsuarioApiUrl()
	fmt.Printf("[MCP LOG] Fazendo POST para a RH-API em: %s\n", url)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		fmt.Printf("[MCP LOG] Erro fatal no POST: %v\n", err)
		return "", fmt.Errorf("erro ao comunicar com a API Java (RH): %v", err)
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)
	fmt.Printf("[MCP LOG] Resposta da RH-API: Status %d - %s\n", resp.StatusCode, string(respData))
	return fmt.Sprintf("Funcionário cadastrado com sucesso!\nResposta da API: %s", string(respData)), nil
}

func ListarUsuarios() (string, error) {
	url := getUsuarioApiUrl()
	fmt.Printf("[MCP LOG] Fazendo GET para listar usuarios em: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("erro ao comunicar com a API Java (RH): %v", err)
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)
	return fmt.Sprintf("Lista de Funcionários:\n%s", string(respData)), nil
}

func BuscarPorCargo(cargo string) (string, error) {
	urlStr := fmt.Sprintf("%s?cargo=%s", getUsuarioApiUrl(), url.QueryEscape(cargo))
	fmt.Printf("[MCP LOG] Fazendo GET para buscar usuarios por cargo (%s) em: %s\n", cargo, urlStr)
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("erro ao comunicar com a API Java (RH): %v", err)
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)
	return fmt.Sprintf("Funcionários com cargo '%s':\n%s", cargo, string(respData)), nil
}

func BuscarIdPorNome(nome string) (string, error) {
	urlStr := fmt.Sprintf("%s?nome=%s", getUsuarioApiUrl(), url.QueryEscape(nome))
	fmt.Printf("[MCP LOG] Fazendo GET para buscar usuario por nome (%s) em: %s\n", nome, urlStr)
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("erro ao comunicar com a API Java (RH): %v", err)
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)
	
	var encontrados []map[string]interface{}
	if err := json.Unmarshal(respData, &encontrados); err != nil {
		return "", fmt.Errorf("erro ao decodificar JSON da API Java: %v", err)
	}
	
	if len(encontrados) == 0 {
		return "NÃO ENCONTRADO", nil
	}
	if len(encontrados) > 1 {
		return "AMBIGUIDADE: Múltiplos funcionários com esse nome.", nil
	}
	
	idStr := fmt.Sprintf("%v", encontrados[0]["id"])
	// Remove decimais caso o JSON decodifique como float64
	if strings.HasSuffix(idStr, ".0") {
		idStr = strings.TrimSuffix(idStr, ".0")
	}
	return fmt.Sprintf("ID: %s", idStr), nil
}

func ConsultarUsuario(id int) (string, error) {
	url := fmt.Sprintf("%s/%d", getUsuarioApiUrl(), id)
	fmt.Printf("[MCP LOG] Fazendo GET para consultar usuario %d em: %s\n", id, url)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("erro ao comunicar com a API Java (RH): %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "Funcionário não encontrado.", nil
	}

	respData, _ := io.ReadAll(resp.Body)
	return fmt.Sprintf("Dados do funcionário:\n%s", string(respData)), nil
}

func DeletarUsuario(id int) (string, error) {
	url := fmt.Sprintf("%s/%d", getUsuarioApiUrl(), id)
	fmt.Printf("[MCP LOG] Fazendo DELETE para remover usuario %d em: %s\n", id, url)
	reqHttp, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		return "", fmt.Errorf("erro ao comunicar com a API Java (RH): %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "Funcionário não encontrado para deleção.", nil
	}

	return "Funcionário deletado com sucesso.", nil
}

func AtualizarUsuario(id int, nome, cargo string, salario float64) (string, error) {
	url := fmt.Sprintf("%s/%d", getUsuarioApiUrl(), id)
	fmt.Printf("[MCP LOG] Fazendo PUT para atualizar usuario %d em: %s\n", id, url)
	
	payload := map[string]interface{}{
		"nome":    nome,
		"cargo":   cargo,
		"salario": salario,
	}
	bodyBytes, _ := json.Marshal(payload)

	reqHttp, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisicao PUT: %v", err)
	}
	reqHttp.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		return "", fmt.Errorf("erro ao comunicar com a API Java (RH): %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "Funcionário não encontrado para atualização.", nil
	}

	respData, _ := io.ReadAll(resp.Body)
	return fmt.Sprintf("Funcionário atualizado com sucesso!\nNovos dados: %s", string(respData)), nil
}
