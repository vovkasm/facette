package server

import "net/http"

func (server *Server) serveManage(writer http.ResponseWriter, request *http.Request) {
	// TODO: handle HTTP authorization before anything else

	// Dispatch management API routes
	if routeMatch(request.URL.Path, urlManagePath+"library/backup") {
		server.serveManageLibraryBackup(writer, request)
	} else if routeMatch(request.URL.Path, urlManagePath+"library/refresh") {
		server.serveManageLibraryRefresh(writer, request)
	} else {
		server.serveResponse(writer, nil, http.StatusNotFound)
	}
}
