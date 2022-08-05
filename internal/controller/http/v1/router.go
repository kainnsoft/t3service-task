package v1

import (
	"net/http"
	gp "team3-task/internal/controller/grpc"
	"team3-task/pkg/logging"
)

func NewRouter(mux *http.ServeMux, t TaskHandlerInterface, grpcClient *gp.GClient, log *logging.ZeroLogger) {
	// Routers
	NewTaskRouter(mux, t, grpcClient, log)
}
