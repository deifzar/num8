package controller8

import (
	"deifzar/num8/pkg/model8"

	"github.com/gin-gonic/gin"
)

type Controller8NumateInterface interface {
	NumateScan(*gin.Context)
	NumateDomain(*gin.Context)
	NumateHostname(*gin.Context)
	NumateEndpoint(*gin.Context)
	HealthCheck(c *gin.Context)
	ReadinessCheck(c *gin.Context)
	// ConfigureEngine(model8.PostOptionsScan8) (model8.Model8Options8Interface, model8.Model8Results8Interface, error)
	ConfigureEngine(model8.PostOptionsScan8) (model8.Model8Options8Interface, string, error)
	RunNumate(bool, []model8.Endpoint8, model8.Model8Options8Interface, string)
	// CommitResults will insert the issues found into the database. This function internally parses the slice of `securityissues8` into a slice of `historyissue8` DB model.
	// Returns one boolean value that flags if new security issues have been found and one string value with the highest severity risk finding: critical, high or normal
	CommitResults(securityIssues []model8.SecurityIssues8, e8 []model8.Endpoint8) (bool, string, error)
	handleNotificationErrorOnFullscan(fullscan bool, message string, urgency string)
	// RunNumate([]model8.Httpendpoint8, model8.Model8Options8Interface, model8.Model8Results8Interface)
	// RabbitMQBringConsumerBackAndPublishMessage() error
}
