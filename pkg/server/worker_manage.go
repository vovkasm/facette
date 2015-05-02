package server

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/facette/facette/pkg/logger"
	"github.com/facette/facette/pkg/worker"
)

const (
	urlManagePath string = "/api/v1/manage/"
)

func workerManageInit(w *worker.Worker, args ...interface{}) {
	var server = args[0].(*Server)

	logger.Log(logger.LevelDebug, "manageWorker", "init")

	// Worker properties:
	// 0: server instance (*Server)
	w.Props = append(w.Props, server)

	w.ReturnErr(nil)
}

func workerManageShutdown(w *worker.Worker, args ...interface{}) {
	logger.Log(logger.LevelDebug, "manageWorker", "shutdown")

	w.SendJobSignal(jobSignalShutdown)

	w.ReturnErr(nil)
}

func workerManageRun(w *worker.Worker, args ...interface{}) {
	var server = w.Props[0].(*Server)

	defer w.Shutdown()

	logger.Log(logger.LevelDebug, "manageWorker", "starting")

	// Prepare router
	router := NewRouter(server)

	router.HandleFunc(urlManagePath, server.serveManage)

	http.Handle(urlManagePath, router)

	// TODO: refactor listener setup as serveWorker's code is identical

	// Start serving HTTP requests
	netType := "tcp"
	address := server.Config.MgmtBindAddr
	for _, scheme := range [...]string{"tcp", "tcp4", "tcp6", "unix"} {
		prefix := scheme + "://"

		if strings.HasPrefix(address, prefix) {
			netType = scheme
			address = strings.TrimPrefix(address, prefix)
			break
		}
	}

	listener, err := net.Listen(netType, address)
	if err != nil {
		w.ReturnErr(err)
		return
	}

	logger.Log(logger.LevelInfo, "manageWorker", "listening on %s", server.Config.MgmtBindAddr)

	if netType == "unix" {
		// Change owning user and group
		if server.Config.MgmtSocketUser >= 0 || server.Config.MgmtSocketGroup >= 0 {
			logger.Log(logger.LevelDebug, "manageWorker", "changing ownership of unix socket to UID %v and GID %v",
				server.Config.MgmtSocketUser, server.Config.MgmtSocketGroup)
			err = os.Chown(address, server.Config.MgmtSocketUser, server.Config.MgmtSocketGroup)
			if err != nil {
				listener.Close()
				w.ReturnErr(err)
				return
			}
		}

		// Change mode
		if server.Config.MgmtSocketMode != nil {
			mode, err := strconv.ParseUint(*server.Config.MgmtSocketMode, 8, 32)
			if err != nil {
				logger.Log(logger.LevelError, "manageWorker", "manage_socket_mode is invalid")
				listener.Close()
				w.ReturnErr(err)
				return
			}

			logger.Log(logger.LevelDebug, "manageWorker", "changing file permissions mode of unix socket to %04o", mode)
			err = os.Chmod(address, os.FileMode(mode))
			if err != nil {
				listener.Close()
				w.ReturnErr(err)
				return
			}
		}
	}

	go http.Serve(listener, nil)

	for {
		select {
		case cmd := <-w.ReceiveJobSignals():
			switch cmd {
			case jobSignalShutdown:
				logger.Log(logger.LevelInfo, "manageWorker", "received shutdown command, stopping job")

				listener.Close()

				logger.Log(logger.LevelInfo, "manageWorker", "server listener closed")

				w.State = worker.JobStopped

				return

			default:
				logger.Log(logger.LevelInfo, "manageWorker", "received unknown command, ignoring")
			}
		}
	}

	w.ReturnErr(nil)
}
