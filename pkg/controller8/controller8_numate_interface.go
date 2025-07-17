package controller8

import (
	"deifzar/num8/pkg/model8"
	"deifzar/num8/pkg/orchestrator8"

	"github.com/gin-gonic/gin"
)

type Controller8NumateInterface interface {
	NumateScan(*gin.Context)
	NumateDomain(*gin.Context)
	NumateHostname(*gin.Context)
	NumateEndpoint(*gin.Context)
	// ConfigureEngine(model8.PostOptionsScan8) (model8.Model8Options8Interface, model8.Model8Results8Interface, error)
	ConfigureEngine(model8.PostOptionsScan8) (model8.Model8Options8Interface, string, error)
	RunNumate(bool, orchestrator8.Orchestrator8Interface, []model8.Endpoint8, model8.Model8Options8Interface, string)
	CommitResults([]model8.SecurityIssues8, []model8.Endpoint8) (bool, error)
	// RunNumate([]model8.Httpendpoint8, model8.Model8Options8Interface, model8.Model8Results8Interface)
	// RabbitMQBringConsumerBackAndPublishMessage() error
}
