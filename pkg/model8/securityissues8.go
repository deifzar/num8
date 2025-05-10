package model8

import (
	"bufio"
	"database/sql/driver"
	"deifzar/num8/pkg/log8"
	"deifzar/num8/pkg/utils"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/gofrs/uuid/v5"
	gojq "github.com/itchyny/gojq"
)

type Cookie8 struct {
	Domain     string    `json:"domain,omitempty"`
	Expiration time.Time `json:"expiration,omitempty"`
	Name       string    `json:"name,omitempty"`
	Path       string    `json:"path,omitempty"`
	Value      string    `json:"value,omitempty"`
}

type Parameter8 struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type HTTPMessage8 struct {
	Comment         string       `json:"comment,omitempty"`
	Cookies         []Cookie8    `json:"cookies,omitempty"`
	Highlight       string       `json:"highlight,omitempty"`
	Host            string       `json:"host,omitempty"`
	Method          string       `json:"method,omitempty"`
	Parameters      []Parameter8 `parameters:"url,omitempty"`
	Port            string       `json:"port,omitempty"`
	Protocol        string       `json:"protocol,omitempty"`
	Request         string       `json:"request,omitempty"`
	Response        string       `json:"response,omitempty"`
	ResponseHeaders []string     `json:"responseHeaders,omitempty"`
	StatusCode      int          `json:"statusCode,omitempty"`
	Url             string       `json:"url,omitempty"`
}

type Issue8 struct {
	IssueBackground       string         `json:"issueBackground,omitempty"`
	IssueDetail           string         `json:"issueDetail,omitempty"`
	IssueName             string         `json:"issueName,omitempty"`
	IssueType             string         `json:"issueType,omitempty"`
	Port                  int            `json:"port,omitempty"`
	Protocol              string         `json:"protocol,omitempty"`
	RemediationBackground string         `json:"remediationBackground,omitempty"`
	RemediationDetail     string         `json:"remediationDetail,omitempty"`
	Severity              string         `json:"severity,omitempty"`
	Confidence            string         `json:"confidence,omitempty"`
	Host                  string         `json:"host,omitempty"`
	HttpMessages          []HTTPMessage8 `json:"httpMessages,omitempty"`
	Url                   string         `json:"url,omitempty"`
}

type SecurityIssues8 struct {
	Url            string    `json:"url,omitempty"`
	Issues         []Issue8  `json:"issues"`
	HttpEndpointID uuid.UUID `json:"httpendpointID,omitempty"`
}

func (i Issue8) Value() (driver.Value, error) {
	return json.Marshal(i)
}

func (i *Issue8) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &i)
		return nil
	case string:
		json.Unmarshal([]byte(v), &i)
		return nil
	default:
		err := errors.New("issue8 - type assertion to []byte failed")
		log8.BaseLogger.Debug().Msg(err.Error())
		return err
	}
}

// ParseNum8ScanResults read the Nuclei results JSON file and parse those results into the `securityissue8` model.
// The JSON file is a line by line of JSON objects that needs to be converted into an array of JSON objects.
func ParseNum8ScanResults(resultfile string) ([]SecurityIssues8, error) {
	var err error

	query, err := gojq.Parse(". | group_by(.url) | map({url: first.url, httpendpointID:\"00000000-0000-0000-0000-000000000000\", issues:(group_by(.\"template-id\") | map({issueName: first.\"template-id\", severity: first.info.severity, host: first.host, port: (if first.port != null then first.port | tonumber else null end), protocol: first.scheme, issueBackground: first.info.description, issueType: first.type, issueDetail: ((first.info.description | (if . != null then \"Description: \" + . + \"\\\\r\\\\n\" else null end)) + (first.\"matcher-status\" | (if . == true then \"Matched found: Yes\\\\r\\\\n\" else \"Matched found: False\\\\r\\\\n\" end)) + (first.\"matcher-name\" | (if . != null then \"Found the following key indicator: \"+ . + \"\\\\r\\\\n\" else null end)) + (first.info.classification.\"cve-id\" | (if . != null then \"CVE: \"+ . + \"\\\\r\\\\n\" else null end)) + (first.info.reference | (if type == \"array\" then \"References: \" + (. | join(\", \")) + \"\\\\r\\\\n\" else null end)) + (first.\"template-id\" | (if . != null then \"Template ID: \" + . + \"\\\\r\\\\n\" else null end)) + (first.\"template-url\" | if . != null then \"Nuclei template URL: \"+ . + \"\\\\r\\\\n\" else null end) + (first.\"template-path\" | if . != null then \"Nuclei template path: \" + . + \"\\\\r\\\\n\" else null end)), remediationBackground: first.info.remediation | (if . != null then . else null end), remediationDetail: first.info.remediation | (if . != null then . else null end), httpMessages: (group_by(.request) | map({request: first.request | @base64, response: first.response | @base64}))}))})")
	if err != nil {
		log8.BaseLogger.Debug().Stack().Msg(err.Error())
		log8.BaseLogger.Error().Msgf("error with JQ query, file %s", resultfile)
		return nil, err
	}

	readFile, err := os.Open(resultfile)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Error().Msgf("error opening num8 results JSON file: %s", resultfile)
		return nil, err
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var results_num8 []any
	var line any
	for fileScanner.Scan() {
		err = json.Unmarshal(fileScanner.Bytes(), &line)
		if err != nil {
			log8.BaseLogger.Debug().Msg(err.Error())
			log8.BaseLogger.Error().Msgf("error unmarshaling lines of the num8 results JSON file: %s", resultfile)
			return nil, err
		}
		results_num8 = append(results_num8, line)
	}

	JSONOutput, err := utils.RunGoJQQuery(results_num8, query)
	if err != nil {
		log8.BaseLogger.Error().Msgf("error running 'jq' filter in `ParseScanResults` with JSON file: : %s", resultfile)
	}

	b, err := json.Marshal(JSONOutput)
	if err != nil {
		log8.BaseLogger.Debug().Stack().Msg(err.Error())
		log8.BaseLogger.Error().Msgf("error with Marshal function. Source: `ParseScanResults` with JSON file : %s", resultfile)
		return nil, err

	}
	var results []SecurityIssues8
	err = json.Unmarshal(b, &results)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		log8.BaseLogger.Error().Msgf("error with Unmarshal function. Source: `ParseScanResults` with JSON file : %s", resultfile)
		return nil, nil
	}

	return results, nil
}
