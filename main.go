package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

const (
	StatusOK                  = 200
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

var statusTexts = map[int]string{
	StatusOK:                  "OK",
	StatusNotFound:            "Not Found",
	StatusInternalServerError: "Internal Server Error",
}

type HTTPRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
}

type HTTPResponse struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       string
}

type RouteHandler func(HTTPRequest) HTTPResponse

var routes = make(map[string]RouteHandler)

func main() {
	setupRoutes()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("HTTP Server is running on port 8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func setupRoutes() {
	routes["GET /home"] = func(req HTTPRequest) HTTPResponse {
		return HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "text/html"},
			Body:       "<h1>Home</h1>",
		}
	}

	routes["POST /api/data"] = func(req HTTPRequest) HTTPResponse {
		return HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"message": "Dados recebidos", "body": "` + req.Body + `"}`,
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	// textproto implementa um suporte genérico para protocolos de texto, como HTTP.
	tp := textproto.NewReader(reader)

	requestLine, err := tp.ReadLine()
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Println("Error reading request line:", err)
		}
		return
	}

	log.Println("Request line:", requestLine)

	// Divide a linha de requisição em método, caminho e versão do protocolo
	parts := strings.SplitN(requestLine, " ", 3)
	if len(parts) != 3 {
		log.Println("Invalid request line:", requestLine)
		return
	}
	method, path := parts[0], parts[1]

	headers := make(map[string]string)
	// Lê todas as linhas de header até encontrar uma linha vazia.
	for {
		headerLine, err := tp.ReadLine()
		if err != nil {
			log.Println("Error reading header line:", err)
			break
		}

		if headerLine == "" {
			break
		}

		parts := strings.SplitN(headerLine, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	log.Println("Headers:", headers)

	body := ""
	if headers["Content-Length"] != "" {
		contentLength, err := strconv.Atoi(headers["Content-Length"])
		if err != nil {
			log.Println("Error parsing Content-Length:", err)
			return
		}

		bodyBytes := make([]byte, contentLength)
		_, err = reader.Read(bodyBytes)
		if err != nil {
			log.Println("Error reading body:", err)
			return
		}

		body = string(bodyBytes)
	}

	request := HTTPRequest{
		Method:  method,
		Path:    path,
		Headers: headers,
		Body:    body,
	}

	var response HTTPResponse
	if handler, ok := routes[fmt.Sprintf("%s %s", method, path)]; ok {
		response = handler(request)
	} else {
		response = HTTPResponse{
			StatusCode: StatusNotFound,
			StatusText: statusTexts[StatusNotFound],
			Headers:    map[string]string{"Content-Type": "text/html"},
			Body:       "<h1>Not Found</h1>",
		}
	}

	responseString := fmt.Sprintf("HTTP/1.1 %d %s\r\n", response.StatusCode, response.StatusText)
	for key, value := range response.Headers {
		responseString += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	responseString += "\r\n"
	if len(response.Body) > 0 {
		responseString += response.Body
	}

	conn.Write([]byte(responseString))
}
