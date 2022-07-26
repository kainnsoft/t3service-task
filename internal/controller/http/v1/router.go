package v1

import (
	"net/http"
	gp "team3-task/internal/controller/grpc"
	"team3-task/pkg/logging"
)

func NewRouter(mux *http.ServeMux, t TaskHandlerInterface, grpcClient *gp.GrpcClient, log *logging.ZeroLogger) {
	// Swagger
	//swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	//handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	//handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	//handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	{
		NewTaskRouter(mux, t, grpcClient, log)
	}
}
