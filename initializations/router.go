package initializations

import (
	"bmt_payment_service/internal/routers"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()

	// // Routers
	moMoRouter := routers.PaymentServiceRouterGroup.MoMo

	mainGroup := r.Group("/v1")
	{
		moMoRouter.InitMoMoRouter(mainGroup)
	}

	return r
}
