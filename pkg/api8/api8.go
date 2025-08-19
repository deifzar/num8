package api8

import (
	"database/sql"

	"deifzar/num8/pkg/configparser"
	"deifzar/num8/pkg/controller8"
	"deifzar/num8/pkg/db8"
	"deifzar/num8/pkg/log8"
	"deifzar/num8/pkg/orchestrator8"

	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
)

type Api8 struct {
	Cnfg   *viper.Viper
	DB     *sql.DB
	Router *gin.Engine
}

func (a *Api8) Init() error {
	v, err := configparser.InitConfigParser()
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Info().Msg("Error initialising the config parser.")
		return err
	}
	location := v.GetString("Database.location")
	port := v.GetInt("Database.port")
	schema := v.GetString("Database.schema")
	database := v.GetString("Database.database")
	username := v.GetString("Database.username")
	password := v.GetString("Database.password")

	var db db8.Db8
	db.InitDatabase8(location, port, schema, database, username, password)
	conn, err := db.OpenConnection()
	if err != err {
		log8.BaseLogger.Error().Msg("Error connecting into DB.")
		return err
	}

	orchestrator8, err := orchestrator8.NewOrchestrator8()
	if err != nil {
		log8.BaseLogger.Error().Msg("Error connecting to the RabbitMQ server.")
		return err
	}
	err = orchestrator8.InitOrchestrator()
	if err != nil {
		log8.BaseLogger.Error().Msg("Error bringin up the RabbitMQ exchanges.")
		return err
	}
	err = orchestrator8.ActivateQueueByService("num8")
	if err != nil {
		log8.BaseLogger.Error().Msg("Error bringing up the RabbitMQ queues for the `num8` service.")
		return err
	}
	orchestrator8.CreateHandleAPICallByService("num8")
	orchestrator8.ActivateConsumerByService("num8")

	a.Cnfg = v
	a.DB = conn
	return nil
}

func (a *Api8) Routes() {
	r := gin.Default()
	contrNumate := controller8.NewController8Numate(a.DB, a.Cnfg)
	r.MaxMultipartMemory = 8 << 20 //8 MiB
	r.GET("/scan", contrNumate.NumateScan)
	r.POST("/domain/:id/scan", contrNumate.NumateDomain)
	r.POST("/hostname/:id/scan", contrNumate.NumateHostname)
	r.POST("/endpoint/:id/scan", contrNumate.NumateEndpoint)
	a.Router = r
}

func (a *Api8) Run(addr string) {
	a.Router.Run(addr)
}
