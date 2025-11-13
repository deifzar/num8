package api8

import (
	"database/sql"
	"net/http"

	"deifzar/num8/pkg/cleanup8"
	"deifzar/num8/pkg/configparser"
	"deifzar/num8/pkg/controller8"
	"deifzar/num8/pkg/db8"
	"deifzar/num8/pkg/log8"
	"deifzar/num8/pkg/orchestrator8"

	"github.com/gin-gonic/gin"

	"os"
	"time"

	"github.com/spf13/viper"
)

type Api8 struct {
	Cnfg   *viper.Viper
	DB     *sql.DB
	Router *gin.Engine
}

func (a *Api8) Init() error {
	// Create configs, log and tmp directories if they don't exist
	if err := os.MkdirAll("configs", 0750); err != nil {
		log8.BaseLogger.Error().Err(err).Msg("Failed to create configs directory")
		return err
	}
	if err := os.MkdirAll("log", 0750); err != nil {
		log8.BaseLogger.Error().Err(err).Msg("Failed to create log directory")
		return err
	}
	if err := os.MkdirAll("tmp", 0750); err != nil {
		log8.BaseLogger.Error().Err(err).Msg("Failed to create tmp directory")
		return err
	}

	// Clean up old files in tmp directory (older than 24 hours)
	cleanup := cleanup8.NewCleanup8()
	if err := cleanup.CleanupDirectory("tmp", 24*time.Hour); err != nil {
		log8.BaseLogger.Error().Err(err).Msg("Failed to cleanup tmp directory")
		// Don't return error here as cleanup failure shouldn't prevent startup
	}
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

	a.Cnfg = v
	a.DB = conn
	return nil
}

// InitializeConsumerAfterReady starts a goroutine that waits for the API service
// to become ready (via /ready endpoint) before initializing RabbitMQ queues and consumers.
// This prevents consumers from receiving messages before the API can handle them.
func (a *Api8) InitializeConsumerAfterReady() {
	go func() {
		locationService := a.Cnfg.GetString("ORCHESTRATORM8.Services.num8")
		requestURL := locationService + "/health"

		log8.BaseLogger.Info().Msg("Waiting for API service to become ready before activating RabbitMQ consumer...")

		// Poll the /ready endpoint until the service is healthy
		maxRetries := 60 // 5 minutes total (60 * 5 seconds)
		retryCount := 0
		for {
			resp, err := http.Get(requestURL)
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				log8.BaseLogger.Info().Msg("API service is ready. Initializing RabbitMQ consumer...")
				break
			}
			if resp != nil {
				resp.Body.Close()
			}

			retryCount++
			if retryCount >= maxRetries {
				log8.BaseLogger.Error().Msg("Timeout waiting for API service to become ready. Consumer will not be activated.")
				return
			}

			time.Sleep(5 * time.Second)
		}

		// Initialize RabbitMQ orchestrator
		orchestrator8, err := orchestrator8.NewOrchestrator8()
		if err != nil {
			log8.BaseLogger.Error().Err(err).Msg("Error connecting to the RabbitMQ server.")
			return
		}

		err = orchestrator8.InitOrchestrator()
		if err != nil {
			log8.BaseLogger.Error().Err(err).Msg("Error bringing up the RabbitMQ exchanges.")
			return
		}

		err = orchestrator8.ActivateQueueByService("num8")
		if err != nil {
			log8.BaseLogger.Error().Err(err).Msg("Error bringing up the RabbitMQ queues for the `num8` service.")
			return
		}

		err = orchestrator8.ActivateConsumerByService("num8")
		if err != nil {
			log8.BaseLogger.Error().Err(err).Msg("Error activating consumer with dedicated connection for the `num8` service.")
			return
		}

		log8.BaseLogger.Info().Msg("RabbitMQ consumer successfully activated for num8 service.")
	}()
}

func (a *Api8) Routes() {
	r := gin.Default()
	contrNumate := controller8.NewController8Numate(a.DB, a.Cnfg)
	r.MaxMultipartMemory = 8 << 20 //8 MiB
	r.GET("/scan", contrNumate.NumateScan)
	r.POST("/domain/:id/scan", contrNumate.NumateDomain)
	r.POST("/hostname/:id/scan", contrNumate.NumateHostname)
	r.POST("/endpoint/:id/scan", contrNumate.NumateEndpoint)

	// Health live probes
	r.GET("/health", contrNumate.HealthCheck)
	r.GET("/ready", contrNumate.ReadinessCheck)

	a.Router = r
}

func (a *Api8) Run(addr string) {
	a.Router.Run(addr)
}
