package initializations

import (
	"bmt_payment_service/global"
	"fmt"
)

func Run() {
	loadConfigs()
	initPostgreSql()
	initRedis()

	r := initRouter()

	r.Run(fmt.Sprintf("0.0.0.0:%s", global.Config.Server.ServerPort))
}
