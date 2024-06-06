package main

// #include "preamble.h"
import "C"
import (
	"context"
	"fmt"
	"io"
	"net/http"
	"unsafe"

	"github.com/google/uuid"
)

type ClientContext struct {
	ID              string
	Request         *http.Request
	ResponseChannel chan *Response
}

type Response struct {
	ID         string
	StatusCode int
	Headers    string
	Content    []byte
}

type ServerOptions struct {
	TLSCertFile string
	TLSKeyFile  string
}

type Server struct {
	Error      error
	Port       int
	Options    ServerOptions
	HTTPServer *http.Server
	Started    bool
	Requests   chan C.struct_Request
	Clients    map[string]*ClientContext
}

var server *Server

func InitServer(port int, options ServerOptions) *Server {
	s := &Server{
		Port:    port,
		Options: options,
		HTTPServer: &http.Server{
			Addr:                         fmt.Sprintf(":%d", port),
			DisableGeneralOptionsHandler: true,
		},
		Requests: make(chan C.struct_Request, 1000),
		Clients:  make(map[string]*ClientContext),
	}
	s.HTTPServer.Handler = s

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req_ *http.Request) {
	// Process the request and send to the s.Requests channel.
	uid := uuid.New().String()
	ctx := &ClientContext{
		ID:              uid,
		Request:         req_,
		ResponseChannel: make(chan *Response),
	}
	s.Clients[uid] = ctx
	contentBuf, err := io.ReadAll(req_.Body)
	if err != nil {
		// Respond with status code 400.
		w.WriteHeader(400)
		w.Write([]byte("unable to read request body"))
	} else {
		req := C.struct_Request{
			ID:      C.CString(uid),
			Method:  C.CString(req_.Method),
			URL:     C.CString(req_.URL.String()),
			Headers: C.CString(HeadersToString(req_.Header)),
		}
		req.Content = (*C.char)(C.CBytes(contentBuf))
		req.ContentLength = C.size_t(len(contentBuf))

		s.Requests <- req

		// Wait for response, then send it out.
		select {
		case resp := <-ctx.ResponseChannel:
			StringToHeaders(resp.Headers, w.Header())
			w.WriteHeader(resp.StatusCode)
			w.Write(resp.Content)
		}
	}
}

func (s *Server) Start() {
	s.Started = true
	if s.Options.TLSCertFile != "" && s.Options.TLSKeyFile != "" {
		s.Error = s.HTTPServer.ListenAndServeTLS(s.Options.TLSCertFile, s.Options.TLSKeyFile)
	} else {
		s.Error = s.HTTPServer.ListenAndServe()
	}
	s.Started = false

	// Signal that shutdown is done.
	s.Requests <- C.struct_Request{
		ID: C.CString("SHUTDOWN"),
	}
}

func (s *Server) Stop() {
	s.HTTPServer.Shutdown(context.Background())
}

//export StartServer
func StartServer(port int, options C.struct_ServerOptions) C.enum_Status {
	if server != nil && server.Started {
		return C.E_SERVER_STARTED
	}

	opts := ServerOptions{}
	if options.TLSCertFile != nil {
		opts.TLSCertFile = C.GoString(options.TLSCertFile)
	}
	if options.TLSKeyFile != nil {
		opts.TLSKeyFile = C.GoString(options.TLSKeyFile)
	}

	server = InitServer(port, opts)
	go server.Start()
	return C.E_OK
}

//export ServerError
func ServerError() *C.char {
	if server != nil && server.Error != nil {
		return C.CString(server.Error.Error())
	} else {
		return nil
	}
}

//export StopServer
func StopServer() {
	if server != nil && server.Started {
		server.Stop()
	}
}

//export ReadRequest
func ReadRequest(outReq *C.struct_Request) C.enum_Status {
	select {
	case req := <-server.Requests:
		if outReq != nil {
			*outReq = req
		}
		return C.E_OK

	default:
		return C.E_NO_MSG
	}
}

//export SendResponse
func SendResponse(resp C.struct_Response) C.enum_Status {
	if server == nil || !server.Started {
		return C.E_SERVER_NOT_STARTED
	}
	respID := C.GoString(resp.ID)
	ctx := server.Clients[respID]

	ctx.ResponseChannel <- &Response{
		ID:         C.GoString(resp.ID),
		StatusCode: int(resp.StatusCode),
		Headers:    C.GoString(resp.Headers),
		Content:    C.GoBytes(unsafe.Pointer(resp.Content), C.int(resp.ContentLength)),
	}

	delete(server.Clients, respID)

	return C.E_OK
}
