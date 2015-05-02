package server

import (
	"net/http"

	"github.com/facette/facette/pkg/utils"
)

func (server *Server) serveManageLibraryBackup(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		return

	case "PUT":
		if utils.HTTPGetContentType(request) != "application/json" {
			server.serveResponse(writer, serverResponse{mesgUnsupportedMediaType}, http.StatusUnsupportedMediaType)
			return
		}
		return

	default:
		server.serveResponse(writer, serverResponse{mesgMethodNotAllowed}, http.StatusMethodNotAllowed)
	}
}

func (server *Server) serveManageLibraryRefresh(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		if err := server.Library.Refresh(); err != nil {
			server.serveResponse(writer, serverResponse{mesgUnhandledError}, http.StatusInternalServerError)
			return
		}

		server.serveResponse(writer, serverResponse{"library refreshed successfully"}, http.StatusOK)
		return

	default:
		server.serveResponse(writer, serverResponse{mesgMethodNotAllowed}, http.StatusMethodNotAllowed)
	}
}
