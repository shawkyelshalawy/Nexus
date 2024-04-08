package http

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"
)

type HttpRequest struct {
	Method   string
	Path     string
	Protocol string
	Headers  map[string]string
	Body     string
}
type HttpStatus string

var dir string

const (
	StatusOk                  HttpStatus = "200 OK"
	StatusNotFound            HttpStatus = "404 Not Found"
	StatusInternalServerError HttpStatus = "500 Internal Server Error"
)

const (
	textPlainContentType   = "text/plain"
	octetStreamContentType = "application/octet-stream"
)
const (
	httpProtocol = "HTTP"
	httpVersion  = "1.1"
	crlf         = "\r\n\r\n"
	lf           = "\n"
)

func ParseRequest(request string) (req HttpRequest) {
	reqDetails := strings.Split(request, "\r\n")
	staus := strings.Split(reqDetails[0], " ")
	req = HttpRequest{
		Method:   staus[0],
		Path:     staus[1],
		Protocol: staus[2],
	}
	req.Headers = make(map[string]string)
	for i := 1; i < len(reqDetails); i++ {
		if reqDetails[i] == "" {
			req.Body = reqDetails[i+1]
			break
		}
		header := strings.Split(reqDetails[i], ": ")
		req.Headers[header[0]] = header[1]
	}
	return
}

type HttpServer struct {
	conn net.Conn
}

func (http HttpServer) ReadRequest() HttpRequest {
	buf := make([]byte, 4096)
	content, err := http.conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}
	req := string(buf[:content])
	fmt.Println(req)
	return ParseRequest(req)
}

func (http HttpServer) getStatusLine(httpStatus HttpStatus) string {
	return fmt.Sprintf("%s/%s %s", httpProtocol, httpVersion, httpStatus)
}
func (http HttpServer) getContentLines(content string) string {
	res := fmt.Sprintf("Content-Type: text/plain%s", lf)
	res += fmt.Sprintf("Content-Length: %d%s", len(content), lf)
	res += fmt.Sprintf("%s%s%s", lf, content, crlf)
	return res
}
func (http HttpServer) Respond(status HttpStatus) {
	http.RespondWithContent(status, nil)
}

func (http HttpServer) RespondWithContent(status HttpStatus, content *string) {
	res := http.getStatusLine(status)
	if content != nil {
		res += lf + http.getContentLines(*content)
	} else {
		res += crlf
	}
	_, err := http.conn.Write([]byte(res))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
	http.conn.Close()
}

func (http HttpServer) handleFileRequest(filename string) {
	filePath := path.Join(dir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Respond(StatusNotFound)
		} else {
			fmt.Println("Error opening file: ", err.Error())
			http.Respond(StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file: ", err.Error())
		http.Respond(StatusInternalServerError)
		return
	}
	contentStr := string(content)
	http.RespondWithContent(StatusOk, &contentStr)
}

// accepting concurrent connections
func Listen(port string, dir *string) *HttpServer {
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go func() {
			http := &HttpServer{
				conn: conn,
			}
			req := http.ReadRequest()
			url := req.Path
			if req.Path == "/" {
				http.Respond(StatusOk)
			} else if strings.HasPrefix(req.Path, "/echo") {
				content := strings.ReplaceAll(req.Path, "/echo/", "")
				http.RespondWithContent(StatusOk, &content)
			} else if req.Path == "/user-agent" {
				content := req.Headers["User-Agent"]
				http.RespondWithContent(StatusOk, &content)
			} else if strings.Contains(url, "/files/") {
				fileName := url[7:]
				fmt.Println("file name is", fileName)
				if dir != nil {
					file := fmt.Sprintf("%s%s", *dir, fileName)
					data, readFileError := os.ReadFile(file)
					if readFileError != nil {
						_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n\r\n"))
					} else {
						content := string(data[:])
						_, err = conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length:%+v\r\n\r\n%v\r\n", len(content), content)))
					}
				} else {
					_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				}
			} else {
				http.Respond(StatusNotFound)
			}
		}()

	}
}
