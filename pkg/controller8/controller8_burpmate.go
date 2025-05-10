package controller8

import (
	"bufio"
	"crypto/tls"
	"deifzar/num8/pkg/log8"
	"deifzar/num8/pkg/model8"
	"deifzar/num8/pkg/utils"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/itchyny/gojq"
)

type Controller8Burpmate struct {
	BurpAPILocation   string
	BurpProxyLocation string
}

func NewController8Burpmate(bAPI string, bProxy string) Controller8BurpmateInterface {
	return &Controller8Burpmate{
		BurpAPILocation:   bAPI,
		BurpProxyLocation: bProxy,
	}
}

func (m *Controller8Burpmate) GetSitemapByURLPrefix(urlPrefix string) (*model8.Sitemap8, error) {
	var s8 *model8.Sitemap8

	_, err := url.Parse(m.BurpAPILocation)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Error().Msg("Error parsing the URL where BurpM8 is located.")
		return s8, err
	}

	var myClient = &http.Client{Timeout: 10 * time.Second}

	response, err := myClient.Get(m.BurpAPILocation + "/burp/target/sitemap?urlPrefix=" + urlPrefix)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Error().Msgf("Error in GET request aiming to extract sitemap for %s.", urlPrefix)
		return s8, err
	}
	defer response.Body.Close()

	// query, err := gojq.Parse(". | group_by(.host) | map({host: first.host,per_host:(group_by(.port) | map({port: first.port, per_port:(group_by(.info.severity) | map({severity: first.info.severity, per_severity:(group_by(.template)|map({template: first.template,type:first.type,info:first.info.name, description:first.info.description, found:first.\"matched-at\"}))}))}))})")
	query, err := gojq.Parse("if (.messages[0].port == 443) then .messages | map(select((.statusCode >= 200 and .statusCode < 300) or .statusCode == 405 or .statusCode == 415) | {endpoint:.url | sub(\"(?<x>^[^/]+[^:]+):\\\\443(?<y>.*)\"; \"\\(.x)\\(.y)\"; \"\"), method:.method, parameters: [.parameters[] | select(.type == \"PARAM_URL\" or .type == \"PARAM_BODY\") | .name] | sort}) | {sitemap:.} else .messages | map(select((.statusCode >= 200 and .statusCode < 300) or .statusCode == 405 or .statusCode == 415) | {endpoint:.url, method:.method, parameters: [.parameters[] | select(.type == \"PARAM_URL\" or .type == \"PARAM_BODY\") | .name] | sort}) | {sitemap:.} end")
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Error().Msg("Something wrong with the 'jq' filter after extracting all the Burp sitemap.")
		return s8, err
	}

	var respJSON any
	err = json.NewDecoder(response.Body).Decode(&respJSON)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Error().Msgf("Error decoding JSON sitemap for %s.", urlPrefix)
		return s8, err
	}

	JSONOutput, err := utils.RunGoJQQuery(respJSON, query)
	if err != nil {
		log8.BaseLogger.Error().Msgf("Error running 'jq' filter in `GetSitemapByURLPrefix`")
		return nil, nil
	}
	b, err := json.Marshal(JSONOutput)
	if err != nil {
		log8.BaseLogger.Debug().Stack().Msg(err.Error())
		log8.BaseLogger.Warn().Msgf("Error with Marshal function. Origing: the Burp sitemap JSON object after applying 'jq' filter response for the URL: %s", urlPrefix)
		return nil, err

	}
	err = json.Unmarshal(b, &s8)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return nil, nil
	}
	//set hash for all resources found
	s8.SetHash()
	//uniq resources
	s8.Uniq()
	return s8, err
}

// Get GetFilteredSitemap will filter out the following resources in the Sitemap:
// Empty responses
// Listing "Content-type" resouces in the configuration file.
// Listing Response HTTP Codes in the configuration file.
func (m *Controller8Burpmate) GetSitemapFilteredOut(urlPrefix string, contenttype []string, statuscode []string) (*model8.Sitemap8, error) {
	var s8 *model8.Sitemap8

	_, err := url.Parse(m.BurpAPILocation)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Error().Msg("Error parsing the URL where BurpM8 is located.")
		return nil, err
	}

	var myClient = &http.Client{Timeout: 10 * time.Second}

	response, err := myClient.Get(m.BurpAPILocation + "/burp/target/sitemap?urlPrefix=" + urlPrefix)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Error().Msgf("Error in GET request aiming to extract sitemap for %s.", urlPrefix)
		return nil, err
	}
	defer response.Body.Close()

	var respJSON, respEmptyJSON any

	emptySitemap := []byte(`{"messages":[]}`)
	json.Unmarshal(emptySitemap, &respEmptyJSON)

	err = json.NewDecoder(response.Body).Decode(&respJSON)
	if err != nil {
		log8.BaseLogger.Debug().Stack().Msg(err.Error())
		log8.BaseLogger.Error().Msgf("Error decoding JSON sitemap for %s.", urlPrefix)
		return nil, err
	}

	// Check if sitemap is empty
	if !utils.EqualAny(respJSON, emptySitemap) {
		// Filter out empty responses
		filter, err := gojq.Parse("del(.messages[] | select(.responseHeaders | length == 0))")
		if err != nil {
			log8.BaseLogger.Debug().Stack().Msg(err.Error())
			log8.BaseLogger.Error().Msg("Something wrong with the 'jq' filter to clean the Sitemap of empty responses.")
			return nil, err
		}

		respJSON, err = utils.RunGoJQQuery(respJSON, filter)
		if err != nil {
			log8.BaseLogger.Error().Msgf("Error running 'jq' filter in `GetFilteredSitemap` when filtering out empty resources.")
			return nil, err
		}

		if !utils.EqualAny(respJSON, emptySitemap) {
			// Filter out content types responses
			lengthC := len(contenttype)
			lengthS := len(statuscode)
			i := 0
			j := 0
			for (i < lengthC) && (!utils.EqualAny(respJSON, emptySitemap)) {
				filter, err = gojq.Parse("del(.messages[] | select (.responseHeaders[] == \"" + contenttype[i] + "\"))")
				if err != nil {
					log8.BaseLogger.Debug().Stack().Msg(err.Error())
					log8.BaseLogger.Error().Msgf("Something wrong with the 'jq' filter to clean the Sitemap out of specific content-type: %s", contenttype[i])
					return nil, err
				}
				respJSON, err = utils.RunGoJQQuery(respJSON, filter)
				if err != nil {
					log8.BaseLogger.Error().Msgf("Error running 'jq' filter in `GetFilteredSitemap` when filtering out content-type %s", contenttype[i])
					return nil, err
				}
				i++
			}

			for (j < lengthS) && (!utils.EqualAny(respJSON, emptySitemap)) {
				// Filter out status code responses
				filter, err = gojq.Parse("del(.messages[] | select (.statusCode == " + statuscode[j] + "))")
				if err != nil {
					log8.BaseLogger.Debug().Stack().Msg(err.Error())
					log8.BaseLogger.Error().Msgf("Something wrong with the 'jq' filter to clean the Sitemap out of specific status-code: %s", statuscode[j])
					return nil, err
				}
				respJSON, err = utils.RunGoJQQuery(respJSON, filter)
				if err != nil {
					log8.BaseLogger.Error().Msgf("Error running 'jq' filter in `GetFilteredSitemap` when filtering out content-type %s", statuscode[j])
					return nil, err
				}
				j++
			}
		}
	}

	query, err := gojq.Parse("{sitemap: [.messages[] | if .port == 443 then {endpoint:.url | sub(\"(?<x>^[^/]+[^:]+):443(?<y>.*)\"; \"\\(.x)\\(.y)\"; \"\"), method:.method, parameters: [.parameters[] | select(.type == \"PARAM_URL\" or .type == \"PARAM_BODY\") | .name] | sort} | . else {endpoint:.url, method:.method, parameters: [.parameters[] | select(.type == \"PARAM_URL\" or .type == \"PARAM_BODY\") | .name] | sort} end ]}")
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Error().Msg("Something wrong with the 'jq' filter after filtering the Burp sitemap.")
		return nil, err
	}

	JSONOutput, err := utils.RunGoJQQuery(respJSON, query)
	if err != nil {
		log8.BaseLogger.Error().Msgf("Error running 'jq' filter in `GetFilteredSitemap` when cleaning out.")
		return nil, nil
	}
	b, err := json.Marshal(JSONOutput)
	if err != nil {
		log8.BaseLogger.Debug().Stack().Msg(err.Error())
		log8.BaseLogger.Warn().Msgf("Error with Marshal function. Origing: the Burp sitemap JSON object after applying 'jq' filter response for the URL: %s", urlPrefix)
		return nil, err

	}
	err = json.Unmarshal(b, &s8)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return nil, nil
	}
	//set hash for all resources found
	s8.SetHash()
	//uniq resources
	s8.Uniq()
	return s8, err
}

func (m *Controller8Burpmate) SendSitemap(sitemap *model8.Sitemap8) {

	proxy, _ := url.Parse(m.BurpProxyLocation)

	var myClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxy)},
		Timeout: 10 * time.Second,
	}
	for _, r8 := range sitemap.Sitemap {
		rawRequestExist := r8.RawRequest != ""
		if rawRequestExist {
			b := bufio.NewReader(strings.NewReader(r8.RawRequest))
			req, err := http.ReadRequest(b)
			if err != nil {
				rawRequestExist = false
			} else {
				u, err := url.Parse(r8.Endpoint)
				if err != nil {
					rawRequestExist = false
				}
				req.URL = u
				req.RequestURI = ""
				_, err = myClient.Do(req)
				if err != nil {
					rawRequestExist = false
				}
			}
		}
		if !rawRequestExist {
			myClient.Get(r8.Endpoint)
		}
	}
}
