package routers

import (
	"bmt_payment_service/db/sqlc"
	"bmt_payment_service/global"
	"bmt_payment_service/internal/controllers"
	"bmt_payment_service/internal/implementations/momo"

	"github.com/gin-gonic/gin"
)

type MoMoRouter struct {
}

func (m *MoMoRouter) InitMoMoRouter(router *gin.RouterGroup) {
	sqlStore := sqlc.NewStore(global.Postgresql)
	moMoService := momo.NewMomoPayment(
		sqlStore,
		global.Config.ServiceSetting.MoMoSetting.EndPoint,
		global.Config.ServiceSetting.MoMoSetting.PartnerCode,
		global.Config.ServiceSetting.MoMoSetting.AccessKey,
		global.Config.ServiceSetting.MoMoSetting.SecretKey,
		global.Config.ServiceSetting.MoMoSetting.RedirectURL,
		global.Config.ServiceSetting.MoMoSetting.IPNURL,
	)
	moMoController := controllers.NewMoMoController(moMoService)

	moMoRouter := router.Group("/momo")
	{
		moMoCustomerRouter := moMoRouter.Group("/customer")
		{
			moMoCustomerRouter.POST("/create_payment_url", moMoController.CreatePaymentURL)
			moMoCustomerRouter.GET("/verify_payment", moMoController.VerifyPaymentCallback)
		}
	}
}
