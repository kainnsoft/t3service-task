package v1

import (
	"net/http"
	app_interface "team3-task/internal/app/interface"
	"team3-task/pkg/logging"
)

func NewRouter(mux *http.ServeMux, t TaskHandlerInterface, grpcClient app_interface.AuthAccessChecker, log *logging.ZeroLogger) {
	// Routers
	NewTaskRouter(mux, t, grpcClient, log)
}
