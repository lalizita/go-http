# Go HTTP Server

Um servidor HTTP simples implementado do zero em Go, sem dependências externas. Este projeto demonstra como construir um servidor HTTP básico usando apenas a biblioteca padrão do Go, manipulando conexões TCP diretamente.

## 🚀 Características

- Implementação HTTP do zero usando TCP
- Sistema de roteamento simples
- Suporte a múltiplos métodos HTTP (GET, POST, etc.)
- Parsing de headers e body
- Respostas HTTP formatadas

## 📋 Pré-requisitos

- Go 1.16 ou superior

## 🔧 Como executar

1. Clone o repositório ou navegue até o diretório do projeto:
```bash
cd go-http
```

2. Execute o servidor:
```bash
go run main.go
```

O servidor estará rodando em `http://localhost:8080`

## 📖 Como usar

### Rotas disponíveis

Atualmente, o servidor possui a seguinte rota configurada:

- **GET /home** - Retorna uma página HTML simples

### Exemplo de uso

Com o servidor rodando, você pode testar usando `curl`:

```bash
# Testar a rota /home
curl http://localhost:8080/home

# Ver os headers da resposta
curl -i http://localhost:8080/home

# Testar uma rota inexistente (retorna 404)
curl http://localhost:8080/notfound
```

Ou abra no navegador: `http://localhost:8080/home`

## 🏗️ Estrutura do Projeto

```
go-http/
├── main.go      # Código principal do servidor HTTP
└── README.md    # Este arquivo
```

## 📝 Adicionando novas rotas

Para adicionar uma nova rota, edite a função `setupRoutes()` no arquivo `main.go`:

```go
func setupRoutes() {
    // Rota existente
    routes["GET /home"] = func(req HTTPRequest) HTTPResponse {
        return HTTPResponse{
            StatusCode: 200,
            Headers:    map[string]string{"Content-Type": "text/html"},
            Body:       "<h1>Home</h1>",
        }
    }
    
    // Nova rota GET
    routes["GET /about"] = func(req HTTPRequest) HTTPResponse {
        return HTTPResponse{
            StatusCode: 200,
            Headers:    map[string]string{"Content-Type": "text/html"},
            Body:       "<h1>Sobre</h1><p>Esta é a página sobre.</p>",
        }
    }
    
    // Nova rota POST
    routes["POST /api/data"] = func(req HTTPRequest) HTTPResponse {
        return HTTPResponse{
            StatusCode: 200,
            Headers:    map[string]string{"Content-Type": "application/json"},
            Body:       `{"message": "Dados recebidos", "body": "` + req.Body + `"}`,
        }
    }
}
```

## 🔍 Como funciona

1. **Inicialização**: O servidor cria um listener TCP na porta 8080
2. **Conexão**: Para cada nova conexão, uma goroutine é criada para processar a requisição
3. **Parsing**: A requisição HTTP é parseada linha por linha:
   - Linha de requisição (método, path, versão)
   - Headers
   - Body (se presente)
4. **Roteamento**: O servidor busca um handler correspondente ao método e path
5. **Resposta**: A resposta HTTP é formatada e enviada de volta ao cliente

## 📊 Códigos de Status Suportados

- `200 OK` - Requisição bem-sucedida
- `404 Not Found` - Rota não encontrada
- `500 Internal Server Error` - Erro interno do servidor

## 🎯 Objetivos de Aprendizado

Este projeto é ideal para entender:
- Como o protocolo HTTP funciona em baixo nível
- Manipulação de conexões TCP em Go
- Parsing de protocolos de texto
- Estrutura de requisições e respostas HTTP
