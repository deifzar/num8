package controller8

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"deifzar/num8/pkg/db8"
	"deifzar/num8/pkg/log8"
	"deifzar/num8/pkg/model8"
	"deifzar/num8/pkg/orchestrator8"
	"encoding/hex"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	outputnuclei "github.com/projectdiscovery/nuclei/v3/pkg/output"

	"github.com/gofrs/uuid/v5"
)

type Controller8Numate struct {
	Db   *sql.DB
	Cnfg *viper.Viper
}

func NewController8Numate(db *sql.DB, cnfg *viper.Viper) Controller8NumateInterface {
	return &Controller8Numate{Db: db, Cnfg: cnfg}
}

func (m *Controller8Numate) ConfigureEngine(p model8.PostOptionsScan8) (model8.Model8Options8Interface /*model8.Model8Results8Interface*/, string, error) {
	options8 := model8.NewModel8Options8()
	// results8 := model8.NewModel8Results8()
	if p.Options.Filters != nil {
		for _, f := range p.Options.Filters {
			options8.AddOption(nuclei.WithTemplateFilters(nuclei.TemplateFilters(f)))
		}
	}
	if p.Options.T != nil {
		turl := nuclei.WithTemplatesOrWorkflows(
			nuclei.TemplateSources{
				Templates: p.Options.T,
			},
		)
		options8.AddOption(turl)
	}
	if p.Options.TURL != nil {
		turl := nuclei.WithTemplatesOrWorkflows(
			nuclei.TemplateSources{
				RemoteTemplates: p.Options.TURL,
			},
		)
		options8.AddOption(turl)
	}
	if p.Options.W != nil {
		turl := nuclei.WithTemplatesOrWorkflows(
			nuclei.TemplateSources{
				Workflows: p.Options.W,
			},
		)
		options8.AddOption(turl)
	}
	if p.Options.WURL != nil {
		turl := nuclei.WithTemplatesOrWorkflows(
			nuclei.TemplateSources{
				RemoteWorkflows: p.Options.WURL,
			},
		)
		options8.AddOption(turl)
	}
	if m.Cnfg.GetString("NUM8.Proxy") != "" {
		options8.AddOption(nuclei.WithProxy([]string{m.Cnfg.GetString("NUM8.Proxy")}, false))
	}
	wo8 := model8.NewModel8WriterOptions8()
	outputFileName, err := wo8.SetDefaultWriterOptions8()
	if err != nil {
		return nil, "", err
	}
	// results8.SetOutputfilename(outputFileName)
	sw, err := outputnuclei.NewWriter(wo8.GetWriterOptions8()...)
	if err != nil {
		return nil, "", err
	}
	sw.DisableStdout = true
	options8.AddOption(nuclei.UseOutputWriter(sw))
	return options8, outputFileName, nil
}

// NumateScan will run only with what the configuration.yaml file contains: the `turl` property.
// The scans will run across all the domains and only through their root HTTP Endpoints
func (m *Controller8Numate) NumateScan(c *gin.Context) {
	// Check that RabbitMQ relevant Queue is available.
	// If relevant queue does not exist, inform the user that there is one Naabum8 running at this moment and advise the user to wait for the latest results.
	orchestrator8, err := orchestrator8.NewOrchestrator8()
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Fatal().Msg("Error connecting to the RabbitMQ server.")
	}
	amqp8 := orchestrator8.GetAmqp()
	queue_consumer := m.Cnfg.GetStringSlice("ORCHESTRATORM8.num8.Queue")
	qargs_consumer := m.Cnfg.GetStringMap("ORCHESTRATORM8.num8.Queue-arguments")
	exchange := m.Cnfg.GetStringSlice("ORCHESTRATORM8.asmm8.Queue")[0]
	if amqp8.ExistQueue(queue_consumer[1], qargs_consumer) {
		DB := m.Db
		endpoint8 := db8.NewDb8Endpoint8(DB)
		e8, err := endpoint8.GetAllHTTPEndpoints()
		if err != nil {
			// move on and call asmm8 scan
			orchestrator8.PublishMessageToExchangeAndCloseChannelConnection(exchange, "cptm8.asmm8.get.scan")
			c.JSON(http.StatusBadGateway, gin.H{"status": "error", "msg": "Num8 Scan failed - Error fetching the endpoints."})
			log8.BaseLogger.Debug().Msg(err.Error())
			log8.BaseLogger.Warn().Msg("HTTP Repose 500 - Num8 Scan failed - Error fetching the endpoints.")
			return
		}
		if len(e8) < 1 {
			// move on and call asmm8 scan
			orchestrator8.PublishMessageToExchangeAndCloseChannelConnection(exchange, "cptm8.asmm8.get.scan")
			c.JSON(http.StatusOK, gin.H{"msg": "Num8 scan API call success. No targets in scope."})
			log8.BaseLogger.Info().Msg("Num8 scan API call success. No targets in scope")
			return
		}
		var post model8.PostOptionsScan8
		// post.Options.T = m.Cnfg.GetStringSlice("NUM8.t")
		// post.Options.TURL = m.Cnfg.GetStringSlice("NUM8.turl")
		options8, outputFileName, err := m.ConfigureEngine(post)
		if err != nil {
			// move on and call asmm8 scan
			orchestrator8.PublishMessageToExchangeAndCloseChannelConnection(exchange, "cptm8.asmm8.get.scan")
			c.JSON(http.StatusBadGateway, gin.H{"status": "error", "msg": "Numate Scan failed - scan configuration failed"})
			log8.BaseLogger.Debug().Msg(err.Error())
			log8.BaseLogger.Warn().Msg("HTTP Repose 500 - Num8 Scan failed - scan configuration failed")
			return
		}
		// // cancel consumer
		// err = orchestrator8.DeactivateConsumerByService("num8")
		// if err != nil {
		// 	// move on and call asmm8 scan
		// 	orchestrator8.PublishMessageToExchangeAndCloseChannelConnection(exchange, "cptm8.asmm8.get.scan")
		// 	log8.BaseLogger.Error().Msg("HTTP 500 Response - Num8 Full scans failed - Error cancelling the RabbitMQ consumer for `num8` before launching scan.")
		// 	c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": "Num8 Full scans failed. Error cancelling the RabbitMQ consumer."})
		// 	return
		// }
		// bring the queue back with no consumers
		err = orchestrator8.ActivateQueueByService("num8")
		if err != nil {
			// move on and call asmm8 scan
			orchestrator8.PublishMessageToExchangeAndCloseChannelConnection(exchange, "cptm8.asmm8.get.scan")
			log8.BaseLogger.Fatal().Msg("HTTP 500 Response - Num8 Scans failed - Error bringing up the RabbitMQ queues for the Num8 service.")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": "Num8 Scans failed. Error bringing up the RabbitMQ queues for the Num8 service."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Num8 scan API call success... Check notifications for scan updates."})
		log8.BaseLogger.Info().Msg("Num8 scan API call success.")
		// run active.
		go m.RunNumate(true, orchestrator8, e8, options8, outputFileName)
	} else {
		// move on and call asmm8 scan
		orchestrator8.PublishMessageToExchangeAndCloseChannelConnection(exchange, "cptm8.asmm8.get.scan")
		c.JSON(http.StatusForbidden, gin.H{"status": "forbidden", "msg": "Num8 Scans failed - Launching Num8 Scan is not possible at this moment due to non-existent RabbitMQ queues."})
		log8.BaseLogger.Info().Msg("Num8 Scan API call forbidden")
		return
	}
}

// NumateDomain will run with what the configuration.yaml file contains and the POST options.
// POST body contains more particular options for the scan. Among those possible options, Workflow scans are of special interest.
// The scans will run through only the root/base HTTP Endpoints.
func (m *Controller8Numate) NumateDomain(c *gin.Context) {
	DB := m.Db
	var post model8.PostOptionsScan8
	var uri model8.Domain8Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Num8 Scan Domain failed - Check URL parameters."})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 400 - Num8 Scan Domain failed - Check URL parameters.")
		return
	}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Num8 scan domain failed - Check body request."})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 400 - Num8 Scan Domain failed - Check body request.")
		return
	}
	endpoint8 := db8.NewDb8Endpoint8(DB)
	id, err := uuid.FromString(uri.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Num8 Scan Domain failed - Check UUID URL parameters."})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 400 - Num8 Scan Domain failed - Check UUID URL parameters.")
		return
	}
	e8, err := endpoint8.GetAllHTTPByDomainID(id)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "msg": "Num8 Scan Domain failed - Somehing wrong fetching all endpoints by domain"})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 500 - Num8 Scan Domain failed - Something wrong fetching endpoints by domain")
		return
	}
	if len(e8) < 1 {
		c.JSON(http.StatusOK, gin.H{"msg": "OK! Num8 Scan Domain are about to commence.\nCheck Discord for any update about the scans."})
		log8.BaseLogger.Info().Msg("Num8 Scan Domain success. No domain in scope.")
		return
	}
	options8, outputFileName, err := m.ConfigureEngine(post)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "msg": "Num8 scan domain failed - scan configuration failed"})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 500 - Num8 scan domain failed - Configure scan failed")
		return
	}
	go m.RunNumate(false, nil, e8, options8, outputFileName)
	c.JSON(http.StatusOK, gin.H{"msg": "OK! Num8 Scan Domain are about to commence.\nCheck Discord for any update about the scans."})
	log8.BaseLogger.Info().Msg("Num8 Scan Domain running.")
}

// NumateHostname will run with what the configuration.yaml file contains and the POST options.
// POST body contains more particular options for the scan. Among those possible options, Workflow scans are of special interest.
// The scans will run through each in-scope URL resource found in Burp sitemap.
func (m *Controller8Numate) NumateHostname(c *gin.Context) {
	DB := m.Db
	var post model8.PostOptionsScan8
	var uri model8.Hostname8Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Num8 Scan Hostname failed - Check URL parameters."})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 400 - Num8 Scan Hostname failed - Check URL parameters.")
		return
	}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Num8 Scan Hostname failed - Check body request."})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 400 - Num8 Scan Hostname failed - Check body request.")
		return
	}
	endpoint8 := db8.NewDb8Endpoint8(DB)
	id, err := uuid.FromString(uri.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Num8 Scan Hostname failed - Check UUID URL parameters."})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 400 - Num8 Scan Hostname failed - Check UUID URL parameters.")
		return
	}
	e8, err := endpoint8.GetAllHTTPByHostnameID(id)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "msg": "Num8 Scan Hostname failed - something wrong fetching endpoints by hostname"})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 500 - Num8 Scan Hostname failed - something wrong fetching endpoints by hostname")
		return
	}
	if len(e8) < 1 {
		c.JSON(http.StatusOK, gin.H{"msg": "OK! Num8 Scan Hostname are about to commence.\nCheck Discord for any update about the scans."})
		log8.BaseLogger.Info().Msg("Num8 Scan Hostname success. No domain in scope.")
		return
	}
	options8, outputFileName, err := m.ConfigureEngine(post)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "msg": "Num8 Scan Hostname failed - Configure scan failed"})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 500 - Num8 Scan Hostname failed - scan configuration failed")
		return
	}
	go m.RunNumateThoroughly(e8, options8, outputFileName)
	c.JSON(http.StatusOK, gin.H{"msg": "Num8 Scan Hostname running. Check Discord for any update about the scans."})
	log8.BaseLogger.Info().Msg("Num8 Scan Hostname running.")
}

// NumateHostname will run with what the configuration.yaml file contains and the POST options.
// POST body contains more particular options for the scan. Among those possible options, Workflow scans are of special interest.
// The scans will run through each in-scope URL resource found in Burp sitemap.
func (m *Controller8Numate) NumateEndpoint(c *gin.Context) {
	DB := m.Db
	var post model8.PostOptionsScan8
	var uri model8.Endpoint8Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Num8 Scan Endpoint failed - Check URL parameters."})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 400 - Num8 Scan Endpoint failed - Check URL parameters.")
		return
	}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Num8 Scan Endpoint failed - Check body parameters."})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 400 - Num8 Scan Endpoint failed - Check body parameters.")
		return
	}
	id, err := uuid.FromString(uri.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Num8 Scan Endpoint failed - Check UUID URL parameters."})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 400 - Num8 Scan Endpoint failed - Check URL parameters.")
		return
	}
	endpoint8 := db8.NewDb8Endpoint8(DB)
	e8, err := endpoint8.GetOneEndpointByID(id)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "msg": "Num8 Scan Endpoint failed - something wrong fetching the endpoint"})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 500 - Num8 Scan Endpoint failed - Failed to get endpoint")
		return
	}
	if reflect.ValueOf(e8).IsZero() {
		c.JSON(http.StatusOK, gin.H{"msg": "OK! Num8 Scan Endpoint are about to commence.\nCheck Discord for any update about the scans."})
		log8.BaseLogger.Info().Msg("Num8 Scan Endpoint success. No domain in scope.")
		return
	}
	options8, outputFileName, err := m.ConfigureEngine(post)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "msg": "Num8 Scan Endpoint failed - scan configuration failed"})
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Warn().Msg("HTTP Repose 500 - Num8 Scan Endpoint failed - Configure scan failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "OK! Num8 Scan Endpoint are about to commence.\nCheck Discord for any update about the scans."})
	log8.BaseLogger.Info().Msg("Num8 Scan Endpoint running.")
	go m.RunNumateThoroughly([]model8.Endpoint8{e8}, options8, outputFileName)
}

// func (m *Controller8Numate) RunNumate(e8 []model8.Httpendpoint8, o8 model8.Model8Options8Interface, r8 model8.Model8Results8Interface) {
// 	defer m.RabbitMQBringConsumerBack()

// 	num8Discord := model8.NewModel8Discord8(
// m.Config.GetString("Discord.webhookURL"),
// m.Config.GetString("Discord.webhookID"),
// m.Config.GetString("Discord.webhookName"),
// m.Config.GetString("Discord.webhookToken"),
// m.Config.GetString("Discord.botToken"))
// 	err := num8Discord.InitialiseChannelID()
// 	if err != nil {
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		return
// 	}
// 	discordBot, err := discordgo.New("Bot " + num8Discord.GetBotToken())
// 	if err != nil {
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		return
// 	}
// 	discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "Num8 scans are running!")
// 	// create nuclei engine with options
// 	ne, err := nuclei.NewNucleiEngine(o8.GetOptions()...)
// 	if err != nil {
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		discordBot.ChannelMessageSend(num8Discord.GetChannelID(), err.Error())
// 		return
// 	}
// 	defer ne.Close()
// 	// load targets and optionally probe non http/https targets
// 	var targets []string
// 	for _, e := range e8 {
// 		targets = append(targets, e.Endpoint)
// 	}
// 	ne.LoadTargets(targets, false)
// 	err = ne.ExecuteWithCallback(nil)
// 	if err != nil {
// 		// panic(err)
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		discordBot.ChannelMessageSend(num8Discord.GetChannelID(), err.Error())
// 		return
// 	}
// 	discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "Num8 scans have finished!")
// 	if err = r8.SetResultEventFromOutputfilename(); err != nil {
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		discordBot.ChannelMessageSend(num8Discord.GetChannelID(), err.Error())
// 		return
// 	}
// 	if err = num8Discord.SetChatMessages(r8.GetResultEvent()); err != nil {
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		discordBot.ChannelMessageSend(num8Discord.GetChannelID(), err.Error())
// 		return
// 	}
// 	// log.Println(num8Discord.ChatMessages)
// 	if num8Discord.GetChatMessages() != nil {
// 		for _, messagesHost := range num8Discord.GetChatMessages() {
// 			if messagesHost.Per_host != nil {
// 				discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "------------------------------------------------------------------")
// 				discordBot.ChannelMessageSend(num8Discord.GetChannelID(), ":desktop:  Vulnerabilities for host: "+messagesHost.Host+"  :desktop:")
// 				discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "------------------------------------------------------------------")
// 				for _, messagesHostPort := range messagesHost.Per_host {
// 					if messagesHostPort.Per_port != nil {
// 						for _, messagesHostPortSeverity := range messagesHostPort.Per_port {
// 							if messagesHostPortSeverity.Per_severity != nil {
// 								var colorEmoji string
// 								switch riskLevel := messagesHostPortSeverity.Severity; riskLevel {
// 								case "info":
// 									colorEmoji = ":blue_circle:"
// 								case "low":
// 									colorEmoji = ":green_circle:"
// 								case "medium":
// 									colorEmoji = ":orange_circle:"
// 								case "high":
// 									colorEmoji = ":red_circle:"
// 								case "criticial":
// 									colorEmoji = ":exclamation::red_circle:"
// 								default:
// 									colorEmoji = ":question:"

// 								}
// 								numFindginsPerSeverity := strconv.Itoa(len(messagesHostPortSeverity.Per_severity))
// 								discordBot.ChannelMessageSend(num8Discord.GetChannelID(), ":space_invader: **"+strings.ToUpper(numFindginsPerSeverity)+"** "+messagesHostPortSeverity.Severity+" severity risk level findings for "+messagesHost.Host+":"+messagesHostPort.Port)
// 								discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "------------------------------------------------------------------")
// 								for _, messagesHostPortSeverityFinding := range messagesHostPortSeverity.Per_severity {
// 									discordBot.ChannelMessageSend(num8Discord.GetChannelID(), colorEmoji+" Template ("+messagesHostPortSeverityFinding.Type+"): "+messagesHostPortSeverityFinding.Template)
// 									discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "Info: "+messagesHostPortSeverityFinding.Info)
// 									discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "Description: "+messagesHostPortSeverityFinding.Description)
// 								}
// 							} else {
// 								discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "0 "+messagesHostPortSeverity.Severity+" severity risk level finding for "+messagesHost.Host+":"+messagesHostPort.Port)
// 							}
// 						}
// 					} else {
// 						discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "Empty results for "+messagesHost.Host+":"+messagesHostPort.Port)
// 					}
// 				}
// 			} else {
// 				discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "Empty results for host: "+messagesHost.Host)
// 			}
// 		}
// 	} else {
// 		discordBot.ChannelMessageSend(num8Discord.GetChannelID(), "No results found.")
// 	}
// 	// RabbitMQ publish message for ASMM8
// 	m.RabbitMQPublishMessage()
// }

func (m *Controller8Numate) RunNumate(fullscan bool, orch8 orchestrator8.Orchestrator8Interface, e8 []model8.Endpoint8, o8 model8.Model8Options8Interface, outputFileName string) {

	log8.BaseLogger.Info().Msg("Num8 scans are running!")
	// create nuclei engine with options
	ne, err := nuclei.NewNucleiEngineCtx(context.Background(), o8.GetOptions()...)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return
	}
	defer ne.Close()
	// load targets and optionally probe non http/https targets
	var targets []string
	for _, e := range e8 {
		targets = append(targets, e.Endpoint)
	}
	ne.LoadTargets(targets, false)
	err = ne.ExecuteWithCallback(nil)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Info().Msg("there was an error in the execution of nuclei")
		return
	}
	log8.BaseLogger.Info().Msg("Num8 scans have finished!")
	// outputFileName = "./tmp/result-2025-3-4-11-30-14559570398"
	securityIssues, err := model8.ParseNum8ScanResults(outputFileName)
	if err != nil {
		log8.BaseLogger.Error().Msg("Num8 scan results parsed with errors")
		log8.BaseLogger.Debug().Msg(err.Error())
		return
	} else {
		log8.BaseLogger.Info().Msg("Num8 scan results parsed ok!")
	}

	// Commit latest Num8 results into the DB
	err = m.CommitResults(securityIssues, e8)
	if err != nil {
		log8.BaseLogger.Error().Msg("error after attempt to commit results")
		return
	}
	if fullscan {
		// call asmm8 scan
		exchange := m.Cnfg.GetStringSlice("ORCHESTRATORM8.asmm8.Queue")[0]
		orch8.PublishMessageToExchangeAndCloseChannelConnection(exchange, "cptm8.asmm8.get.scan")
	}
}

func (m *Controller8Numate) RunNumateThoroughly(e8 []model8.Endpoint8, o8 model8.Model8Options8Interface, outputFileName string) {

	log8.BaseLogger.Info().Msg("Num8 scans are running!")
	// create nuclei engine with options
	ne, err := nuclei.NewNucleiEngineCtx(context.Background(), o8.GetOptions()...)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return
	}
	defer ne.Close()
	var targets []string

	controllerBurpmate := NewController8Burpmate(m.Cnfg.GetString("NUM8.BurpAPILocation"), m.Cnfg.GetString("NUM8.BurpProxyLocation"))
	// BurpController -> get and filter the sitemap with worthly resources to scan
	contenttype := m.Cnfg.GetStringSlice("NUM8.Sitemap.Filter.ContentType")
	statuscode := m.Cnfg.GetStringSlice("NUM8.Sitemap.Filter.StatusCode")
	for _, endpoint := range e8 {
		sitemap, err := controllerBurpmate.GetSitemapFilteredOut(endpoint.Endpoint, contenttype, statuscode)
		if err != nil {
			log8.BaseLogger.Debug().Msg(err.Error())
			continue
		}
		for _, resource := range sitemap.Sitemap {
			targets = append(targets, resource.Endpoint)
		}
	}
	ne.LoadTargets(targets, false)
	err = ne.ExecuteWithCallback(nil)
	if err != nil {
		log8.BaseLogger.Debug().Stack().Msg(err.Error())
		return
	}
	log8.BaseLogger.Info().Msg("Num8 scans have finished!")

	securityIssues, err := model8.ParseNum8ScanResults(outputFileName)
	if err != nil {
		log8.BaseLogger.Debug().Stack().Msg(err.Error())
		return
	}

	// Commit latest Num8 results into the DB
	err = m.CommitResults(securityIssues, e8)
	if err != nil {
		log8.BaseLogger.Error().Msg("error after attempt to commit results")
		return
	}
}

// CommitResults will insert the issues found into the database. This function internally parses the slice of `securityissues8` into a slice of `historyissue8` DB model.
func (m *Controller8Numate) CommitResults(securityIssues []model8.SecurityIssues8, e8 []model8.Endpoint8) error {
	var historyissues []model8.Historyissue8
	// prepare historyissues slice
	for _, si := range securityIssues {
		var baseURL, baseURL2 string
		if si.Url != "" {
			u, err := url.Parse(si.Url)
			if err != nil {
				log8.BaseLogger.Debug().Stack().Msg(err.Error())
				log8.BaseLogger.Error().Msgf("Error parsing URL `%s`", si.Url)
				continue
			}
			baseURL = u.Scheme + "://" + u.Host
			// Find Httpendpoint ID
			for _, e := range e8 {
				if e.Endpoint == baseURL {
					si.HttpEndpointID = e.Id
					break
				}
			}
			// if the httpendpoint is part of the scope for this test round
			if !si.HttpEndpointID.IsNil() {
				params, err := url.ParseQuery(u.RawQuery)
				if err != nil {
					log8.BaseLogger.Debug().Stack().Msg(err.Error())
					log8.BaseLogger.Error().Msgf("Error parsing URL parameters: `%s`", si.Url)
					continue
				}
				// create issue signature
				var prefix_signature string
				var signature_issue string
				h := sha512.New()
				prefix_signature = baseURL + "," + u.Path
				for p := range params {
					prefix_signature = prefix_signature + "," + p
				}
				for _, issue := range si.Issues {
					signature_issue = prefix_signature + "|" + issue.IssueName + "|" + issue.IssueDetail
					h.Write([]byte(signature_issue))
					hash := h.Sum(nil)
					hi := model8.Historyissue8{
						Endpointid: si.HttpEndpointID,
						Url:        si.Url,
						Signature:  hex.EncodeToString(hash),
						Issue:      issue,
						Status:     model8.Unreviewed,
					}
					historyissues = append(historyissues, hi)
				}
			}
		} else {
			for _, issue := range si.Issues {
				if issue.Port > 0 {
					switch issue.Port {
					case 80:
						baseURL = "http://" + issue.Host
					case 443:
						baseURL = "https://" + issue.Host
					default:
						baseURL = "http://" + issue.Host + ":" + strconv.Itoa(issue.Port)
						baseURL2 = "https://" + issue.Host + ":" + strconv.Itoa(issue.Port)
					}
				} else {
					baseURL = "http://" + issue.Host
					baseURL2 = "https://" + issue.Host
				}
				for _, e := range e8 {
					if e.Endpoint == baseURL || e.Endpoint == baseURL2 {
						baseURL = e.Endpoint
						si.HttpEndpointID = e.Id
						break
					}
				}
				// if the url is part of the scope for this test round
				if !si.HttpEndpointID.IsNil() {
					signature := baseURL + "|" + issue.IssueName + "|" + issue.IssueDetail
					h := sha512.New()
					h.Write([]byte(signature))
					hash := h.Sum(nil)
					hi := model8.Historyissue8{
						Endpointid: si.HttpEndpointID,
						Url:        baseURL,
						Signature:  hex.EncodeToString(hash),
						Issue:      issue,
						Status:     model8.Unreviewed,
					}
					historyissues = append(historyissues, hi)
				}
			}
		}
	}
	dbHistoryissue8 := db8.NewDb8Historyissue8(m.Db)

	// Fetch all FP signatures from database and delete the ones found from the list
	currentHistoryIusses, err := dbHistoryissue8.GetAllHistoryIssuesByStatus(model8.Falsepositive)
	if err != nil {
		log8.BaseLogger.Error().Msgf("Error fetching False Positives security issues from DB")
		return err
	}
	// Fetch all I signatures from database and delete the ones found from the list
	ignored, err := dbHistoryissue8.GetAllHistoryIssuesByStatus(model8.Ignored)
	if err != nil {
		log8.BaseLogger.Error().Msgf("Error fetching Ignored security issues from DB")
		return err
	}
	currentHistoryIusses = append(currentHistoryIusses, ignored...)
	// Fetch all V signatures from database and delete the ones found from the list
	verified, err := dbHistoryissue8.GetAllHistoryIssuesByStatus(model8.Verified)
	if err != nil {
		log8.BaseLogger.Error().Msgf("Error fetching Verified security issues from DB")
		return err
	}
	currentHistoryIusses = append(currentHistoryIusses, verified...)
	historyissues = model8.DifferenceHistoryissues8(historyissues, currentHistoryIusses)
	err = dbHistoryissue8.InsertBatch(historyissues)
	if err != nil {
		log8.BaseLogger.Error().Msgf("Error inserting the new security issues found into the DB")
	}
	return err
}

// func (m *Controller8Numate) RabbitMQBringConsumerBackAndPublishMessage() error {
// 	// RabbitMQ queue and consumer for asmm8 should be back to be available.
// 	orchestrator8, err := orchestrator8.NewOrchestrator8()
// 	amqp8 := orchestrator8.GetAmqp()
// 	defer amqp8.CloseChannel()
// 	defer amqp8.CloseConnection()
// 	if err != nil {
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		log8.BaseLogger.Fatal().Msg("Error connecting to the RabbitMQ server.")
// 		return err
// 	}
// 	orchestrator8.CreateHandleCPTM8()
// 	orchestrator8.ActivateConsumerByService("num8")

// 	// Publish message
// 	queue := m.Cnfg.GetStringSlice("ORCHESTRATORM8.asmm8.Queue")
// 	log8.BaseLogger.Info().Msg("RabbitMQ publishing message to queue for ASMM8 service.")
// 	err = amqp8.Publish(queue[0], "cptm8.asmm8.get.scan", "")
// 	if err != nil {
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		log8.BaseLogger.Error().Msgf("rabbitMQ publishing message to queue for ASMM8 service failed")
// 		return err
// 	}
// 	log8.BaseLogger.Info().Msg("RabbitMQ publishing message to ASMM8 queue service success.")
// 	return nil
// }

// func (m *Controller8Numate) RabbitMQBringConsumerBack() error {

// 	// RabbitMQ queue and consumer for num8 should be back to be available.
// 	orchestrator8, err := orchestrator8.NewOrchestrator8()
// 	if err != nil {
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		log8.BaseLogger.Fatal().Msg("Error connecting to the RabbitMQ server.")
// 		return err
// 	}
// 	orchestrator8.CreateHandleCPTM8()
// 	orchestrator8.ActivateConsumerByService("num8")
// 	return nil
// }

// func (m *Controller8Numate) RabbitMQPublishMessage() error {
// 	orchestrator8, err := orchestrator8.NewOrchestrator8()
// 	if err != nil {
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		log8.BaseLogger.Fatal().Msg("Error connecting to the RabbitMQ server.")
// 		return err
// 	}
// 	amqp8 := orchestrator8.GetAmqp()
// 	defer amqp8.CloseChannel()
// 	defer amqp8.CloseConnection()
// 	queue := m.Cnfg.GetStringSlice("ORCHESTRATORM8.asmm8.Queue")
// 	log8.BaseLogger.Info().Msg("RabbitMQ publishing message to queue for ASMM8 service.")
// 	err = amqp8.Publish(queue[0], "cptm8.asmm8.get.scan", "")
// 	if err != nil {
// 		log8.BaseLogger.Debug().Msg(err.Error())
// 		log8.BaseLogger.Error().Msgf("rabbitMQ publishing message to queue for ASMM8 service failed")
// 		return err
// 	}
// 	log8.BaseLogger.Info().Msg("RabbitMQ publishing message to ASMM8 queue service success.")
// 	return nil
// }
